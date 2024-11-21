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

type ReplyPostUseCase struct{
	PostDAO dao.PostDAOInterface
}

func NewReplyPostUseCase(pd dao.PostDAOInterface) *ReplyPostUseCase{
	return &ReplyPostUseCase{
		PostDAO: pd,
	}
}

func (uc *ReplyPostUseCase) ReplyPost(parentId string, pr PostRegister) error{

	// ULIDの生成
	entropy := rand.New(rand.NewSource(time.Now().UnixNano())) // 乱数生成器の作成
	ulid := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)   // ULID
	
	//返信を作成
	reply := model.Post{
		UserId : pr.UserId,
		PostId : ulid.String(),
		Content : pr.Content,
		ImgUrl : pr.ImgUrl,
		CreatedAt : time.Now(),
		EditedAt : sql.NullTime{Valid: false},  
		DeletedAt : sql.NullTime{Valid: false},  
		ParentPostId : sql.NullString{String : parentId, Valid: true},
	}

	if err := uc.PostDAO.CreatePost(reply); err != nil {
		log.Printf("failed to save reply: %v", err)
		return err
	}

	return nil
}