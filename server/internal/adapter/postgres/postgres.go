package postgres

import (
	"context"
	"fmt"
	"server/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
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

func (p *Postgres) CreateUser(ctx context.Context, user *domain.User) error {
	const query = `
	INSERT INTO users 
	VALUES(
		:id,
		:first_name,
		:last_name,
		:email,
		:pwd_hash,
		:role,
		:created_at
	)`
	_, err := p.conn.NamedExecContext(ctx, query, user)

	return err
}

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

func (p *Postgres) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const query = `
	DELETE FROM users
	WHERE id=$1
	`

	_, err := p.conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) UpdateProfile(ctx context.Context, id uuid.UUID, firstName, lastName string) error {
	const query = `
	UPDATE users
	SET first_name = :first_name, last_name = :last_name
	WHERE id = :id
	`
	arg := struct {
		ID        uuid.UUID `db:"id"`
		FirstName string    `db:"first_name"`
		LastName  string    `db:"last_name"`
	}{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	}

	_, err := p.conn.NamedExecContext(ctx, query, arg)
	if err != nil {
		return err
	}

	return nil
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

func (p *Postgres) CreateGroup(ctx context.Context, group *domain.Group) error {
	const query = `
	INSERT INTO groups
	VALUES(:id, :name)
	`

	_, err := p.conn.NamedExecContext(ctx, query, group)
	if err != nil {
		return err
	}

	return nil
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

func (p *Postgres) UpdateGroup(ctx context.Context, id uuid.UUID, name string) error {
	const query = `
	UPDATE groups
	SET name = :name
	WHERE id = :id
	`
	arg := struct {
		ID   uuid.UUID `db:"id"`
		Name string    `db:"name"`
	}{
		ID:   id,
		Name: name,
	}

	_, err := p.conn.NamedExecContext(ctx, query, arg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	const query = `
	DELETE FROM groups
	WHERE id = $1
	`

	_, err := p.conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
