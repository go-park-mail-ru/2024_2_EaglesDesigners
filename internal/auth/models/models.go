package models

type contextKey string

const (
	UserIDKey    contextKey = "userId"
	UserKey      contextKey = "user"
	MuxParamsKey contextKey = "muxParams"
)
