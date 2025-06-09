package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log/slog"
	"strings"
	"sync"
	"time"
)

var pgOnce sync.Once

type DB struct {
	Client *pgxpool.Pool
}

type DBConfig struct {
	Host              string `mapstructure:"PG_HOST"`
	Port              int    `mapstructure:"PG_PORT"`
	UserName          string `mapstructure:"PG_USERNAME"`
	Password          string `mapstructure:"PG_PASSWORD"`
	DBName            string `mapstructure:"PG_DBNAME"`
	MaxConns          int32
	MinConns          int32
	MaxConnLifeTime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

func LoadConfig(configFile string) (*DBConfig, error) {
	var cfg DBConfig
	// Set file name for environment configuration
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv() // Read environment variables
	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	// Unmarshal into DBConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func WithPgxConfig(dbConfig *DBConfig) *pgx.ConnConfig {
	// Create the dsn string
	connString := strings.TrimSpace(fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d",
		dbConfig.UserName, dbConfig.Password, dbConfig.DBName,
		dbConfig.Host, dbConfig.Port))
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		slog.Error("Error parsing connection config", slog.String("error", err.Error()))
	}
	return config
}

func NewPg(ctx context.Context, dbConfig *DBConfig, pgxConfig *pgx.ConnConfig) (*pgxpool.Pool, error) {
	var db *pgxpool.Pool
	// Parse the pool configuration from connection string
	config, err := pgxpool.ParseConfig(pgxConfig.ConnString())
	if err != nil {
		slog.Error("Error parsing pool config", slog.String("error", err.Error()))
		return nil, err
	}
	// Apply pool-specific configurations
	config.MaxConns = dbConfig.MaxConns
	config.MinConns = dbConfig.MinConns
	config.MaxConnLifetime = dbConfig.MaxConnLifeTime
	config.MaxConnIdleTime = dbConfig.MaxConnIdleTime
	config.HealthCheckPeriod = dbConfig.HealthCheckPeriod
	// Ensure singleton instance of the pool
	pgOnce.Do(func() {
		db, err = pgxpool.NewWithConfig(ctx, config)
	})
	// Verify the connection
	if err = db.Ping(ctx); err != nil {
		slog.Error("Unable to ping database", slog.String("error", err.Error()))
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	slog.Info("Successfully connected to database")
	return db, nil
}

func (db *DB) MonitorPoolStats() {
	stats := db.Client.Stat()

	slog.Info("Pool stats",
		slog.Int("total_connections", int(stats.TotalConns())),
		slog.Int("acquired_connections", int(stats.AcquiredConns())),
		slog.Int("idle_connections", int(stats.IdleConns())),
		slog.Int("max_connections", int(stats.MaxConns())),
	)
}
