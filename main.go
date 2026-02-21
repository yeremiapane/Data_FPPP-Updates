package main

import (
	"os"
	"os/signal"
	"syscall"

	"DataFPP_Dingtalk/config"
	"DataFPP_Dingtalk/internal/scheduler"
	"DataFPP_Dingtalk/internal/service"
	"DataFPP_Dingtalk/internal/sheets"
	"DataFPP_Dingtalk/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	client, err := sheets.NewClient(cfg.CredentialsPath)
	if err != nil {
		logger.Error("Failed to create sheets client: %v", err)
		os.Exit(1)
	}

	reader := sheets.NewReader(client, cfg.SourceSheetID)
	writer := sheets.NewWriter(client, cfg.DestSheetID)
	syncService := service.NewSyncService(reader, writer)

	// Run sync once at startup
	logger.Info("Running initial sync...")
	if err := syncService.SyncData(); err != nil {
		logger.Error("Initial sync failed: %v", err)
	}

	// Setup scheduler for daily 08:00 WIB
	s, err := scheduler.Setup(syncService, cfg.Timezone)
	if err != nil {
		logger.Error("Failed to setup scheduler: %v", err)
		os.Exit(1)
	}
	s.Start()
	logger.Info("Scheduler started. Next sync at 08:00 %s daily.", cfg.Timezone)

	// Wait for shutdown signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	logger.Info("Shutting down...")
	if err := s.Shutdown(); err != nil {
		logger.Error("Scheduler shutdown error: %v", err)
	}
}
