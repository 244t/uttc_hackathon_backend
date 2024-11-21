package usecase

import (
	"myproject/dao"
	"myproject/model"
	"log"
	"time"
)



type UnLikePostUseCase struct{
	PostDAO dao.PostDAOInterface
}

//LikePostUseCaseのファクトリ関数
func NewUnLikePostUseCase(postDao dao.PostDAOInterface) *UnLikePostUseCase{
	return &UnLikePostUseCase{
		PostDAO: postDao,
	}
}

func (uc *UnLikePostUseCase) UnLikePost(l Like) error {

	like := model.Like{
		UserId : l.UserId,
		PostId : l.PostId,
		CreatedAt : time.Now(),
	}

	if err := uc.PostDAO.UnLikePost(like); err != nil {
		log.Printf("failed to save like: %v", err)
		return err
	}

	return nil
}