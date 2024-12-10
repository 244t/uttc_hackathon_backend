package usecase

import (
	"myproject/dao"
	"myproject/model"
)

type GetPostUseCase struct {
	PostDAO dao.PostDAOInterface
}

func NewGetPostUseCase(pd dao.PostDAOInterface) *GetPostUseCase{
	return &GetPostUseCase{PostDAO: pd}
}

func (gp *GetPostUseCase) GetPost(postId string)(model.PostWithReplyCounts,error){
	return gp.PostDAO.GetPost(postId)
}