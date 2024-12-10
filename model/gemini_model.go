package model

// TextGenerationRequest は、クライアントから受け取るリクエストの構造体です。
type TextGenerationRequest struct {
	Text string `json:"text"`
}

// TextGenerationResponse は、Vertex AIからの予測結果を格納する構造体です。
type TextGenerationResponse struct {
	SuggestedText string `json:"suggestedText"`
}

type EmbeddingRequest struct{
	UserId string `json:"user_id"`
	Content string `json:"content"`
}

// 埋め込み結果の構造体
type EmbeddingResult struct {
	UserID   string    `json:"user_id"`
	Count    int       `json:"count"`
	Embedding []float32 `json:"embedding"`
}

type FindSimilarRequest struct{
	SearchWord   string       `json:"search_word"`
}