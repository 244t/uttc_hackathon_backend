package usecase

import (
	"context"
	"myproject/dao"
	"myproject/model"
)

type EmbeddingGenerationUseCase struct{
	VertexAiDAO dao.VertexAiDAOInterface
}

func NewEmbeddingGenerationUseCase(v dao.VertexAiDAOInterface) *EmbeddingGenerationUseCase {
	return &EmbeddingGenerationUseCase{
		VertexAiDAO: v,
	}
}

func (uc *EmbeddingGenerationUseCase) EmbeddingGeneration(ctx context.Context,er model.EmbeddingRequest) error{
	return uc.VertexAiDAO.EmbeddingGeneration(ctx,er)
}