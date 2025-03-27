package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"grpc-auth/internal/config"
)

// SQL-запрос на вставку задачи
const (
	registerUserQuery = `INSERT INTO users (email, username, password_Hash, first_Name, last_Name, is_Active, role ) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	checkUserQuery    = `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
)

type repository struct {
	pool *pgxpool.Pool
}

type Repository interface {
	RegisterUser(ctx context.Context, user UserRegister) error
	CheckUserExists(ctx context.Context, username string) (bool, error)
}

// NewRepository - создание нового экземпляра репозитория с подключением к PostgreSQL
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
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repository{pool}, nil
}

func (r *repository) RegisterUser(ctx context.Context, user UserRegister) error {
	_, err := r.pool.Exec(
		ctx, registerUserQuery, user.Email, user.Username, user.PasswordHash,
		user.FirstName, user.LastName, user.IsActive, user.Role)

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
