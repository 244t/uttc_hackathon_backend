// package usecase

// import (
// 	"context"
// 	"myproject/dao"
// 	"myproject/model"
// )

// type NextTextGenerationUseCase struct{
// 	VertexAiDAO dao.VertexAiDAOInterface
// }

// func NewNextTextGenerationUseCase(v dao.VertexAiDAOInterface) *NextTextGeneration{
// 	return &NextTextGeneration{
// 		VertexAiDAO: v,
// 	}
// }

// func (uc *NextTextGenerationUseCase) NextTextGeneration(ctx context.Context,text string){
// 	return usecase.VertexAiDAO.NextTextGeneration(ctx,text)
// }

package usecase

import (
	"context"
	"cloud.google.com/go/vertexai/genai"
	"myproject/dao"
)

type NextTextGenerationUseCase struct {
	VertexAiDAO dao.VertexAiDAOInterface
}

func NewNextTextGenerationUseCase(v dao.VertexAiDAOInterface) *NextTextGenerationUseCase {
	return &NextTextGenerationUseCase{
		VertexAiDAO: v,
	}
}

// NextTextGeneration は、指定されたテキストをもとに新しいテキストを生成します。
func (uc *NextTextGenerationUseCase) NextTextGeneration(ctx context.Context, text string) (*genai.Part, error) {
	return uc.VertexAiDAO.NextTextGeneration(ctx, text)
}
