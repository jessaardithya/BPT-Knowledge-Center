package models

import "time"

type Document struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Filename     string    `json:"filename"`
	FileURL      string    `json:"file_url"` // MinIO URL
	ContentType  string    `json:"content_type"`
	UploadedAt   time.Time `json:"uploaded_at"`
	ElementCount int       `json:"element_count"`
	DocType      string    `json:"_type"` // Required for Search Index mapping

	Chunks []DocumentChunk `json:"chunks"`
}

type DocumentChunk struct {
	ChunkID  string                 `json:"chunk_id"`
	Text     string                 `json:"text"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
	Vector   []float32              `json:"vector"` // 🆕 Added Vector Field
}
