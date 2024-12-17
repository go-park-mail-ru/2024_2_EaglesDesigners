package models

import "github.com/google/uuid"

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

type GetStickerPackResponse struct {
	Photo string   `json:"photo" valid:"-" example:"url1"`
	URLs  []string `json:"stickers" valid:"-" example:"url1,url2,url3"`
}

type StickerPack struct {
	Photo  string    `json:"photo" valid:"-" example:"url1"`
	PackID uuid.UUID `json:"id" valid:"-"`
}

type StickerPacks struct {
	Packs []StickerPack `json:"packs" valid:"-"`
}
