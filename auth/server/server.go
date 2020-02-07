package server

import (
	"context"
	"github.com/gin-gonic/gin"
	authhttp "github.com/lazy-bees/borsch/auth/delivery/http"
	"github.com/lazy-bees/borsch/auth/repository/memorystorage"
	"github.com/lazy-bees/borsch/auth/usecase"
	"github.com/lazy-bees/borsch/auth/usecase/authusecase"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	httpServer *http.Server
	uc         usecase.UseCase
}

func New() *Server {
	return &Server{
		uc: authusecase.NewAuthUseCase(
			memorystorage.NewUserMemoryStorage(),
			viper.GetString("hash_salt"),
			[]byte(viper.GetString("signing_key")),
			viper.GetDuration("token_ttl"),
		),
	}
}

func (s *Server) Run(port string) error {
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	authhttp.RegisterHTTPEndpoints(router, s.uc)

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
