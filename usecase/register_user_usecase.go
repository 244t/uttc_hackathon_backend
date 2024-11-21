package usecase

import (
	"errors"
	"myproject/dao"
	"myproject/model"
	"log"
)

type RegisterUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

//RegisterUserUseCaseのファクトリ関数
func NewRegisterUserUseCase(tweetDao dao.TweetDAOInterface) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		TweetDAO: tweetDao,
	}
}

func (uc *RegisterUserUseCase) RegisterUser(ur model.Profile) error{
	// 入力のバリデーション
	if ur.Name == "" || len(ur.Name) > 50  {
		return errors.New("validation failed: invalid user data")
	}

	// ユーザーを保存
	if err := uc.TweetDAO.RegisterUser(ur); err != nil {
		log.Printf("failed to save user: %v", err)
		return err
	}

	return nil
}