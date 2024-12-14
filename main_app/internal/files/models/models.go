package models

type UploadFileResponse struct {
	FileID string `json:"file_id"`
}

type FileMetaData struct {
	Filename    string
	ContentType string
	FileSize    int64
}

type Payload struct {
	URL      string
	Filename string
	Size     int64
}
