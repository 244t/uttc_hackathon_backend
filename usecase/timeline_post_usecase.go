package usecase

import (
	"myproject/dao"
	"myproject/model"
)

type TimelineUseCase struct {
	PostDAO dao.PostDAOInterface
}

func NewTimelineUseCase(pd dao.PostDAOInterface) *TimelineUseCase{
	return &TimelineUseCase{PostDAO: pd}
}

// func (gp *TimelineUseCase) Timeline(userId string)([]model.PostWithReplyCounts,error){
// 	return gp.PostDAO.Timeline(userId)
// }
func (gp *TimelineUseCase) Timeline(userId string, limit, page int) ([]model.PostWithReplyCounts, int, error) {
    return gp.PostDAO.Timeline(userId, limit, page)
}
