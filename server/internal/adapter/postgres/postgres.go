package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	conn *sqlx.DB
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
	}, nil
}
