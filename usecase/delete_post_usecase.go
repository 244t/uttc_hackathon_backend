package usecase

import (
	"myproject/dao"
	"myproject/model"
	"log"
	"database/sql"
	"time"
)

type DeletePostUseCase struct{
	PostDAO dao.PostDAOInterface
}

func NewDeletePostUseCase(postDao dao.PostDAOInterface) *DeletePostUseCase{
	return &DeletePostUseCase{
		PostDAO: postDao,
	}
}

func (uc *DeletePostUseCase) DeletePost(postId string) error {
	
	deletePost := model.Delete{
		PostId : postId,
		DeletedAt : sql.NullTime{Time: time.Now(), Valid: true},
	}

	if err := uc.PostDAO.DeletePost(deletePost); err != nil {
		log.Printf("failed to delete post: %v", err)
		return err
	}

	return nil
}