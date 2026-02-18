package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"server/internal/adapter/postgres/sqlc"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	conn   *sqlx.DB
	sqlc   *sqlc.Queries
	goqu   *goqu.Database
	tables goquTableNames
}

func (p *Postgres) MakeGroup() *Group {
	return &Group{
		goqu:   p.goqu,
		tables: p.tables,
	}
}

func (p *Postgres) MakeQuizItem() *QuizItem {
	return &QuizItem{
		goqu:   p.goqu,
		tables: p.tables,
	}
}

func (p *Postgres) MakeDomain() *Domain {
	return &Domain{
		Postgres: p,
	}
}

func (p *Postgres) MakeQuery() *Query {
	return &Query{
		Postgres: p,
	}
}

type Options struct {
	User     string
	Password string
	DB       string
	SSLMode  string
	Port     int
	Host     string
}

func (o *Options) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", o.Host, o.Port, o.User, o.Password, o.DB, o.SSLMode)
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
		conn:   conn,
		sqlc:   sqlc.New(conn),
		goqu:   goqu.New("postgres", conn),
		tables: newGoquTables(),
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

type goquTableNames struct {
	SchoolGroup       exp.IdentifierExpression
	SchoolSubject     exp.IdentifierExpression
	AccountProfile    exp.IdentifierExpression
	AccountStudent    exp.IdentifierExpression
	AccountCredential exp.IdentifierExpression
	AccountTeacher    exp.IdentifierExpression
	QuizInfo          exp.IdentifierExpression
	QuizAttempt       exp.IdentifierExpression
	QuizAssignment    exp.IdentifierExpression
}

func newGoquTables() goquTableNames {
	school := goqu.S("school")
	account := goqu.S("account")
	quiz := goqu.S("quiz")
	return goquTableNames{
		SchoolGroup:       school.Table("group"),
		SchoolSubject:     school.Table("subject"),
		AccountProfile:    account.Table("profile"),
		AccountStudent:    account.Table("student"),
		AccountTeacher:    account.Table("teacher"),
		AccountCredential: account.Table("credential"),
		QuizInfo:          quiz.Table("info"),
		QuizAttempt:       quiz.Table("attempt"),
		QuizAssignment:    quiz.Table("assignment"),
	}
}
