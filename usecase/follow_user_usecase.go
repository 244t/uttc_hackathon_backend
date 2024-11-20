package usecase

import(
	"myproject/dao"
	"myproject/model"
	"log"
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
	follow := model.Follow{
		UserId: followRegister.UserId,
		FollowingId: followRegister.FollowingId,
	}

	if err := uc.TweetDAO.Follow(follow); err != nil {
		log.Printf("failed to save follow: %v", err)
		return err
	}

	return nil
}