package postgres

import (
	"context"
	"fmt"
	"server/internal/logger"
	domain "server/internal/new_domain"
	repository "server/internal/new_repository"

	"github.com/jackc/pgx"
)

type (
	User struct {
		tx *pgx.Tx
	}
	Repository struct {
		tx *pgx.Tx
	}
	UnitOfWork struct {
		conn *pgx.Conn
	}
)

func New(conf pgx.ConnConfig) (*UnitOfWork, error) {
	conn, err := pgx.Connect(conf)
	if err != nil {
		return nil, err
	}
	return &UnitOfWork{
		conn: conn,
	}, nil
}

func (uof *UnitOfWork) Handle(f func(repository.Repository) error) error {
	tx, err := uof.conn.Begin()
	if err != nil {
		return err
	}
	if err := f(&Repository{tx: tx}); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Repository) User() repository.User {
	return &User{
		tx: r.tx,
	}
}

func (u *User) Save(ctx context.Context, user *domain.User) (domain.UserID, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(user),
	)
	log.Debug("Called a save postgres repostiory method")
	if user == nil {
		return 0, nil
	}
	if err := u.checkTransaction(); err != nil {
		return 0, err
	}
	const (
		query1 = `
			INSERT INTO 
				users.credentials(login, email, password_hash) 
			VALUES ($1, $2, $3)
			RETURNING id`
		query2 = `
			INSERT INTO
				users.profiles(
					id,
					first_name,
					last_name,
					middle_name,
					role,
					access
				)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
	)
	cred := user.Profile().Credential
	if cred == nil {
		return 0, fmt.Errorf("credential is nil")
	}
	row := u.tx.QueryRow(query1, cred.Login, cred.Email, cred.PwdHash)
	var id domain.UserID
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	profile := user.Profile()
	_, err := u.tx.Exec(query2, id, profile.FirstName, profile.LastName, profile.MiddleName, "none", "user")
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *User) SaveStudent(context.Context, *domain.Student) (domain.StudentID, error) {
	return 0, nil
}
func (u *User) SaveTeacher(context.Context, *domain.Teacher) (domain.TeacherID, error) {
	return 0, nil
}

func (u *User) Get(context.Context, *repository.UserFilter) ([]domain.UserData, error) {
	return nil, nil
}

func (u *User) Update(context.Context, domain.UserData) error {
	return nil
}

func (u *User) Delete(context.Context, domain.UserID) error {
	return nil
}

func (u *User) checkTransaction() error {
	if u.tx == nil || u.tx.Status() != pgx.TxStatusInProgress {
		return fmt.Errorf("transtaction is not in progress")
	}
	return nil
}
