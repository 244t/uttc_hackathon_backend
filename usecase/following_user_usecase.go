package usecase

import(
	"myproject/dao"
	"myproject/model"
)

type GetFollowingUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

func NewGetFollowingUserUseCase(tweetDao dao.TweetDAOInterface) *GetFollowingUserUseCase{
	return &GetFollowingUserUseCase{TweetDAO:tweetDao}
}

func (uc *GetFollowingUserUseCase) GetFollowing(userId string)([]model.Profile,error){
	return uc.TweetDAO.GetFollowing(userId)
}