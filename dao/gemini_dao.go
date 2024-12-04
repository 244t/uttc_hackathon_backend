package dao

import (
	"context"
	"fmt"
    "encoding/json"
    "myproject/model"
    "math"
    "sync"

	"github.com/google/generative-ai-go/genai"
    "cloud.google.com/go/storage"

)

type VertexAiDAO struct {
	Client *genai.Client
}

type VertexAiDAOInterface interface {
	NextTextGeneration(ctx context.Context, text string) (*genai.Part, error)
    EmbeddingGeneration(ctx context.Context, er model.EmbeddingRequest) (error)
}

func NewVertexAiDAO(client *genai.Client) *VertexAiDAO {
	return &VertexAiDAO{
		Client: client,
	}
}

func (dao *VertexAiDAO) NextTextGeneration(ctx context.Context, promptText string) (*genai.Part, error) {
    gemini := dao.Client.GenerativeModel("gemini-1.5-flash-002")
    prompt := genai.Text(promptText)
    resp, err := gemini.GenerateContent(ctx, prompt)
    if err != nil {
        return nil, fmt.Errorf("error generating content: %w", err)
    }

    // Candidatesの中から最初のものを取得
    if len(resp.Candidates) == 0 {
        return nil, fmt.Errorf("no generated content received")
    }

    // ContentのPartsに格納された生成されたテキストを取得
    if len(resp.Candidates[0].Content.Parts) == 0 {
        return nil, fmt.Errorf("no content found in response")
    }

    // Parts[0]の情報をそのまま返す
    part := resp.Candidates[0].Content.Parts[0]
    return &part, nil
}

func (dao *VertexAiDAO) EmbeddingGeneration(ctx context.Context, er model.EmbeddingRequest) error{
    em := dao.Client.EmbeddingModel("text-embedding-004")
    // エンベディングを生成
	res, err := em.EmbedContent(ctx, genai.Text(er.Content))
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}
    //embeddingベクトルを取得
	vec := res.Embedding.Values

    err = saveEmbeddingToGCS(ctx, vec, er.UserId)
    if err != nil {
        return fmt.Errorf("failed to save embedding to GCS: %w", err)
    }

    return nil
}

// GCSにファイルが存在するかを確認する関数
func fileExistsInGCS(ctx context.Context, bucketName, objectName string) (bool, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	// オブジェクトの属性を取得して、エラーがErrObjectNotExistであればファイルは存在しない
	_, err = obj.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if object exists: %v", err)
	}

	return true, nil
}

// GCSから既存のembeddingデータを取得する関数
func getEmbeddingFromGCS(ctx context.Context, bucketName, objectName string) (*model.EmbeddingResult, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	// オブジェクトのデータを読み込む
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read object from GCS: %v", err)
	}
	defer reader.Close()

	var embeddingResult model.EmbeddingResult
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&embeddingResult)
	if err != nil {
		return nil, fmt.Errorf("failed to decode embedding data: %v", err)
	}

	return &embeddingResult, nil
}

// ベクトルの平均を計算する関数（並列処理版）
func averageVectorsSync(existingVec, newVec []float32, count int) []float32 {
	if len(existingVec) != len(newVec) {
		return nil
	}

	// 結果用のベクトルを作成
	avg := make([]float32, len(existingVec))

	// 並列処理のためのWaitGroup
	var wg sync.WaitGroup
	// 複数のゴルーチンを使用して並列処理
	numGoroutines := 4 // ゴルーチンの数（適宜調整）
	chunkSize := len(existingVec) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()

			// チャンクごとの処理を行う
			endIdx := startIdx + chunkSize
			if i == numGoroutines-1 {
				endIdx = len(existingVec) // 最後のチャンクは残りの部分を処理
			}

			for j := startIdx; j < endIdx; j++ {
				avg[j] = (existingVec[j]*float32(count) + newVec[j]) / float32(count+1)
			}
		}(i * chunkSize)
	}

	// ゴルーチンが全て完了するまで待機
	wg.Wait()

	return avg
}

// GCSにembeddingデータを保存する関数（ファイルが存在しない場合のみ新規作成、存在する場合は更新）
func saveEmbeddingToGCS(ctx context.Context, embedding []float32, userId string) error {
	bucketName := "term6-kyosuke-nishishita.firebasestorage.app"
	objectName := fmt.Sprintf("vector/%s.json", userId)

	// 既存のファイルがあるかを確認
	exists, err := fileExistsInGCS(ctx, bucketName, objectName)
	if err != nil {
		return fmt.Errorf("failed to check if file exists in GCS: %w", err)
	}

	var embeddingData model.EmbeddingResult

	if exists {
		// 既存のファイルがある場合、GCSからデータを取得
		existingData, err := getEmbeddingFromGCS(ctx, bucketName, objectName)
		if err != nil {
			return fmt.Errorf("failed to retrieve existing embedding data: %w", err)
		}

		// ベクトルの平均を計算
		averageVec := averageVectorsSync(existingData.Embedding, embedding, existingData.Count)
		if averageVec == nil {
			return fmt.Errorf("embedding vectors have different lengths")
		}

		// 新しいembeddingデータを更新
		embeddingData = model.EmbeddingResult{
			UserID:   userId,
			Count:    existingData.Count + 1, // カウントを1増加
			Embedding: averageVec,
		}
	} else {
		// 新規に保存する場合は、埋め込みデータをそのまま使用
		embeddingData = model.EmbeddingResult{
			UserID:   userId,
			Count:    1, // 初回のカウントは1
			Embedding: embedding,
		}
	}

	// 新規または更新されたデータをGCSに保存
	return saveToGCS(ctx, embeddingData, bucketName, objectName)
}

// GCSにembeddingデータを保存する関数
func saveToGCS(ctx context.Context, result model.EmbeddingResult, bucketName, objectName string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)

	// embeddingデータをJSONにエンコード
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal embedding result data: %v", err)
	}

	// GCSに保存するオブジェクトの参照を取得
	object := bucket.Object(objectName)

	// オブジェクトに書き込むためのwriterを作成
	writer := object.NewWriter(ctx)
	defer writer.Close()

	// データを書き込む
	n, err := writer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to GCS: %v", err)
	}

	// 書き込まれたバイト数を表示
	fmt.Printf("Successfully wrote %d bytes to GCS\n", n)
	return nil
}

// コサイン類似度を計算する関数
func cosineSimilarity(vec1, vec2 []float32) (float32, error) {
	if len(vec1) != len(vec2) {
		return 0, fmt.Errorf("vectors must be of the same length")
	}

	// 内積を計算
	var dotProduct float32
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
	}

	// 各ベクトルのノルムを計算
	var normVec1, normVec2 float32
	for i := 0; i < len(vec1); i++ {
		normVec1 += vec1[i] * vec1[i]
		normVec2 += vec2[i] * vec2[i]
	}

	// ノルムの平方根を取る
	normVec1 = float32(math.Sqrt(float64(normVec1)))
	normVec2 = float32(math.Sqrt(float64(normVec2)))

	// コサイン類似度を計算
	if normVec1 == 0 || normVec2 == 0 {
		return 0, fmt.Errorf("one of the vectors is zero-length")
	}
	return dotProduct / (normVec1 * normVec2), nil
}