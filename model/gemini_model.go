// package model

// // TextGenerationRequest は、クライアントから受け取るリクエストの構造体です。
// type TextGenerationRequest struct {
//     Text string `json:"text"`
// }

// // TextGenerationResponse は、Vertex AIからの予測結果を格納する構造体です。
// type TextGenerationResponse struct {
//     SuggestedText string `json:"suggestedText"`
// }

package model

// TextGenerationRequest は、クライアントから受け取るリクエストの構造体です。
type TextGenerationRequest struct {
	Text string `json:"text"`
}

// TextGenerationResponse は、Vertex AIからの予測結果を格納する構造体です。
type TextGenerationResponse struct {
	SuggestedText string `json:"suggestedText"`
}
