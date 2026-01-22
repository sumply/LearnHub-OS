package postgres

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (p *Postgres) DetailedUsers(ctx context.Context) ([]domain.DetailedUser, error) {
	const query = `
SELECT 
	p.id, 
	p.first_name, 
	p.last_name, 
	p.role, 
	p.created_at,
	s.group_id AS s_group_id,
	array_remove(array_agg(t.group_id), NULL) AS t_group_ids,
	array_remove(array_agg(t.subject_id), NULL) AS t_subject_ids
FROM account.profile AS p
LEFT JOIN account.student AS s
	ON s.profile_id = p.id
LEFT JOIN account.teacher AS t
	ON t.profile_id = p.id
GROUP BY 
	p.id, 
	p.first_name,
	p.last_name,
	p.role,
	p.created_at,
	s.group_id
	`

	rows, err := p.conn.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp []domain.DetailedUser
	for rows.Next() {
		var user domain.User
		var sGroupID uuid.UUID
		var tGroupIDs uuid.UUIDs
		var tSubjectIDs uuid.UUIDs
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.CreatedAt,
			&sGroupID,
			pq.Array(&tGroupIDs),
			pq.Array(&tSubjectIDs),
		)
		if err != nil {
			return nil, err
		}

		switch user.Role {
		case domain.RoleStudent:
			student := domain.Student{
				User:  user,
				Group: sGroupID,
			}
			resp = append(resp, &student)
		case domain.RoleTeacher:
			teacher := domain.Teacher{
				User:     user,
				Subjects: tSubjectIDs,
				Groups:   tGroupIDs,
			}
			resp = append(resp, &teacher)
		default:
			resp = append(resp, &user)
		}
	}

	return resp, nil
}

func (p *Postgres) Groups(ctx context.Context) ([]domain.Group, error) {
	const query = `
	SELECT 
		id,
		name
	FROM school.group
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
	FROM school.subject
	`

	var subjects []domain.Subject
	err := p.conn.SelectContext(ctx, &subjects, query)
	if err != nil {
		return nil, err
	}

	return subjects, nil
}
