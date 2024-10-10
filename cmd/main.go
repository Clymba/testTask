package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Clymba/testTask/internal/repository"
	"github.com/Clymba/testTask/internal/service"
	"github.com/Clymba/testTask/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"os/signal"
	"time"

	"github.com/Clymba/testTask/internal/config"
	"github.com/Clymba/testTask/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errc := make(chan error, 1)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger.InitLogger()

	db, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return err
	}

	repo := repository.New(db, repository.Config{Timeout: time.Second})
	srv := service.New(repo)

	app := fiber.New(fiber.Config{
		AppName:      cfg.Name,
		ErrorHandler: ErrorHandler,
	})

	srvServer := server.New(ctx, app, srv)

	go func() {
		if err := srvServer.Listen(cfg.Port); err != nil {
			errc <- err
		}
	}()

	select {
	case err = <-errc:
		return err
	case <-ctx.Done():
		return srvServer.Shutdown()
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	var code int
	var message string

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else {
		code = fiber.StatusInternalServerError
		message = "Internal Server Error"
	}

	fmt.Printf("Error: %s\n", err)

	return c.Status(code).JSON(fiber.Map{"error": message})
}
