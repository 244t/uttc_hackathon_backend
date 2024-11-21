package usecase

import(
	"myproject/dao"
	"myproject/model"
)

type GetUserPostsUserUseCase struct {
	TweetDAO dao.TweetDAOInterface
}

//FollowUserUseCaseのファクトリ関数
func NewGetUserPostsUserUseCase(tweetDao dao.TweetDAOInterface) *GetUserPostsUserUseCase{
	return &GetUserPostsUserUseCase{
		TweetDAO: tweetDao,
	}
}


func (gp *GetUserPostsUserUseCase) GetUserPosts(userId string)([]model.Post,error){
	return gp.TweetDAO.GetUserPosts(userId)
}