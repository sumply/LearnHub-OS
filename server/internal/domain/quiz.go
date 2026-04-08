package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	ID          uuid.UUID
	Title       string
	Summary     string
	OwnerID     uuid.UUID
	SubjectID   uuid.UUID
	Questions   []IQuestion
	GroupIDs    uuid.UUIDs
	Deadline    *time.Time
	MaxAttempts int
	TotalScore  int
	CreatedAt   time.Time
}

func NewQuiz(
	ownerID, subjectID uuid.UUID,
	title, summary string,
	questions []IQuestion,
	groupIDs uuid.UUIDs,
	maxAttempts int,
	deadline *time.Time,
) (*Quiz, error) {
	quizID := uuid.New()

	totalCount := 0
	for i := range questions {
		totalCount += questions[i].Score()
		questions[i].setQuizID(quizID)
	}

	quiz := &Quiz{
		ID:          quizID,
		Title:       strings.TrimSpace(title),
		Summary:     strings.TrimSpace(summary),
		OwnerID:     ownerID,
		SubjectID:   subjectID,
		Questions:   questions,
		Deadline:    deadline,
		GroupIDs:    groupIDs,
		MaxAttempts: maxAttempts,
		TotalScore:  totalCount,
		CreatedAt:   time.Now().UTC(),
	}

	if err := quiz.Validate(); err != nil {
		return nil, err
	}

	return quiz, nil
}

