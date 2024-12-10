package usecase

import (
	"myproject/dao"
	"myproject/model"
)

type NotificationUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

func NewNotificationUserUseCase(tweetDao dao.TweetDAOInterface)*NotificationUserUseCase{
	return &NotificationUserUseCase{TweetDAO:tweetDao}
}

func (uc *NotificationUserUseCase) Notification(userId string) ([]model.NotificationInfo,error){
	return uc.TweetDAO.Notification(userId)
}