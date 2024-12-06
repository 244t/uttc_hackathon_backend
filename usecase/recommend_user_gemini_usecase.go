package usecase

import (
	"context"
	"myproject/dao"
	"myproject/model"
)

type RecommendUserUseCase struct{
	VertexAiDAO dao.VertexAiDAOInterface
}

func NewRecommendUserUseCase(v dao.VertexAiDAOInterface) *RecommendUserUseCase{
	return &RecommendUserUseCase{
		VertexAiDAO : v,
	}
}

func (uc *RecommendUserUseCase) RecommendUser(ctx context.Context, userId string)([]model.Profile, error){
	return uc.VertexAiDAO.RecommendUser(ctx, userId)
}