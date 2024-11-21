package usecase

import (
	"myproject/dao"
	"myproject/model"
	"log"
	"time"
)

type Like struct {
	UserId string `json:"user_id"`
	PostId string `json:"post_id"`
}

type LikePostUseCase struct{
	PostDAO dao.PostDAOInterface
}

//LikePostUseCaseのファクトリ関数
func NewLikePostUseCase(postDao dao.PostDAOInterface) *LikePostUseCase{
	return &LikePostUseCase{
		PostDAO: postDao,
	}
}

func (uc *LikePostUseCase) LikePost(l Like) error {

	like := model.Like{
		UserId : l.UserId,
		PostId : l.PostId,
		CreatedAt : time.Now(),
	}

	if err := uc.PostDAO.LikePost(like); err != nil {
		log.Printf("failed to save like: %v", err)
		return err
	}

	return nil
}