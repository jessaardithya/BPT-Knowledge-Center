package models

import "time"

type Document struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Filename     string          `json:"filename"`
	DisplayName  string          `json:"display_name"`
	FileURL      string          `json:"file_url"`
	ContentType  string          `json:"content_type"`
	UploadedAt   time.Time       `json:"uploaded_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	ElementCount int             `json:"element_count"`
	Version      int             `json:"version"`
	DocType      string          `json:"_type"`
	Category     string          `json:"category"`
	Description  string          `json:"description"`
	Chunks       []DocumentChunk `json:"chunks"`
}

type DocumentChunk struct {
	ChunkID  string                 `json:"chunk_id"`
	Text     string                 `json:"text"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
	Vector   []float32              `json:"vector"`
}
