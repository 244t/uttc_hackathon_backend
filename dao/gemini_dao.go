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

