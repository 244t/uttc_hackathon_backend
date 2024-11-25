package usecase

import (
	"myproject/dao"
	"myproject/model"
)

type GetReplyPostUseCase struct {
	PostDAO dao.PostDAOInterface
}

func NewGetReplyPostUseCase(pd dao.PostDAOInterface) *GetReplyPostUseCase{
	return &GetReplyPostUseCase{PostDAO: pd}
}

func (gp *GetReplyPostUseCase) GetReplyPost(postId string)([]model.PostWithReplyCounts,error){
	return gp.PostDAO.GetReplyPost(postId)
}