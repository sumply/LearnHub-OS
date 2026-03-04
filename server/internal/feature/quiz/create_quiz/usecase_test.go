package create_quiz

import (
	"context"
	"server/internal/domain"
	"server/internal/dto"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type QuizMock struct{}

func (q *QuizMock) Save(context.Context, domain.Quiz) error {
	return nil
}

type TeacherMock struct {
	groups []uuid.UUID
}

func (t *TeacherMock) Get(_ context.Context, id uuid.UUID) (domain.Teacher, error) {
	return domain.Teacher{
		User: domain.User{
			ID: id,
		},
		Groups: t.groups,
	}, nil
}

type IdentityMock struct {
	id   uuid.UUID
	role domain.UserRole
}

func (i *IdentityMock) ID() uuid.UUID {
	return i.id
}

func (i *IdentityMock) Role() domain.UserRole {
	return i.role
}

func makeIdentity(r domain.UserRole) *IdentityMock {
	return &IdentityMock{
		id:   uuid.New(),
		role: r,
	}
}

func makeValidInput(groupIDs []uuid.UUID) Input {
	return Input{
		Title:       "Title",
		Summary:     "Summary",
		SubjectID:   uuid.New(),
		GroupIDs:    groupIDs,
		Deadline:    nil,
		MaxAttempts: 1,
		Questions: []dto.Question{
			{
				Text:  "Text",
				Score: 1,
				Details: dto.QuestionDetails{
					Domain: &domain.SingleChoiceQuestion{
						Options: []string{"true", "false"},
						Correct: "true",
					},
				},
			},
		},
	}
}

func TestUseCase_CreateQuiz(t *testing.T) {
	var (
		ID1 = uuid.New()
		ID2 = uuid.New()
		ID3 = uuid.New()
	)
	tests := []struct {
		name          string
		identity      *IdentityMock
		teacherRepo   *TeacherMock
		input         Input
		expectedError bool
	}{
		{
			name:     "user role is student",
			identity: makeIdentity(domain.RoleStudent),
			teacherRepo: &TeacherMock{
				groups: nil,
			},
			input:         makeValidInput(nil),
			expectedError: true,
		},
		{
			name:     "user role is admin",
			identity: makeIdentity(domain.RoleAdmin),
			teacherRepo: &TeacherMock{
				groups: nil,
			},
			input:         makeValidInput([]uuid.UUID{ID1, ID2, ID3}),
			expectedError: false,
		},
		{
			name:     "teacher has not access to group",
			identity: makeIdentity(domain.RoleTeacher),
			teacherRepo: &TeacherMock{
				groups: []uuid.UUID{ID1, ID2},
			},
			input:         makeValidInput([]uuid.UUID{ID3}),
			expectedError: true,
		},
		{
			name:     "teacher has access to group",
			identity: makeIdentity(domain.RoleTeacher),
			teacherRepo: &TeacherMock{
				groups: []uuid.UUID{ID1, ID2},
			},
			input:         makeValidInput([]uuid.UUID{ID1, ID2}),
			expectedError: false,
		},
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := New(tt.teacherRepo, new(QuizMock))
			_, err := uc.CreateQuiz(context.Background(), tt.identity, tt.input)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
