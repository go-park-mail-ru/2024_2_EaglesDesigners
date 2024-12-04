package models

type UploadFileResponse struct {
	FileID string `json:"file_id"`
}

type FileMetaData struct {
	Filename    string
	ContentType string
	FileSize    int64
}
