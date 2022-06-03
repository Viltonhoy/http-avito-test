package server

import (
	"context"
	"fmt"
	"http-avito-test/internal/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Server struct {
	logger        *zap.Logger
	httpServer    *http.Server
	afterShutdown func()
}

func New() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zap.NewDevelopment: %v", err)
	}
	defer logger.Sync()

	ctx := context.Background()

	var s, _ = storage.NewStore(ctx, logger)
	h := Handler{
		Store: s,
	}

	http.HandleFunc("/read", h.ReadUser)
	http.HandleFunc("/deposit", h.AccountDeposit)
	http.HandleFunc("/transf", h.TransferCommand)
	http.HandleFunc("/history", h.ReadUserHistory)
	http.HandleFunc("/withdrawal", h.AccountWithdrawal)
	port := ":9090"
	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Errorf("failed to listen and serve: %v", err)
		return
	}

}

func (s *Server) Start() error {
	ConnClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		s.logger.Info("shutting down http server")

		err := s.httpServer.Shutdown(context.Background())
		if err != nil {
			s.logger.Error("failed to shutdown http server", zap.Error(err))
			return
		}

		s.logger.Info("http server is stoped")

		close(ConnClosed)
		return
	}()

	s.logger.Info("staring http server")
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("failed to listen and serve: %v", err)
	}

	<-ConnClosed

	s.afterShutdown()

	return nil
}

// func Start() {
// 	ctx, cansel := context.WithCancel(context.Background())
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		listen(ctx)

// 		wg.Done()
// 	}()
// }

// type Service interface {
// 	Init(ctx context.Context) error
// 	Ping(ctx context.Context) error
// 	Close() error
// }

// type (
// 	ServiceKeeper struct {
// 		Services []Service
// 		state    int32
// 	}
// )

// func (s *ServiceKeeper) initAllServices(ctx context.Context) error {
// 	for i := range s.Services {
// 		if err := s.Services[i].Init(ctx); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// const (
// 	srvStateInit int32 = iota
// 	srvStateReady
// 	srvStateRunning
// 	srvStateShutdown
// 	srvStateOff
// )

// func (s *ServiceKeeper) checkState(old, new int32) bool {
// 	return atomic.CompareAndSwapInt32(&s.state, old, new)
// }

// func (s *ServiceKeeper) Init(ctx context.Context) error {
// 	if !s.checkState(srvStateInit, srvStateReady) {
// 		return errBadOrderType
// 	}
// 	return s.initAllServices(ctx)
// }
