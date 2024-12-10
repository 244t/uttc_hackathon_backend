package dao

import (
	"context"
	"fmt"
    "encoding/json"
    "myproject/model"
    "math"
	"database/sql"
    "sync"
	"sort"
	"log"
	

	"github.com/google/generative-ai-go/genai"
    "cloud.google.com/go/storage"

)

type VertexAiDAO struct {
	Client *genai.Client
	DB *sql.DB
}

type VertexAiDAOInterface interface {
	RegisterUser(ctx context.Context,user model.Profile) error
	NextTextGeneration(ctx context.Context, text string) (*genai.Part, error)
    EmbeddingGeneration(ctx context.Context, er model.EmbeddingRequest) (error)
	FindSimilar(ctx context.Context, fs model.FindSimilarRequest) ([]model.Profile,error)
	GetUserProfile(userId string) (model.Profile, error)
	RecommendUser(ctx context.Context,userId string)([]model.Profile,error)
}

func NewVertexAiDAO(client *genai.Client, db *sql.DB) *VertexAiDAO {
	return &VertexAiDAO{
		Client: client,
		DB : db,
	}
}

func (dao *VertexAiDAO) RegisterUser(ctx context.Context,user model.Profile) error{
	_ ,err := dao.DB.Exec("INSERT INTO user (user_id, name, bio,profile_img_url,header_img_url,location) VALUES (?, ?, ?,?,?,?)", user.Id, user.Name, user.Bio,user.ImgUrl,user.HeaderUrl,user.Location)
	content := fmt.Sprintf("私は %sです。自己紹介は%s",user.Name, user.Bio)

	er := model.EmbeddingRequest{
		UserId : user.Id,
		Content : content,
	}
	dao.EmbeddingGeneration(ctx,er)
	return err
}

//ツイートの続きを生成する関数
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

//embedding
func (dao *VertexAiDAO) EmbeddingGeneration(ctx context.Context, er model.EmbeddingRequest) error{
	// Step 1: er.Contentを英訳
    gemini := dao.Client.GenerativeModel("gemini-1.5-flash-002")
    prompt := genai.Text("Please return only the translated text."+"Translate the following tweet to English: " + er.Content )
    resp, err := gemini.GenerateContent(ctx, prompt)
    if err != nil {
        return fmt.Errorf("error translating content: %w", err)
    }

    // 翻訳されたテキストを取得
    if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
        return fmt.Errorf("no translation received")
    }
    translatedText := resp.Candidates[0].Content.Parts[0]

    fmt.Println("Translated text: ", translatedText) // 翻訳結果を確認するためにプリント

    em := dao.Client.EmbeddingModel("text-embedding-004")

    // Step 2: エンベディングを生成
	res, err := em.EmbedContent(ctx, translatedText)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}
    //embeddingベクトルを取得
	vec := res.Embedding.Values
	// 最初の5つの要素をプリント
    for i := 0; i < 5 && i < len(vec); i++ {
        fmt.Println(vec[i])
    }

    err = saveEmbeddingToGCS(ctx, vec, er.UserId)
    if err != nil {
        return fmt.Errorf("failed to save embedding to GCS: %w", err)
    }

    return nil
}

//検索ワードに類似したユーザーを返す関数
func (dao *VertexAiDAO) FindSimilar(ctx context.Context, fs model.FindSimilarRequest) ([]model.Profile,error){
	//Step 1:fs.contentのvectorを求める
    gemini := dao.Client.GenerativeModel("gemini-1.5-flash-002")
	fmt.Println(fs.SearchWord)
    prompt := genai.Text("Please return only the translated text."+"Translate the following tweet to English: " + fs.SearchWord )
    resp, err := gemini.GenerateContent(ctx, prompt)
    if err != nil {
        return nil,fmt.Errorf("error translating content: %w", err)
    }
	// 翻訳されたテキストを取得
    if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
        return nil,fmt.Errorf("no translation received")
    }
    translatedText := resp.Candidates[0].Content.Parts[0]

    fmt.Println("Translated text: ", translatedText) // 翻訳結果を確認するためにプリント
	em := dao.Client.EmbeddingModel("text-embedding-004")
	res, err := em.EmbedContent(ctx, translatedText)
	if err != nil {
		return nil,fmt.Errorf("failed to generate embedding: %w", err)
	}
    //embeddingベクトルを取得
	vec := res.Embedding.Values


	// Step 2: データベースからすべてのuser_idを取得
    var userIds []string
    rows, err := dao.DB.Query("SELECT user_id FROM user")
    if err != nil {
        return nil, fmt.Errorf("failed to fetch user IDs: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var userId string
        if err := rows.Scan(&userId); err != nil {
            return nil, fmt.Errorf("error scanning user ID: %w", err)
        }
        userIds = append(userIds, userId)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over user rows: %w", err)
    }

	//GCSにファイルが存在するかを確認し、存在すればコサイン類似度を計算
    var userScores []struct {
        userId string
        score  float32
    }

    for _, userId := range userIds {
        // Step 3.1: GCSからvectorを取得
        objectName := fmt.Sprintf("vector/%s.json", userId)
        embeddingData, err := getEmbeddingFromGCS(ctx, "term6-kyosuke-nishishita.firebasestorage.app", objectName)
        if err != nil {
            return nil, fmt.Errorf("error fetching user embedding from GCS for userId %s: %w", userId, err)
        }

        // Step 3.2: コサイン類似度を求める
        similarity, err := cosineSimilarity(embeddingData.Embedding, vec)
		fmt.Println("similarity",similarity)
        if err != nil {
            return nil, fmt.Errorf("error calculating cosine similarity for userId %s: %w", userId, err)
        }
		// Step 3.3: 類似度を保存
        userScores = append(userScores, struct {
            userId string
            score  float32
        }{
            userId: userId,
            score:  similarity,
        })
    }

	// **プリント追加**: userScoresの中身を確認
	fmt.Println("User Scores:")
	for _, score := range userScores {
		fmt.Printf("UserId: %s, Score: %f\n", score.userId, score.score)
	}

	// Step 4: 降順にソート
    sort.Slice(userScores, func(i, j int) bool {
        return userScores[i].score > userScores[j].score
    })
	
	// Step 5: 上位5人のプロフィールを取得
    var similarProfiles []model.Profile
    for i := 0; i < 3 && i < len(userScores); i++ {
        profile, err := dao.GetUserProfile(userScores[i].userId)
        if err != nil {
            return nil, fmt.Errorf("error fetching profile for userId %s: %w", userScores[i].userId, err)
        }
        similarProfiles = append(similarProfiles, profile)
    }

    return similarProfiles, nil
}

