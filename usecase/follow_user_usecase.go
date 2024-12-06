package usecase

import(
	"myproject/dao"
	"myproject/model"
	"log"
	"github.com/oklog/ulid"
	"time"
	"math/rand"
)

type FollowRegister struct{
	UserId string `json:"user_id"`
	FollowingId string `json:"following_id"`
}

type FollowUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

//FollowUserUseCaseのファクトリ関数
func NewFollowUserUseCase(tweetDao dao.TweetDAOInterface) *FollowUserUseCase{
	return &FollowUserUseCase{
		TweetDAO: tweetDao,
	}
}

func (uc *FollowUserUseCase) Follow(followRegister FollowRegister) error{
	// ULIDの生成
	entropy := rand.New(rand.NewSource(time.Now().UnixNano())) // 乱数生成器の作成
	ulid := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)   // ULID

	follow := model.Follow{
		UserId: followRegister.UserId,
		FollowingId: followRegister.FollowingId,
	}

	notification := model.Notification{
		NotificationId : ulid.String(),
		UserId: follow.FollowingId,
		Flag : "follow",
		ActionUserId : follow.UserId,
	}

	if err := uc.TweetDAO.Follow(follow,notification); err != nil {
		log.Printf("failed to save follow: %v", err)
		return err
	}

	return nil
}