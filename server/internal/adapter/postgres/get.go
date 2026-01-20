package postgres

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (p *Postgres) Users(ctx context.Context) ([]domain.User, error) {
	const query = `
	SELECT 
		id, 
		first_name, 
		last_name, 
		role, 
		created_at
	FROM users
	`

	rows, err := p.conn.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
		resp = append(resp, user)
	}

	return resp, nil
}

func (p *Postgres) StudentMap(ctx context.Context) (map[uuid.UUID]domain.Student, error) {
	const query = `
	SELECT user_id, group_id
	FROM user_groups
	`
	var result []struct {
		UserID  uuid.UUID `db:"user_id"`
		GroupID uuid.UUID `db:"group_id"`
	}

	err := p.conn.SelectContext(ctx, &result, query)
	if err != nil {
		return nil, err
	}

	students := make(map[uuid.UUID]domain.Student)
	for i := range result {
		students[result[i].UserID] = domain.Student{
			User:  result[i].UserID,
			Group: result[i].GroupID,
		}
	}

	return students, nil
}

func (p *Postgres) TeacherMap(ctx context.Context) (map[uuid.UUID]domain.Teacher, error) {
	const query = `
	SELECT 
		user_id,
		array_agg(group_id) AS group_ids
	FROM user_groups
	GROUP BY user_id
	`

	rows, err := p.conn.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teachers := make(map[uuid.UUID]domain.Teacher)
	for rows.Next() {
		var teacher domain.Teacher
		if err := rows.Scan(&teacher.User, pq.Array(&teacher.Groups)); err != nil {
			return nil, err
		}
		teachers[teacher.User] = teacher
	}

	return teachers, nil
}

func (p *Postgres) Groups(ctx context.Context) ([]domain.Group, error) {
	const query = `
	SELECT 
		id,
		name
	FROM groups
	`

	var groups []domain.Group
	err := p.conn.SelectContext(ctx, &groups, query)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (p *Postgres) Subjects(ctx context.Context) ([]domain.Subject, error) {
	const query = `
	SELECT id, name
	FROM subjects
	`

	var subjects []domain.Subject
	err := p.conn.SelectContext(ctx, &subjects, query)
	if err != nil {
		return nil, err
	}

	return subjects, nil
}
