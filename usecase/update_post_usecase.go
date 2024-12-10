package usecase

import (
	"myproject/dao"
	"myproject/model"
	"log"
	"database/sql"
	"time"
)

type PostUpdate struct {
	PostId string `json:"post_id"`
	Content string `json:"content"`
	ImgUrl string `json:"img_url"`
}

type UpdatePostUseCase struct{
	PostDAO dao.PostDAOInterface
}

func NewUpdatePostUSeCase(postDao dao.PostDAOInterface) *UpdatePostUseCase{
	return &UpdatePostUseCase{
		PostDAO: postDao,
	}
}

func (uc *UpdatePostUseCase) UpdatePost(pu PostUpdate) error {

	updatePost := model.Update{
		PostId : pu.PostId,
		Content: pu.Content,
		ImgUrl : pu.ImgUrl,
		EditedAt : sql.NullTime{Time: time.Now(), Valid: true},
	}

	if err := uc.PostDAO.UpdatePost(updatePost); err != nil {
		log.Printf("failed to update post: %v", err)
		return err
	}

	return nil
}