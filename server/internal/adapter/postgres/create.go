package postgres

import (
	"context"
	"encoding/json"
	"fmt"
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
	INSERT INTO school.group(id, name)
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
	INSERT INTO school.subject(id, name)
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

func (p *Postgres) CreateQuiz(ctx context.Context, quiz *domain.Quiz) error {
	tx, err := p.conn.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	err = p.insertQuizTx(ctx, tx, quiz)
	if err != nil {
		return err
	}

	err = p.insertQuizContentTx(ctx, tx, quiz)
	if err != nil {
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
	INSERT INTO account.credential(
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
		return fmt.Errorf("insert credential: %w", err)
	}

	return nil
}

func (p *Postgres) insertProfileTx(ctx context.Context, tx *sqlx.Tx, user *domain.User) error {
	const query = `
	INSERT INTO account.profile(
		id,
		first_name,
		last_name,
		role,
		created_at
	) 
	VALUES(
		:id,
		:first_name,
		:last_name,
		:role,
		:created_at
	)`

	_, err := tx.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("insert profile: %w", err)
	}

	return nil
}

func (p *Postgres) insertStudentTx(ctx context.Context, tx *sqlx.Tx, student *domain.Student) error {
	const query = `
	INSERT INTO account.student(
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
	INSERT INTO account.teacher(
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

func (p *Postgres) insertQuizTx(ctx context.Context, tx *sqlx.Tx, quiz *domain.Quiz) error {
	const query = `
	INSERT INTO quiz.info(
		id,
		title,
		summary,
		subject_id,
		owner_id,
		total_score,
		created_at
	)
	VALUES (
		:id,
		:title,
		:summary,
		:subject_id,
		:owner_id,
		:total_score,
		:created_at
	)
	`

	_, err := tx.NamedExecContext(ctx, query, quiz)
	if err != nil {
		return fmt.Errorf("insert quiz: %w", err)
	}

	return nil
}

func (p *Postgres) insertQuizContentTx(ctx context.Context, tx *sqlx.Tx, quiz *domain.Quiz) error {
	const query = `
	INSERT INTO quiz.content(
		quiz_id,
		content
	)
	VALUES(
		$1,
		$2
	)
	`

	content, err := json.Marshal(quiz.Content)
	if err != nil {
		return fmt.Errorf("insert quiz content: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, quiz.ID, content)
	if err != nil {
		return fmt.Errorf("insert quiz content: %w", err)
	}

	return nil
}
