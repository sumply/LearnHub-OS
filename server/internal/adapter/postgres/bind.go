package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (p *Postgres) BindTeacherToGroups(ctx context.Context, teacherID uuid.UUID, groupIDs uuid.UUIDs) error {
	const query = `
	INSERT INTO user_groups(
		user_id,
		group_id
	)
	VALUES(
		$1,
		unnest($2::UUID[])
	)
	`

	ext := p.selectExecuter(ctx)

	_, err := ext.ExecContext(ctx, query, teacherID, pq.Array(groupIDs))
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) BindTeacherToSubjects(ctx context.Context, teacherID uuid.UUID, subjectIDs uuid.UUIDs) error {
	const query = `
	INSERT INTO user_subjects(
		user_id,
		subject_id
	)
	VALUES(
		$1,
		unnest($2::UUID[])
	)`

	ext := p.selectExecuter(ctx)

	_, err := ext.ExecContext(ctx, query, teacherID, pq.Array(subjectIDs))
	if err != nil {
		return err
	}

	return nil
}
