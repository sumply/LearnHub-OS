package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"server/internal/adapter/postgres/sqlc"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	conn *sqlx.DB
	sqlc *sqlc.Queries
}

type Options struct {
	User     string
	Password string
	DB       string
	SSLMode  string
}

func (o *Options) String() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", o.User, o.Password, o.DB, o.SSLMode)
}

func New(opt Options) (*Postgres, error) {
	conn, err := sqlx.Connect("postgres", opt.String())
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{
		conn: conn,
		sqlc: sqlc.New(conn),
	}, nil
}

func (p *Postgres) Transaction() Transaction {
	return Transaction{
		conn: p.conn,
	}
}

func (p *Postgres) selectExecuter(ctx context.Context) executer {
	tx, ok := txFromContext(ctx)
	if !ok {
		return p.conn
	}
	return tx
}

type executer interface {
	sqlx.ExtContext
	NamedExecContext(context.Context, string, any) (sql.Result, error)
}
