package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/flyingbisons/vacation-tracker-float-sync/internal/float"
	"github.com/flyingbisons/vacation-tracker-float-sync/internal/integrator"
	"github.com/flyingbisons/vacation-tracker-float-sync/internal/integrator/store"
	"github.com/flyingbisons/vacation-tracker-float-sync/internal/vacation"
	"github.com/syumai/workers/cloudflare"
	"github.com/syumai/workers/cloudflare/cron"
	"github.com/syumai/workers/cloudflare/d1"
)

func handleErr(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(msg))
}

func main() {
	cron.ScheduleTask(task)

	//http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	//	ctx := req.Context()
	//	logger := newLogger()
	//	err := run(logger, ctx)
	//	if err != nil {
	//		handleErr(w, http.StatusInternalServerError, err.Error())
	//		return
	//	}
	//
	//	w.Write([]byte("OK"))
	//})
	//
	//workers.Serve(nil) // use http.DefaultServeMux
}

func run(logger *slog.Logger, ctx context.Context) error {
	c, err := d1.OpenConnector(ctx, "VTFLOAT")
	if err != nil {
		return fmt.Errorf("failed to initialize DB: %v", err)
	}

	// get current time
	currentTime, err := getDatetime(ctx)
	if err != nil {
		return fmt.Errorf("failed to get time: %v", err)
	}

	// prepare integrator
	vtKey := cloudflare.Getenv(ctx, "VACATION_TRACKER_API_KEY")
	vClient := vacation.NewClient(vtKey, vacation.ApiUrl, currentTime.Now)
	floatKey := cloudflare.Getenv(ctx, "FLOAT_API_KEY")
	fClient := float.NewClient(floatKey, float.ApiUrl)
	db := sql.OpenDB(c)
	repository := store.NewRequestRepositoryDB(db, currentTime.Now)
	manager, err := integrator.New(repository, fClient, vClient, logger)

	//run integrator
	err = manager.Sync(ctx)
	if err != nil {
		return fmt.Errorf("failed to sync: %v", err)
	}

	return nil
}

func task(ctx context.Context, event *cron.Event) error {
	logger := newLogger()
	logger.Info("running task")
	return run(logger, ctx)
}

func newLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)
	return logger
}
