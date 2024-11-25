package usecase

import (
	"myproject/dao"
)

type LikesResponse struct {
    LikeCount int      `json:"like_count"`
    UserIds   []string `json:"user_ids"`
}

type GetLikesPostUseCase struct {
	PostDAO dao.PostDAOInterface
}

func NewGetLikesPostUseCase(pd dao.PostDAOInterface) *GetLikesPostUseCase{
	return &GetLikesPostUseCase{PostDAO: pd}
}

func (gp *GetLikesPostUseCase) GetLikesPost(postId string)([]string,error){
	return gp.PostDAO.GetLikesPost(postId)
}