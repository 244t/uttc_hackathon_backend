package usecase

import (
	"errors"
	"myproject/dao"
	"myproject/model"
	"context"
	"log"
)

type RegisterUserUseCase struct {
	VertexAiDAO dao.VertexAiDAOInterface
}

//RegisterUserUseCaseのファクトリ関数
func NewRegisterUserUseCase(v dao.VertexAiDAOInterface) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		VertexAiDAO: v,
	}
}

func (uc *RegisterUserUseCase) RegisterUser(ctx context.Context,ur model.Profile) error{
	// 入力のバリデーション
	if ur.Name == "" || len(ur.Name) > 50  {
		return errors.New("validation failed: invalid user data")
	}

	// ユーザーを保存
	if err := uc.VertexAiDAO.RegisterUser(ctx,ur); err != nil {
		log.Printf("failed to save user: %v", err)
		return err
	}

	return nil
}