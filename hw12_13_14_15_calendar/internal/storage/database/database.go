package database

import (
	"context"
	"fmt"

	"github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/AHTI6IOTIK/hw_otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	conf   *config.DatabaseConf
	logger *logger.Logger
	pg     *pgxpool.Pool
}

type IDatabase interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	DB() *pgxpool.Pool
}

func New(
	conf *config.DatabaseConf,
	logger *logger.Logger,
) *DB {
	return &DB{
		conf:   conf,
		logger: logger,
	}
}

func (d *DB) DB() *pgxpool.Pool {
	return d.pg
}

func (d *DB) Connect(ctx context.Context) error {
	connCfg, err := pgxpool.ParseConfig(d.conf.GetDsn())
	if err != nil {
		return fmt.Errorf("parse database dsn: %w", err)
	}

	conn, err := pgxpool.NewWithConfig(ctx, connCfg)
	if err != nil {
		return fmt.Errorf("database open connect: %w", err)
	}

	pingErr := conn.Ping(ctx)
	if pingErr != nil {
		return fmt.Errorf("ping database connect: %s", pingErr.Error())
	}

	d.pg = conn

	return nil
}

func (d *DB) Close(_ context.Context) error {
	d.pg.Close()

	return nil
}
