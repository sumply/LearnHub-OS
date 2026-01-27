package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/domain"
	"server/internal/query"
	"time"

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

func (p *Postgres) Quiz(ctx context.Context, id uuid.UUID) (query.Quiz, error) {
	const sqlQuery = `
SELECT 
	q.id AS quiz_id, 
	q.title AS quiz_title,
	q.summary AS quiz_summary,
	q.total_score AS quiz_total_score,
	q.created_at AS quiz_created_at,
	c.content AS quiz_content,
	u.id AS owner_id,
	u.first_name AS owner_first_name,
	u.last_name AS owner_last_name,
	u.role AS owner_role,
	u.created_at AS owner_created_at,
	s.id AS subject_id,
	s.name AS subject_name
FROM quiz.info AS q
JOIN quiz.content AS c
	ON c.quiz_id = q.id
JOIN account.profile AS u 
	ON u.id = q.owner_id
JOIN school.subject AS s
	ON s.id = q.subject_id
WHERE q.id = $1
	`

	result := p.conn.QueryRowxContext(ctx, sqlQuery, id)

	var row QuizRow
	err := result.StructScan(&row)
	if err != nil {
		return query.Quiz{}, err
	}

	return row.ToQuery(), nil
}

func (p *Postgres) QuizWithoutAnswers(ctx context.Context, id uuid.UUID) (query.Quiz, error) {
	const sqlQuery = `
SELECT 
	q.id AS quiz_id, 
	q.title AS quiz_title,
	q.summary AS quiz_summary,
	q.total_score AS quiz_total_score,
	q.created_at AS quiz_created_at,
	c.content AS quiz_content,
	u.id AS owner_id,
	u.first_name AS owner_first_name,
	u.last_name AS owner_last_name,
	u.role AS owner_role,
	u.created_at AS owner_created_at,
	s.id AS subject_id,
	s.name AS subject_name
FROM quiz.info AS q
JOIN quiz.content AS c
	ON c.quiz_id = q.id
JOIN account.profile AS u 
	ON u.id = q.owner_id
JOIN school.subject AS s
	ON s.id = q.subject_id
WHERE q.id = $1
	`

	row := p.conn.QueryRowxContext(ctx, sqlQuery, id)

	var result QuizRow
	err := row.StructScan(&result)
	if err != nil {
		return query.Quiz{}, fmt.Errorf("struct scan: %w", err)
	}

	quiz := result.ToQuery()
	quiz.DeleteAnswers()
	return quiz, nil
}

func (p *Postgres) QuizItems(ctx context.Context) ([]query.QuizItem, error) {
	const sqlQuery = `
SELECT 
	q.id AS quiz_id, 
	q.title AS quiz_title,
	q.summary AS quiz_summary,
	q.total_score AS quiz_total_score,
	q.created_at AS quiz_created_at,
	u.id AS owner_id,
	u.first_name AS owner_first_name,
	u.last_name AS owner_last_name,
	u.role AS owner_role,
	u.created_at AS owner_created_at,
	s.id AS subject_id,
	s.name AS subject_name
FROM quiz.info AS q
JOIN account.profile AS u 
	ON u.id = q.owner_id
JOIN school.subject AS s
	ON s.id = q.subject_id
	`

	rows, err := p.conn.QueryxContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	var items []query.QuizItem

	for rows.Next() {
		var row QuizItemRow
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		item := row.ToQuery()

		items = append(items, item)
	}

	return items, nil
}

func (p *Postgres) DomainQuiz(ctx context.Context, id uuid.UUID) (domain.Quiz, error) {
	const sqlQuery = `
SELECT 
	i.id,
	i.title,
	i.summary,
	i.owner_id,
	i.subject_id,
	i.total_score,
	i.created_at,
	c.content
FROM quiz.info AS i
JOIN quiz.content AS c
	ON c.quiz_id = i.id
WHERE i.id = $1
	`

	row := p.conn.QueryRowxContext(ctx, sqlQuery, id)
	var result DomainQuizRow
	err := row.StructScan(&result)
	if err != nil {
		return domain.Quiz{}, err
	}

	return result.ToDomain(), nil
}

type DomainQuizRow struct {
	ID         uuid.UUID `db:"id"`
	Title      string    `db:"title"`
	Summary    string    `db:"summary"`
	OwnerID    uuid.UUID `db:"owner_id"`
	SubjectID  uuid.UUID `db:"subject_id"`
	TotalScore int       `db:"total_score"`
	CreatedAt  time.Time `db:"created_at"`
	Content    []byte    `db:"content"`
}

func (d *DomainQuizRow) ToDomain() domain.Quiz {
	var content []domain.QuestionAggregate
	json.Unmarshal(d.Content, &content)
	return domain.Quiz{
		ID:         d.ID,
		Title:      d.Title,
		Summary:    d.Summary,
		OwnerID:    d.OwnerID,
		SubjectID:  d.SubjectID,
		TotalScore: d.TotalScore,
		Content:    content,
		CreatedAt:  d.CreatedAt,
	}
}

type QuizRow struct {
	QuizItemRow
	Content []byte `db:"quiz_content"`
}

func (q *QuizRow) ToQuery() query.Quiz {
	var content []query.QuizQuestion
	json.Unmarshal(q.Content, &content)
	/* <- Отличие только в этой строчке. У меня вопросы хранятся в jsonb
	for i := range content {
		delete(content[i].Payload, "correct")
	}
	*/
	query := query.Quiz{
		QuizItem: q.QuizItemRow.ToQuery(),
		Content:  content,
	}
	return query
}

type QuizItemRow struct {
	QuizID         uuid.UUID `db:"quiz_id"`
	QuizTitle      string    `db:"quiz_title"`
	QuizSummary    string    `db:"quiz_summary"`
	QuizTotalScore int       `db:"quiz_total_score"`
	QuizCreatedAt  time.Time `db:"quiz_created_at"`
	OwnerID        uuid.UUID `db:"owner_id"`
	OwnerFirstName string    `db:"owner_first_name"`
	OwnerLastName  string    `db:"owner_last_name"`
	OwnerRole      string    `db:"owner_role"`
	OwnerCreatedAt time.Time `db:"owner_created_at"`
	SubjectID      uuid.UUID `db:"subject_id"`
	SubjectName    string    `db:"subject_name"`
}

func (q *QuizItemRow) ToQuery() query.QuizItem {
	return query.QuizItem{
		ID:      q.QuizID,
		Title:   q.QuizTitle,
		Summary: q.QuizSummary,
		Owner: query.User{
			ID:        q.OwnerID,
			FirstName: q.OwnerFirstName,
			LastName:  q.OwnerLastName,
			Role:      q.OwnerRole,
		},
		Subject: query.Subject{
			ID:   q.SubjectID,
			Name: q.SubjectName,
		},
		TotalScore: q.QuizTotalScore,
		CreatedAt:  q.QuizCreatedAt,
	}
}
