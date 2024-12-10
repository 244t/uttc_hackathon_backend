package usecase

import (
	"myproject/dao"
	"myproject/model"
)

type GetProfileUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

func NewGetProfileUserUseCase(tweetDao dao.TweetDAOInterface) *GetProfileUserUseCase{
	return &GetProfileUserUseCase{TweetDAO:tweetDao}
}

func (uc *GetProfileUserUseCase) GetUserProfile(userId string) (model.Profile,error){
	return uc.TweetDAO.GetUserProfile(userId)
}