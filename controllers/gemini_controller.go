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
	EmbeddingGenerationUseCase *usecase.EmbeddingGenerationUseCase
}

// NewGeminiController は、GeminiControllerの新しいインスタンスを作成します。
func NewGeminiController(db dao.VertexAiDAOInterface) *GeminiController {
	nextTextGenerationUseCase := usecase.NewNextTextGenerationUseCase(db)
	embeddingGenerationUseCase := usecase.NewEmbeddingGenerationUseCase(db)
	return &GeminiController{
		NextTextGenerationUseCase: nextTextGenerationUseCase,
		EmbeddingGenerationUseCase : embeddingGenerationUseCase,
	}
}

// Tweetの続きを生成
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

// Embeddingを生成
func (c *GeminiController) EmbeddingGeneration(w http.ResponseWriter, r *http.Request){
	var request model.EmbeddingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err := c.EmbeddingGenerationUseCase.EmbeddingGeneration(ctx,request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed embedding: %v", err), http.StatusInternalServerError)
		return
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