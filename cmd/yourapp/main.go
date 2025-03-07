package main

import (
	"context"
	"embed"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/saylorsolutions/x/httpx"
	"github.com/saylorsolutions/x/signalx"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
	"yourapp/cmd/yourapp/internal/routes"
	"yourapp/foundation/urlprefix"
)

var (
	//go:embed static
	staticAssets embed.FS
	version      = "noversion"
	gitHash      = "nohash"
)

func main() {
	logger := log.Default()
	logger.Printf("Starting yourapp v%s: %s\n", version, gitHash)
	ctx := signalx.SignalCtx(context.Background(), os.Interrupt, syscall.SIGTERM)
	if err := run(ctx, logger); err != nil {
		logger.Fatalln("Failed to run yourapp:", err)
	}
	logger.Println("yourapp server shut down gracefully")
}

func run(ctx context.Context, logger *log.Logger) error {
	db, err := initData()
	if err != nil {
		logger.Println("[ERR] Failed to connect to database:", err)
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	authSvc, err := initAuth(log.Default(), db)
	if err != nil {
		logger.Println("[ERR] Failed to initialize auth store:", err)
		return err
	}
	ro := &routes.Router{
		Log:          logger,
		AuthSvc:      authSvc,
		Pool:         db,
		StaticAssets: staticAssets,
	}
	handler := httpx.Wrap(ro.ServeMux(),
		httpx.LoggingMiddleware(httpx.StdLogger(logger)),
		httpx.RecoveryMiddleware(panicHandlerFunc(func(cause any) {
			logger.Println("[ERR] Panic encountered:", cause)
		})),
	)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: urlprefix.Group(handler),
	}

	log.Println("Starting yourapp server...")
	if err := httpx.ListenAndServeCtx(ctx, srv, 5*time.Second); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return err
		}
		log.Println("[ERR] Error running down server:", err)
	}
	return nil
}

type panicHandlerFunc func(cause any)

func (f panicHandlerFunc) Handle(cause any) {
	f(cause)
}
