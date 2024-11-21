package usecase

import (
	"myproject/dao"
	"myproject/model"
)

type GetLikesPostUseCase struct {
	PostDAO dao.PostDAOInterface
}

func NewGetLikesPostUseCase(pd dao.PostDAOInterface) *GetLikesPostUseCase{
	return &GetLikesPostUseCase{PostDAO: pd}
}

func (gp *GetLikesPostUseCase) GetLikesPost(postId string)(model.Likes,error){
	return gp.PostDAO.GetLikesPost(postId)
}