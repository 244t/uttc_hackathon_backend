package usecase

import (
	"myproject/dao"
	"myproject/model"
)

type Search struct {
	SearchWord string `json:"search_word"`
}

type SearchUserUseCase struct{
	TweetDAO dao.TweetDAOInterface
}

//CreatePostUseCaseのファクトリ関数
func NewSearchUserUseCase(tweetDao dao.TweetDAOInterface) *SearchUserUseCase{
	return &SearchUserUseCase{
		TweetDAO: tweetDao,
	}
}

func (uc *SearchUserUseCase) SearchUser(sw Search) ([]model.Profile,error) {
	return uc.TweetDAO.SearchUser(sw.SearchWord);
}