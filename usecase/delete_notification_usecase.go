package usecase

import (
	"myproject/dao"
	"log"
)

type DeleteNotificationUserUseCase struct{
	TweetDAO dao.TweetDAOInterface
}

func NewDeleteNotificationUserUseCase(tweetDao dao.TweetDAOInterface) *DeleteNotificationUserUseCase{
	return &DeleteNotificationUserUseCase{
		TweetDAO: tweetDao,
	}
}

func (uc *DeleteNotificationUserUseCase) DeleteNotification(ni string) error{
	if err := uc.TweetDAO.DeleteNotification(ni); err != nil {
		log.Printf("failed to delete notification: %v", err)
		return err
	}
	return nil
}