package server

import (
	"http-avito-test/internal/storage"
)

type Handler struct {
	Store *storage.Storage
}
