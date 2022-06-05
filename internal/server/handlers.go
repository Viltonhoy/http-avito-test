package server

import "go.uber.org/zap"

type Handler struct {
	Logger *zap.Logger
	Store  Storager
}
