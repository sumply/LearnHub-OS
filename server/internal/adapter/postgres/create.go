package postgres

import (
	"context"
	"server/internal/domain"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (p *Postgres) CreateUser(ctx context.Context, user *domain.User) error {
	tx, err := p.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	if err := p.insertCredentialTx(ctx, tx, user); err != nil {
		return err
	}

	if err := p.insertProfileTx(ctx, tx, user); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}

func (p *Postgres) CreateGroup(ctx context.Context, group *domain.Group) error {
	const query = `
	INSERT INTO groups
	VALUES(:id, :name)
	`

	ext := p.selectExecuter(ctx)

	_, err := ext.NamedExecContext(ctx, query, group)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateSubject(ctx context.Context, subject *domain.Subject) error {
	const query = `
	INSERT INTO subjects(id, name)
	VALUES(:id, :name)
	`

	ext := p.selectExecuter(ctx)

	_, err := ext.NamedExecContext(ctx, query, subject)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateStudent(ctx context.Context, student *domain.Student) error {
	tx, err := p.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	if err := p.insertCredentialTx(ctx, tx, &student.User); err != nil {
		return err
	}

	if err := p.insertProfileTx(ctx, tx, &student.User); err != nil {
		return err
	}

	if err := p.insertStudentTx(ctx, tx, student); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateTeacher(ctx context.Context, teacher *domain.Teacher) error {
	tx, err := p.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	if err := p.insertCredentialTx(ctx, tx, &teacher.User); err != nil {
		return err
	}

	if err := p.insertProfileTx(ctx, tx, &teacher.User); err != nil {
		return err
	}

	if err := p.insertTeacherTx(ctx, tx, teacher); err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) insertCredentialTx(ctx context.Context, tx *sqlx.Tx, user *domain.User) error {
	const query = `
	INSERT INTO account.credentials(
		id,
		email,
		pwd_hash
	) 
	VALUES(
		:id,
		:email,
		:pwd_hash
	)`

	_, err := tx.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) insertProfileTx(ctx context.Context, tx *sqlx.Tx, user *domain.User) error {
	const query = `
	INSERT INTO account.credentials(
		id,
		email,
		pwd_hash
	) 
	VALUES(
		:id,
		:email,
		:pwd_hash
	)`

	_, err := tx.NamedExecContext(ctx, query, user)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) insertStudentTx(ctx context.Context, tx *sqlx.Tx, student *domain.Student) error {
	const query = `
	INSERT INTO account.students(
		profile_id,
		group_id
	)
	VALUES(
		:id,
		:group_id
	)
	`

	_, err := tx.NamedExecContext(ctx, query, student)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) insertTeacherTx(ctx context.Context, tx *sqlx.Tx, teacher *domain.Teacher) error {
	const queryTeachers = `
	INSERT INTO account.teachers(
		profile_id,
		group_id,
		subject_id
	)
	VALUES(
		$1,
		unnest($2::UUID[]),
		unnest($3::UUID[])
	)
	`

	_, err := tx.ExecContext(ctx, queryTeachers, teacher.ID, pq.Array(teacher.Groups), pq.Array(teacher.Subjects))
	if err != nil {
		return err
	}

	return nil
}
