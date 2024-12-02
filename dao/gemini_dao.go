// package dao

// import (
// 	"context"
// 	"fmt"

// 	"cloud.google.com/go/vertexai/genai"
// 	"google.golang.org/api/option"
// 	"myproject/model"
// )

// const (
//     location  = "asia-northeast1"
//     modelName = "gemini-1.5-flash-002"
//     projectID = "term6-kyosuke-nishishita"
// )


// type VertexAiDAO struct{
// 	Client *genai.Client
// }

// type VertexAiDAOInterface interface{
// 	NextTextGeneration(ctx context.Context,text string) (*model.TextGenerationResponse,error)
// }

// func NewVertexAiDAO (client *vertexai.PredictionClient) *VertexAiDAO,error{
// 	client, err := genai.NewClient(ctx, projectID, location, option.WithoutAuthentication())
// 	if err != nil {
//         return nil, fmt.Errorf("failed to create client: %w", err)
//     }
// 	return &VertexAiDAO{Client: client},nil
// }


// // GenerateText はGeminiモデルを使用してテキスト生成を行います。
// func (dao *VertexAiDAO) NextTextGeneration(ctx context.Context, promptText string) (string, error) {
//     gemini := dao.Client.GenerativeModel(modelName)
//     prompt := genai.Text(promptText)
//     resp, err := gemini.GenerateContent(ctx, prompt)
//     if err != nil {
//         return "", fmt.Errorf("error generating content: %w", err)
//     }

//     // レスポンスから生成されたテキストを取得
//     if len(resp.GeneratedContent) == 0 {
//         return "", fmt.Errorf("no generated content received")
//     }

//     return resp.GeneratedContent[0], nil
// }

// // Close クライアントを閉じる
// func (dao *GeminiDAO) Close() error {
//     if dao.client != nil {
//         return dao.client.Close()
//     }
//     return nil
// }

package dao

import (
	"context"
	"fmt"
	"cloud.google.com/go/vertexai/genai"
)

const (
	location  = "asia-northeast1"
	modelName = "gemini-1.5-flash-002"
	projectID = "your-project-id"
)

type VertexAiDAO struct {
	Client *genai.Client
}

type VertexAiDAOInterface interface {
	NextTextGeneration(ctx context.Context, text string) (*genai.Part, error)
}

func NewVertexAiDAO(client *genai.Client) *VertexAiDAO {
	return &VertexAiDAO{
		Client: client,
	}
}

func (dao *VertexAiDAO) NextTextGeneration(ctx context.Context, promptText string) (*genai.Part, error) {
    gemini := dao.Client.GenerativeModel(modelName)
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

