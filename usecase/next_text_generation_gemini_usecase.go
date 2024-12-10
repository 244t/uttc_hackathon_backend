package usecase

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"myproject/dao"
)

type NextTextGenerationUseCase struct {
	VertexAiDAO dao.VertexAiDAOInterface
}

func NewNextTextGenerationUseCase(v dao.VertexAiDAOInterface) *NextTextGenerationUseCase {
	return &NextTextGenerationUseCase{
		VertexAiDAO: v,
	}
}

// NextTextGeneration は、指定されたテキストをもとに新しいテキストを生成します。
func (uc *NextTextGenerationUseCase) NextTextGeneration(ctx context.Context, text string) (*genai.Part, error) {
	return uc.VertexAiDAO.NextTextGeneration(ctx, text)
}
