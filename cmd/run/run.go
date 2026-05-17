package run

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/TimothyYe/bing-wallpaper/handler"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Run this API service",
	Long:  `Run this API service`,
	Run:   run,
}

// CORSMiddleware permits cross-origin reads from any browser. The endpoint
// is unauthenticated, so Allow-Credentials is intentionally omitted — it
// is incompatible with a wildcard origin and would be rejected by browsers.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Cache-Control")
		c.Header("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func run(cmd *cobra.Command, args []string) {
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/", handler.RootHandler)

	srv := &http.Server{
		Addr:              ":9000",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run the server in the background; ListenAndServe blocks until the
	// listener errors or Shutdown is called.
	go func() {
		log.Println("API is running on :9000...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for SIGINT/SIGTERM, then drain in-flight requests.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("shutdown signal received, draining...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
	log.Println("server stopped")
}
