package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/adapters/dtos"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/adapters/handler"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/adapters/repository"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/core/service"
	"github.com/mosmo1212312121/hexagonal_practice_go/internal/infrastructure"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func main() {
	r := gin.Default()
	db, err := infrastructure.ConnectDB()
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&repository.UserModel{})
	if err != nil {
		panic(err)
	}
	// Health Check
	// reg := prometheus.NewRegistry()
	// reg.MustRegister(
	// 	collectors.NewGoCollector(),
	// 	collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	// )
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.Use(r)
	r.GET("/health", func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		// check db connection
		resp := dtos.HealthCheckResponse{}
		sqlDB, err := db.DB()
		if err != nil {
			resp.Database = dtos.DOWN
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		if err := sqlDB.Ping(); err != nil {
			resp.Database = dtos.DOWN
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Status = dtos.UP
		resp.Database = dtos.UP
		c.JSON(http.StatusOK, resp)
	})

	userRepo := repository.NewUserRepository(db)
	_ = userRepo
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(r, userService)

	ug := r.Group("users")
	ug.POST("", userHandler.RegisterUser)
	ug.GET("/:id", userHandler.GetUserByID)

	//implement Graceful Shutdown
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exited")
}
