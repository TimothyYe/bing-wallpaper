package run

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/TimothyYe/bing-wallpaper/handler"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const defaultPort = 9000

var argPort int

var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Run this API service",
	Long:  `Run this API service`,
	Run:   run,
}

func init() {
	// 0 sentinel means "fall back to PORT env or defaultPort".
	Cmd.Flags().IntVarP(&argPort, "port", "p", 0, fmt.Sprintf("port to listen on (default %d; overrides PORT env)", defaultPort))
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

// resolvePort returns the listen port using the precedence flag > PORT env > default.
func resolvePort() int {
	if argPort > 0 {
		return argPort
	}
	if env := os.Getenv("PORT"); env != "" {
		if p, err := strconv.Atoi(env); err == nil && p > 0 {
			return p
		}
		log.Printf("ignoring invalid PORT=%q", env)
	}
	return defaultPort
}

func run(cmd *cobra.Command, args []string) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery(), CORSMiddleware())
	router.GET("/", handler.RootHandler)
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	addr := fmt.Sprintf(":%d", resolvePort())
	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run the server in the background; ListenAndServe blocks until the
	// listener errors or Shutdown is called.
	go func() {
		log.Printf("API is running on %s...", addr)
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
