package usecase

import (
	"errors"
	"myproject/dao"
	"myproject/model"
	"log"
)


type UpdateProfileUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

//UpdateProfileUserUseCaseのファクトリ関数
func NewUpdateProfileUserUseCase(tweetDao dao.TweetDAOInterface) *UpdateProfileUserUseCase {
	return &UpdateProfileUserUseCase{
		TweetDAO: tweetDao,
	}
}

func (uc *UpdateProfileUserUseCase) UpdateProfile(up model.Profile) error{
	// 入力のバリデーション
	if up.Name == "" || len(up.Name) > 50  {
		return errors.New("validation failed: invalid user data")
	}

	// ユーザーを保存
	if err := uc.TweetDAO.UpdateProfile(up); err != nil {
		log.Printf("failed to save user: %v", err)
		return err
	}

	return nil
}