//userIdがフォローしていないと類似したvectorを持つ上位3件くらい
func (dao *VertexAiDAO) RecommendUser(ctx context.Context,userId string) ([]model.Profile, error) {
	// 1. userIdがフォローしていないユーザーを取得
	profiles, err := dao.GetNonFollowedUsers(userId)
	if err != nil {
		return nil, fmt.Errorf("could not get non-followed users: %w", err)
	}

	// 2. userIdの埋め込みベクトルを取得
	userEmbedding, err := dao.getUserEmbedding(userId)
	if err != nil {
		return nil, fmt.Errorf("could not get user embedding: %w", err)
	}

	// 3. フォローしていないユーザーの中で、最も類似度の高いユーザーを選出
	var similarities []struct {
		Profile   model.Profile
		Similarity float32
	}

	for _, profile := range profiles {
		// それぞれのユーザーの埋め込みベクトルを取得
		profileEmbedding, err := dao.getUserEmbedding(profile.Id)
		if err != nil {
			// 埋め込みが取得できない場合はスキップ
			continue
		}

		// コサイン類似度を計算
		similarity, err := cosineSimilarity(userEmbedding, profileEmbedding)
		if err != nil {
			// 類似度計算エラーが発生した場合はスキップ
			continue
		}

		// 類似度を保存
		similarities = append(similarities, struct {
			Profile   model.Profile
			Similarity float32
		}{
			Profile:   profile,
			Similarity: similarity,
		})
	}

	// 4. 類似度でソート
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	// 上位3件を選出
	top3Profiles := []model.Profile{}
	for i := 0; i < 3 && i < len(similarities); i++ {
		top3Profiles = append(top3Profiles, similarities[i].Profile)
	}

	// 5. 上位3件のユーザーを返す
	return top3Profiles, nil
}

// ユーザーの埋め込みベクトルを取得するメソッド
func (dao *VertexAiDAO) getUserEmbedding(userId string) ([]float32, error) {
	// GCSから埋め込みベクトルを取得する
	bucketName := "term6-kyosuke-nishishita.firebasestorage.app"
	objectName := fmt.Sprintf("vector/%s.json", userId)

	// GCSに保存されている埋め込みデータを取得
	embeddingResult, err := getEmbeddingFromGCS(context.Background(), bucketName, objectName)
	if err != nil {
		return nil, fmt.Errorf("failed to get embedding from GCS for userId %s: %w", userId, err)
	}

	// 埋め込みベクトルを返す
	return embeddingResult.Embedding, nil
}


// GetNonFollowedUsers メソッドを追加
func (dao *VertexAiDAO) GetNonFollowedUsers(userId string) ([]model.Profile, error) {
	// フォローしていないユーザーを取得するSQLクエリ
	query := `
		SELECT u.user_id, u.name, u.bio, u.profile_img_url, u.header_img_url, u.location
		FROM user u
		WHERE u.user_id != ?  -- 自分自身を除外
		AND NOT EXISTS (
			SELECT 1
			FROM follower f
			WHERE f.user_id = ? AND f.following_user_id = u.user_id
		)
	`

	// フォローしていないユーザーのプロフィールを格納するスライス
	var profiles []model.Profile

	// データベースから情報を取得
	rows, err := dao.DB.Query(query, userId, userId)
	if err != nil {
		log.Printf("Error fetching non-followed users for userId %s: %v", userId, err)
		return nil, fmt.Errorf("could not fetch non-followed users: %w", err)
	}
	defer rows.Close()

	// 取得した各行を処理
	for rows.Next() {
		var profile model.Profile
		if err := rows.Scan(&profile.Id, &profile.Name, &profile.Bio, &profile.ImgUrl, &profile.HeaderUrl, &profile.Location); err != nil {
			log.Printf("Error scanning profile for userId %s: %v", userId, err)
			continue
		}
		// プロフィールをスライスに追加
		profiles = append(profiles, profile)
	}

	// エラーがあれば返す
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	// フォローしていないユーザーのプロフィールを返す
	return profiles, nil
}






func (dao *VertexAiDAO) GetUserProfile(userId string) (model.Profile, error) {
	var prof model.Profile
	err := dao.DB.QueryRow("SELECT user_id, name, bio, profile_img_url,header_img_url,location FROM user WHERE user_id = ?", userId).Scan(&prof.Id, &prof.Name, &prof.Bio,&prof.ImgUrl,&prof.HeaderUrl,&prof.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			// ユーザーが見つからなかった場合
			return model.Profile{}, nil  // 空の構造体を返す
		}

		// その他のエラー
		log.Printf("Error fetching user profile for userId %s: %v", userId, err)
		return model.Profile{}, fmt.Errorf("could not fetch user profile: %w", err)  // ラップしたエラーを返す
	}

	return prof, nil
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
	fmt.Println(objectName)

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