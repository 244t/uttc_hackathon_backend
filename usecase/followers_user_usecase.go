package usecase

import(
	"myproject/dao"
	"myproject/model"
)

type GetFollowersUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

func NewGetFollowersUserUseCase(tweetDao dao.TweetDAOInterface) *GetFollowersUserUseCase{
	return &GetFollowersUserUseCase{TweetDAO:tweetDao}
}

func (uc *GetFollowersUserUseCase) GetFollowers(userId string)([]model.Profile,error){
	return uc.TweetDAO.GetFollowers(userId)
}