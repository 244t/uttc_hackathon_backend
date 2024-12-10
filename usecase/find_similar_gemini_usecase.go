package usecase

import (
	"context"
	"myproject/dao"
	"myproject/model"
)

type FindSimilarUseCase struct{
	VertexAiDAO dao.VertexAiDAOInterface
}

func NewFindSimilarUseCase(v dao.VertexAiDAOInterface) *FindSimilarUseCase{
	return &FindSimilarUseCase{
		VertexAiDAO: v,
	}
}

func (uc *FindSimilarUseCase) FindSimilar(ctx context.Context, fs model.FindSimilarRequest)([]model.Profile,error){
	return uc.VertexAiDAO.FindSimilar(ctx,fs)
}