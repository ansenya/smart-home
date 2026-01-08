package db

import (
	"context"
	"device-service/internal/config"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Client struct {
	cfg *config.PostgresConfig
	log *slog.Logger

	DB *gorm.DB
}

func NewClient(cfg *config.Container) *Client {
	return &Client{
		cfg: cfg.PostgresConfig,
		log: cfg.Log,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := gorm.Open(postgres.Open(c.cfg.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent)},
	)
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(c.cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(c.cfg.MaxIdleConns)
	sqlDB.SetConnMaxIdleTime(c.cfg.ConnMaxIdleTime)
	sqlDB.SetConnMaxLifetime(c.cfg.ConnMaxLifetime)

	c.DB = db

	return nil
}

func (c *Client) Close() error {
	if c.DB == nil {
		return nil
	}
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
