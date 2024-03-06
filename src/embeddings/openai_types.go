package embeddings

import "github.com/jake-landersweb/gollm/src/ltypes"

// https://platform.openai.com/docs/api-reference/embeddings/create

type OpenAIEmbeddingRequest struct {
	Input          []string `json:"input" binding:"required"` // Can be string or []string
	Model          string   `json:"model" binding:"required"`
	EncodingFormat string   `json:"encoding_format,omitempty"` // Defaults to "float", can be "float" or "base64"
	Dimensions     int      `json:"dimensions,omitempty"`      // Optional, supported in text-embedding-3 and later models
	User           string   `json:"user,omitempty"`            // Optional unique identifier for end-user
}

type OpenAIEmbeddingResponse struct {
	Object string                `json:"object"`
	Data   []OpenAIEmbeddingData `json:"data"`
	Model  string                `json:"model"`
	Usage  ltypes.GPTUsage       `json:"usage"`
	Error  *ltypes.GPTError      `json:"error"`
}

type OpenAIEmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}
