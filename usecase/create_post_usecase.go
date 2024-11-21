package usecase

import (
	"myproject/dao"
	"myproject/model"
	"log"
	"github.com/oklog/ulid"
	"time"
	"database/sql"
	"math/rand"
)

type PostRegister struct {
	UserId string `json:"user_id"`
	Content string `json:"content"`
	ImgUrl string `json:"img_url"`
}

type CreatePostUseCase struct{
	PostDAO dao.PostDAOInterface
}

//CreatePostUseCaseのファクトリ関数
func NewCreatePostUseCase(postDao dao.PostDAOInterface) *CreatePostUseCase{
	return &CreatePostUseCase{
		PostDAO: postDao,
	}
}

func (uc *CreatePostUseCase) CreatePost(postRegister PostRegister) error {

	// ULIDの生成
	entropy := rand.New(rand.NewSource(time.Now().UnixNano())) // 乱数生成器の作成
	ulid := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)   // ULID
	
	//新しいポストを作成
	newPost := model.Post{
		UserId : postRegister.UserId,
		PostId : ulid.String(),
		Content : postRegister.Content,
		ImgUrl : postRegister.ImgUrl,
		CreatedAt : time.Now(),
		EditedAt : sql.NullTime{Valid: false},  
		DeletedAt : sql.NullTime{Valid: false},  
		ParentPostId : sql.NullString{Valid: false}, 
	}

	if err := uc.PostDAO.CreatePost(newPost); err != nil {
		log.Printf("failed to save post: %v", err)
		return err
	}

	return nil
}