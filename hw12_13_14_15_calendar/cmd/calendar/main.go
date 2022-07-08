package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/server/http"
	sqlstorage "github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "calendar",
		Short: "Calendar service",
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "http",
		Short: "Run http server",
		RunE:  runHTTPServer,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "App version",
		RunE:  runVersion,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE:  runMigrations,
	})

	rootCmd.PersistentFlags().String("config", "/etc/calendar/config.toml", "Path to configuration file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute cmd: %v", err)
	}
}

func runVersion(cmd *cobra.Command, args []string) error {
	printVersion()
	return nil
}

func runHTTPServer(cmd *cobra.Command, args []string) error {
	configFile, err := cmd.Root().PersistentFlags().GetString("config")
	if err != nil {
		return err
	}

	config, err := ReadConfig(configFile)
	if err != nil {
		return err
	}

	logg := logger.New(logger.LevelFromString(config.Logger.Level))
	storage := sqlstorage.New(
		config.Postgres.DSN,
		config.Postgres.MaxOpenConns,
		config.Postgres.MaxIdleConns,
		config.Postgres.ConnMaxLifetime,
		config.Postgres.ConnMaxIdleTime,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := storage.Connect(ctx); err != nil {
		return err
	}

	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, config.HTTP.Host, config.HTTP.Port)

	notifyCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	go func() {
		<-notifyCtx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	return server.Start(notifyCtx)
}

func runMigrations(cmd *cobra.Command, args []string) error {
	configFile, err := cmd.Root().PersistentFlags().GetString("config")
	if err != nil {
		return err
	}

	config, err := ReadConfig(configFile)
	if err != nil {
		return err
	}

	storage := sqlstorage.New(
		config.Postgres.DSN,
		config.Postgres.MaxOpenConns,
		config.Postgres.MaxIdleConns,
		config.Postgres.ConnMaxLifetime,
		config.Postgres.ConnMaxIdleTime,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := storage.Connect(ctx); err != nil {
		return err
	}

	return storage.Migrate(context.Background(), args[0])
}
