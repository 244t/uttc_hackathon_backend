package usecase

import (
	"errors"
	"myproject/dao"
	"myproject/model"
	"log"
	"github.com/oklog/ulid"
	"time"
	"math/rand"
)

type UserRegister struct {
	Name string `json:"name"`
	Bio string `json:"bio"`
	FireBaseId string `json:"firebase_id"`
}

type RegisterUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

//RegisterUserUseCaseのファクトリ関数
func NewRegisterUserUseCase(tweetDao dao.TweetDAOInterface) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		TweetDAO: tweetDao,
	}
}

func (uc *RegisterUserUseCase) RegisterUser(userRegister UserRegister) error{
	// 入力のバリデーション
	if userRegister.Name == "" || len(userRegister.Name) > 50  {
		return errors.New("validation failed: invalid user data")
	}

	// ULIDの生成
	entropy := rand.New(rand.NewSource(time.Now().UnixNano())) // 乱数生成器の作成
	ulid := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)   // ULID

	//新しいユーザーを作成
	user := model.Profile{
		Id : ulid.String(),
		Name: userRegister.Name,
		Bio: userRegister.Bio,
		FireBaseId: userRegister.FireBaseId,
	}

	// ユーザーを保存
	if err := uc.TweetDAO.RegisterUser(user); err != nil {
		log.Printf("failed to save user: %v", err)
		return err
	}

	return nil
}