package controllers

import (
	"net/http"
	"myproject/dao"
	"myproject/model"
	"myproject/usecase"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
)

// GeminiController は、Geminiテキスト生成に関連する処理を行うコントローラーです。
type GeminiController struct {
	RegisterUserUseCase *usecase.RegisterUserUseCase
	NextTextGenerationUseCase *usecase.NextTextGenerationUseCase
	EmbeddingGenerationUseCase *usecase.EmbeddingGenerationUseCase
	FindSimilarUseCase *usecase.FindSimilarUseCase
	RecommendUserUseCase *usecase.RecommendUserUseCase
}

// NewGeminiController は、GeminiControllerの新しいインスタンスを作成します。
func NewGeminiController(db dao.VertexAiDAOInterface) *GeminiController {
	registerUserUseCase := usecase.NewRegisterUserUseCase(db)
	nextTextGenerationUseCase := usecase.NewNextTextGenerationUseCase(db)
	embeddingGenerationUseCase := usecase.NewEmbeddingGenerationUseCase(db)
	findSimilarUseCase := usecase.NewFindSimilarUseCase(db)
	recommendUserUseCase := usecase.NewRecommendUserUseCase(db)
	return &GeminiController{
		RegisterUserUseCase: registerUserUseCase,
		NextTextGenerationUseCase: nextTextGenerationUseCase,
		EmbeddingGenerationUseCase : embeddingGenerationUseCase,
		FindSimilarUseCase : findSimilarUseCase,
		RecommendUserUseCase : recommendUserUseCase,
	}
}
// ユーザー登録
func (c *GeminiController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRegister model.Profile
	if err := json.NewDecoder(r.Body).Decode(&userRegister); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	if err := c.RegisterUserUseCase.RegisterUser(ctx,userRegister); err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

func (c *GeminiController) FindSimilar(w http.ResponseWriter, r *http.Request){
	var request model.FindSimilarRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	users, err := c.FindSimilarUseCase.FindSimilar(ctx, request)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to Find Similar: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func (c *GeminiController) RecommendUser(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	userId := vars["userId"]
	ctx := r.Context()
	users, err := c.RecommendUserUseCase.RecommendUser(ctx, userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to recommend users: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
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