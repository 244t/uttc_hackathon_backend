// package controllers

// import (
// 	"net/http"
// 	"myproject/usecase"
// 	"myproject/dao"
// 	"encoding/json"
// 	"github.com/gorilla/mux"
// )

// type GeminiController struct{
// 	NextTextGenerationUseCase *usecase.NextTextGenerationUseCase
// }

// func NewGeminiController(db dao.VertexAiDAO) *GeminiController{
// 	nextTextGenerationUseCase := usecase.NewNextTextGenerationUseCase(db)

// 	return &GeminiController{
// 		NextTextGenerationUseCase : nextTextGenerationUseCase,
// 	}
// }

// func (c *GeminiController) NextTextGeneration(w http.ResponseWriter, r* http.Request){
// 	var request model.TextGenerationRequest
// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}
// 	ctx := r.Context()
// 	response, err := c.NextTextGenerationUseCase.NextTextGeneration(ctx,request.Text)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to generate text: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// レスポンスをクライアントに返す
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

package controllers

import (
	"net/http"
	"myproject/dao"
	"myproject/model"
	"myproject/usecase"
	"encoding/json"
	"fmt"
)

// GeminiController は、Geminiテキスト生成に関連する処理を行うコントローラーです。
type GeminiController struct {
	NextTextGenerationUseCase *usecase.NextTextGenerationUseCase
}

// NewGeminiController は、GeminiControllerの新しいインスタンスを作成します。
func NewGeminiController(db dao.VertexAiDAOInterface) *GeminiController {
	nextTextGenerationUseCase := usecase.NewNextTextGenerationUseCase(db)
	return &GeminiController{
		NextTextGenerationUseCase: nextTextGenerationUseCase,
	}
}


func (c *GeminiController) NextTextGeneration(w http.ResponseWriter, r *http.Request){
	var request model.TextGenerationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	part, err := c.NextTextGenerationUseCase.NextTextGeneration(ctx, request.Text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate text: %v", err), http.StatusInternalServerError)
		return
	}

	// Part をそのままクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(part); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

// OPTIONSリクエストに対する処理
func (c *GeminiController) CORSOptionsHandler(w http.ResponseWriter, r *http.Request) {
    // 必要なCORSヘッダーを設定
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT,POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.WriteHeader(http.StatusOK) // 200 OKを返す
}