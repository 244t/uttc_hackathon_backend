package usecase

import(
	"myproject/dao"
	"myproject/model"
	"log"
)

type UnFollowRegister struct{
	UserId string `json:"user_id"`
	FollowingId string `json:"following_id"`
}

type UnFollowUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

//FollowUserUseCaseのファクトリ関数
func NewUnFollowUserUseCase(tweetDao dao.TweetDAOInterface) *UnFollowUserUseCase{
	return &UnFollowUserUseCase{
		TweetDAO: tweetDao,
	}
}

func (uc *UnFollowUserUseCase) UnFollow(unFollowRegister UnFollowRegister) error{
	unfollow := model.UnFollow{
		UserId: unFollowRegister.UserId,
		FollowingId: unFollowRegister.FollowingId,
	}

	if err := uc.TweetDAO.UnFollow(unfollow); err != nil {
		log.Printf("failed to save unfollow: %v", err)
		return err
	}

	return nil
}