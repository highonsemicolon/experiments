package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/highonsemicolon/experiments/appointment-booking/config"
	"github.com/highonsemicolon/experiments/appointment-booking/database"
	"github.com/highonsemicolon/experiments/appointment-booking/internal/handler"
	"github.com/highonsemicolon/experiments/appointment-booking/internal/repository"
	"github.com/highonsemicolon/experiments/appointment-booking/internal/service"
	"github.com/highonsemicolon/experiments/appointment-booking/logging"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logging.NewZerologAdapter(logging.LoggingOption{
		Format: "json",
		Level:  "info",
	})

	cfg := &config.Config{}
	if err := config.Load(cfg, config.ConfigLoaderOption{
		Prefix: "booking.",
		Logger: log,
	}); err != nil {
		log.Fatal("failed to load config", err)
	}

	db := database.Connect(cfg)

	coachRepo := repository.NewCoachRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	coachSvc := service.NewCoachService(coachRepo)
	bookingSvc := service.NewBookingService(bookingRepo, coachRepo)

	h := handler.NewHandler(coachSvc, bookingSvc)

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	handler.RegisterHandlers(r, handler.NewStrictHandler(h, nil))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Info("server starting on :" + cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed to start", err)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server forced to shutdown", err)
	}

	log.Info("server exited cleanly")
}
