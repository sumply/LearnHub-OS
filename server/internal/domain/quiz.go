package domain

import (
	"fmt"
	"slices"
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
) (Quiz, error) {
	quizID := uuid.New()

	totalCount := 0
	for i := range questions {
		totalCount += questions[i].Score()
		questions[i].setQuizID(quizID)
	}

	quiz := Quiz{
		ID:          quizID,
		Title:       title,
		Summary:     summary,
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
		return Quiz{}, err
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

func (q *Quiz) CheckAttempt(attempt *Attempt) error {
	if attempt == nil {
		return fmt.Errorf("attempt is nil")
	}

	attempt.Score = 0

	/*
		questionMap := make(map[uuid.UUID]*Question)
		for i := range q.Questions {
			questionMap[q.Questions[i].ID] = &q.Questions[i]
		}

		for i := range attempt.Answers {
			question, ok := questionMap[attempt.Answers[i].QuestionID]
			if !ok {
				return fmt.Errorf("unknow questionID (id=%d)", attempt.Answers[i].QuestionID)
			}
			ok, err := question.Details.ReviewAnswer(attempt.Answers[i].Answer)
			if err != nil {
				return err
			}
			attempt.Answers[i].IsCorrect = ok
			if ok {
				attempt.Answers[i].Score = question.Score
			}
			attempt.Score += attempt.Answers[i].Score
		}
	*/

	return nil
}

type IQuestion interface {
	Validate() error
	ReviewAnswer(any) (bool, error)
	Score() int
	setQuizID(uuid.UUID)
	Accept(QuestionVisitor) error
}

type QuestionVisitor interface {
	VisitSingle(*SingleQuestion) error
	VisitMultiple(*MultipleQuestion) error
	VisitNumeric(*NumericQuestion) error
}

type CommonQuestion[T any] struct {
	ID      uuid.UUID
	QuizID  uuid.UUID
	Title   string
	score   int
	Correct T
}

type SingleQuestion struct {
	CommonQuestion[string]
	Options []string
}

func NewSingleQuestion(title, correct string, options []string, score int) (*SingleQuestion, error) {
	return &SingleQuestion{
		CommonQuestion: CommonQuestion[string]{
			ID:      uuid.New(),
			Title:   title,
			Correct: correct,
			score:   score,
		},
		Options: options,
	}, nil
}

func (s *SingleQuestion) setQuizID(quizID uuid.UUID) {
	s.QuizID = quizID
}

func (s *SingleQuestion) Score() int {
	return s.score
}

func (s *SingleQuestion) Validate() error {
	domainErr := NewError("single_question")
	if len(s.Options) == 0 {
		domainErr.add("options", fmt.Errorf("options is nil"))
	}

	uniqueOptions := make(map[string]interface{})
	for i := range s.Options {
		uniqueOptions[s.Options[i]] = struct{}{}
	}

	if len(uniqueOptions) != len(s.Options) {
		domainErr.add("options", fmt.Errorf("options are not unique"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

func (s *SingleQuestion) ReviewAnswer(answer any) (bool, error) {
	_answer, ok := answer.(string)
	if !ok {
		return false, fmt.Errorf("answer is incorrect type (%T)", answer)
	}
	if s.Correct != _answer {
		return false, nil
	}
	return true, nil
}

func (s *SingleQuestion) Accept(visitor QuestionVisitor) error {
	return visitor.VisitSingle(s)
}

type MultipleQuestion struct {
	CommonQuestion[[]string]
	Options []string
}

func NewMultipleQuestion(title string, correct, options []string, score int) (*MultipleQuestion, error) {
	return &MultipleQuestion{
		CommonQuestion: CommonQuestion[[]string]{
			ID:      uuid.New(),
			Title:   title,
			score:   score,
			Correct: correct,
		},
		Options: options,
	}, nil
}

func (m *MultipleQuestion) setQuizID(quizID uuid.UUID) {
	m.QuizID = quizID
}

func (m *MultipleQuestion) Score() int {
	return m.score
}

func (m *MultipleQuestion) Validate() error {
	domainErr := NewError("multiple_question")

	if len(m.Options) == 0 {
		return fmt.Errorf("options is nil")
	}

	uniqueOptions := make(map[string]interface{})
	for i := range m.Options {
		uniqueOptions[m.Options[i]] = struct{}{}
	}

	if len(uniqueOptions) != len(m.Options) {
		domainErr.add("options", fmt.Errorf("options are not unique"))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

func (m *MultipleQuestion) ReviewAnswer(answer any) (bool, error) {
	var _answer []string
	if a, ok := answer.([]string); ok {
		_answer = a
	} else if a, ok := answer.([]any); ok {
		_answer = make([]string, len(a))
		for i := 0; i < len(a); i++ {
			if s, ok := a[i].(string); ok {
				_answer[i] = s
			}
		}
	}
	if len(m.Correct) != len(_answer) {
		return false, nil
	}
	for i := range m.Correct {
		if !slices.Contains(m.Correct, _answer[i]) {
			return false, nil
		}
	}
	return true, nil
}

func (m *MultipleQuestion) Accept(visitor QuestionVisitor) error {
	return visitor.VisitMultiple(m)
}

type NumericQuestion struct {
	CommonQuestion[float32]
}

func NewNumericQuestion(title string, correct float32, score int) (*NumericQuestion, error) {
	return &NumericQuestion{
		CommonQuestion: CommonQuestion[float32]{
			ID:      uuid.New(),
			Title:   title,
			Correct: correct,
			score:   score,
		},
	}, nil
}

func (n *NumericQuestion) setQuizID(quizID uuid.UUID) {
	n.QuizID = quizID
}

func (n *NumericQuestion) Score() int {
	return n.score
}

func (n *NumericQuestion) Validate() error {
	return nil
}

func (n *NumericQuestion) ReviewAnswer(answer any) (bool, error) {
	_answer, ok := answer.(float32)
	if !ok {
		return false, fmt.Errorf("answer is incorrect type (%T)", answer)
	}
	if n.Correct != _answer {
		return false, nil
	}
	return true, nil
}

func (n *NumericQuestion) Accept(visitor QuestionVisitor) error {
	return visitor.VisitNumeric(n)
}
