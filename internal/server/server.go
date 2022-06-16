package server

import (
	"context"
	"errors"
	"fmt"
	"http-avito-test/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type Server struct {
	logger        *zap.Logger
	httpServer    *http.Server
	afterShutdown func()
}

type ServerConfig struct {
	Host string `env:"ADDR_HOST"`
	Port int    `env:"ADDR_PORT"`
}

func New(logger *zap.Logger, storage *storage.Storage, afterShutdown func(), e Exchanger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("no logger provided")
	}

	cfg := ServerConfig{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	defer logger.Sync()

	mux := http.NewServeMux()

	h := Handler{
		Logger:    logger,
		Store:     storage,
		Exchanger: e,
	}

	mux.HandleFunc("/read", h.ReadUser)
	mux.HandleFunc("/deposit", h.AccountDeposit)
	mux.HandleFunc("/transf", h.TransferCommand)
	mux.HandleFunc("/history", h.ReadUserHistory)
	mux.HandleFunc("/withdrawal", h.AccountWithdrawal)

	httpServer := http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	server := &Server{
		logger:        logger,
		httpServer:    &httpServer,
		afterShutdown: afterShutdown,
	}
	return server, nil
}

func (s *Server) Start() error {
	idleConnClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		s.logger.Info("shutting down http server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error("failed to shutdown http server", zap.Error(err))
			return
		}

		s.logger.Info("http server is stopped")

		close(idleConnClosed)
	}()

	s.logger.Info("starting http server")
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("failed to listen and serve: %v", err)
	}

	<-idleConnClosed

	s.afterShutdown()

	return nil
}
