package repo

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"grpc-auth/internal/config"
)

const (
	registerUserQuery = `INSERT INTO users (email, username, password_hash, first_Name, last_Name) VALUES ($1, $2, $3, $4, $5)`
	checkUserQuery    = `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	getUser           = `SELECT email, username, password_hash, first_name, last_name, created_at, updated_at FROM users WHERE username = $1`
	updateLoginTime   = `UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE username = $1`
)

type repository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	RegisterUser(ctx context.Context, user User) error
	CheckUserExists(ctx context.Context, username string) (bool, error)
	GetUser(ctx context.Context, username string) (*User, error)
	ShuttingDownPostgres() error
	UpdateLoginTime(ctx context.Context, username string) error
}

func NewRepository(ctx context.Context, cfg config.PostgreSQL) (Repository, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s 
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	configConnect, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	configConnect.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, configConnect)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	if err := applyMigrations(pool); err != nil {
		return nil, errors.Wrap(err, "failed to apply migrations")
	}

	return &repository{pool}, nil
}

func applyMigrations(pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	driver, err := pgxMigrate.WithInstance(sqlDB, &pgxMigrate.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to initialize pgx migrate driver")
	}

	migrationsPath := "app/migrations"
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to apply migrations")
	}

	return nil
}

func (r *repository) RegisterUser(ctx context.Context, user User) error {
	_, err := r.pool.Exec(
		ctx, registerUserQuery, user.Email, user.Username, user.PassHash,
		user.FirstName, user.LastName)

	if err != nil {
		return errors.Wrap(err, "failed to register user")
	}
	return nil
}

func (r *repository) CheckUserExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, checkUserQuery, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repository) GetUser(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.pool.QueryRow(ctx, getUser, username).Scan(
		&user.Email,
		&user.Username,
		&user.PassHash,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get user by username")
	}
	return &user, nil
}

func (r *repository) ShuttingDownPostgres() error {
	if r.pool != nil {
		r.pool.Close()
		return nil
	}
	return errors.New("postgres pool is empty")
}

func (r *repository) UpdateLoginTime(ctx context.Context, username string) error {
	_, err := r.pool.Exec(ctx, updateLoginTime, username)
	if err != nil {
		return fmt.Errorf("failed to update login time: %w", err)
	}

	return nil
}
