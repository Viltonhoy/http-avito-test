package server

import (
	"http-avito-test/internal/storage"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestReadUserHostory(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)
		m.EXPECT().ReadUserHistoryList(int64(1), "Data").Return(storage.Transf{
			ID: 1,
		}, nil)
	})
}
