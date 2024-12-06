package usecase

import (
	"myproject/dao"
	"myproject/model"
	"log"
	"github.com/oklog/ulid"
	"time"
	"math/rand"
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
	// ULIDの生成
	entropy := rand.New(rand.NewSource(time.Now().UnixNano())) // 乱数生成器の作成
	ulid := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)   // ULID

	like := model.Like{
		UserId : l.UserId,
		PostId : l.PostId,
		CreatedAt : time.Now(),
	}

	notification := model.Notification{
		NotificationId : ulid.String(),
		UserId: "",
		Flag : "like",
		ActionUserId : "l.UserId",
	}

	if err := uc.PostDAO.LikePost(like,notification); err != nil {
		log.Printf("failed to save like: %v", err)
		return err
	}

	return nil
}