func (q *Quiz) Validate() error {
	domainErr := NewError("quiz")

	if q.OwnerID == uuid.Nil {
		domainErr.add("ownerID", fmt.Errorf("ownerID is empty"))
	}

	if q.SubjectID == uuid.Nil {
		domainErr.add("subjectID", fmt.Errorf("subjectID is empty"))
	}

	q.Title = strings.TrimSpace(q.Title)

	if q.Title == "" {
		domainErr.add("title", fmt.Errorf("title is empty"))
	}

	q.Summary = strings.TrimSpace(q.Summary)

	if q.Summary == "" {
		domainErr.add("summary", fmt.Errorf("summary is empty"))
	}

	if len(q.Questions) == 0 {
		domainErr.add("questions", fmt.Errorf("questions is empty"))
	}

	if len(q.GroupIDs) == 0 {
		domainErr.add("group_ids", fmt.Errorf("group_ids is empty"))
	}

	if q.MaxAttempts <= 0 {
		domainErr.add("max_attempts", fmt.Errorf("max_attempts less 0"))
	}

	if q.Deadline != nil && q.Deadline.Before(time.Now().UTC()) {
		domainErr.add("deadline", fmt.Errorf("deadline is expired"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

type IQuestion interface {
	ID() uuid.UUID
	Validate() error
	Score() int
	setQuizID(uuid.UUID)
	Accept(QuestionVisitor) error
}

type QuestionVisitor interface {
	VisitSingleQuestion(*SingleQuestion) error
	VisitMultipleQuestion(*MultipleQuestion) error
	VisitNumericQuestion(*NumericQuestion) error
}

type CommonQuestion struct {
	id     uuid.UUID
	QuizID uuid.UUID
	Title  string
	score  int
}

func (q *CommonQuestion) ID() uuid.UUID {
	return q.id
}

func (s *CommonQuestion) Score() int {
	return s.score
}

func (s *CommonQuestion) setQuizID(quizID uuid.UUID) {
	s.QuizID = quizID
}

type SingleQuestion struct {
	CommonQuestion
	Correct string
	Options []string
}

func NewSingleQuestion(title, correct string, options []string, score int) (*SingleQuestion, error) {
	q := &SingleQuestion{
		CommonQuestion: CommonQuestion{
			id:    uuid.New(),
			Title: strings.TrimSpace(title),
			score: score,
		},
		Correct: strings.TrimSpace(correct),
		Options: func() []string {
			trimmed := make([]string, 0, len(options))
			for _, opt := range options {
				trimmed = append(trimmed, strings.TrimSpace(opt))
			}
			return trimmed
		}(),
	}

	if err := q.Validate(); err != nil {
		return nil, err
	}

	return q, nil
}

func RestoreSingleQuestion(id uuid.UUID, title, correct string, options []string, score int) (*SingleQuestion, error) {
	q := &SingleQuestion{
		CommonQuestion: CommonQuestion{
			id:    id,
			Title: title,
			score: score,
		},
		Correct: correct,
		Options: options,
	}

	if err := q.Validate(); err != nil {
		return nil, err
	}

	return q, nil
}

func (s *SingleQuestion) Validate() error {
	domainErr := NewError("single_question")

	if s.Title == "" {
		domainErr.add("title", fmt.Errorf("title is empty"))
	}

	if len(s.Options) == 0 {
		domainErr.add("options", fmt.Errorf("options is nil"))
	}

	uniqueOptions := make(map[string]interface{})
	for i := range s.Options {
		if s.Options[i] == "" {
			domainErr.add("option", fmt.Errorf("option is empty"))
		}
		uniqueOptions[s.Options[i]] = struct{}{}
	}

	if len(uniqueOptions) != len(s.Options) {
		domainErr.add("options", fmt.Errorf("options are not unique"))
	}

	if _, ok := uniqueOptions[s.Correct]; !ok {
		domainErr.add("correct", fmt.Errorf("correct is not contains in options"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

func (s *SingleQuestion) Accept(visitor QuestionVisitor) error {
	return visitor.VisitSingleQuestion(s)
}

type MultipleQuestion struct {
	CommonQuestion
	Correct []string
	Options []string
}

func NewMultipleQuestion(title string, correct, options []string, score int) (*MultipleQuestion, error) {
	q := &MultipleQuestion{
		CommonQuestion: CommonQuestion{
			id:    uuid.New(),
			Title: strings.TrimSpace(title),
			score: score,
		},
		Correct: func() []string {
			trimmed := make([]string, 0, len(correct))
			for _, c := range correct {
				trimmed = append(trimmed, strings.TrimSpace(c))
			}
			return trimmed
		}(),
		Options: func() []string {
			trimmed := make([]string, 0, len(options))
			for _, opt := range options {
				trimmed = append(trimmed, strings.TrimSpace(opt))
			}
			return trimmed
		}(),
	}

	if err := q.Validate(); err != nil {
		return nil, err
	}

	return q, nil
}

func RestoreMultipleQuestion(id uuid.UUID, title string, correct, options []string, score int) (*MultipleQuestion, error) {
	q := &MultipleQuestion{
		CommonQuestion: CommonQuestion{
			id:    id,
			Title: title,
			score: score,
		},
		Correct: correct,
		Options: options,
	}

	if err := q.Validate(); err != nil {
		return nil, err
	}

	return q, nil
}

func (m *MultipleQuestion) Validate() error {
	domainErr := NewError("multiple_question")

	if m.Title == "" {
		domainErr.add("title", fmt.Errorf("title is empty"))
	}

	if len(m.Correct) == 0 {
		domainErr.add("correct", fmt.Errorf("correct length is 0"))
	}

	uniqueCorrect := make(map[string]any)
	for i := range m.Correct {
		uniqueCorrect[m.Correct[i]] = struct{}{}
		if m.Correct[i] == "" {
			domainErr.add("correct", fmt.Errorf("correct is empty"))
		}
	}

	if len(uniqueCorrect) != len(m.Correct) {
		domainErr.add("correct", fmt.Errorf("correct is not unique"))
	}

	if len(m.Options) == 0 {
		return fmt.Errorf("options length is 0")
	}

	uniqueOptions := make(map[string]any)
	for i := range m.Options {
		uniqueOptions[m.Options[i]] = struct{}{}
		if m.Options[i] == "" {
			domainErr.add("option", fmt.Errorf("option is empty"))
		}
	}

	if len(uniqueOptions) != len(m.Options) {
		domainErr.add("options", fmt.Errorf("options are not unique"))
	}

	for i := range m.Correct {
		if _, ok := uniqueOptions[m.Correct[i]]; !ok {
			domainErr.add("correct", fmt.Errorf("correct is not contains in options"))
		}
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

func (m *MultipleQuestion) Accept(visitor QuestionVisitor) error {
	return visitor.VisitMultipleQuestion(m)
}

type NumericQuestion struct {
	CommonQuestion
	Correct float32
}

func NewNumericQuestion(title string, correct float32, score int) (*NumericQuestion, error) {
	q := &NumericQuestion{
		CommonQuestion: CommonQuestion{
			id:    uuid.New(),
			Title: title,
			score: score,
		},
		Correct: correct,
	}

	if err := q.Validate(); err != nil {
		return nil, err
	}

	return q, nil
}

func (n *NumericQuestion) Validate() error {
	domainErr := NewError("numeric_question")

	if n.Title == "" {
		domainErr.add("title", fmt.Errorf("title is empty"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

func (n *NumericQuestion) Accept(visitor QuestionVisitor) error {
	return visitor.VisitNumericQuestion(n)
}

func RestoreNumericQuestion(id uuid.UUID, title string, correct float32, score int) (*NumericQuestion, error) {
	q := &NumericQuestion{
		CommonQuestion: CommonQuestion{
			id:    id,
			Title: title,
			score: score,
		},
		Correct: correct,
	}

	if err := q.Validate(); err != nil {
		return nil, err
	}

	return q, nil
}
