package domain

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Attempt struct {
	ID              uuid.UUID
	QuizID          uuid.UUID
	UserID          uuid.UUID
	SelectedAnswers []IAnswer
	TotalScore      int
	StartedAt       time.Time
	EndedAt         *time.Time
}

func NewAttempt(quiz *Quiz, userID uuid.UUID) (*Attempt, error) {
	if quiz == nil {
		panic("quiz is panic")
	}

	err := NewError("attempt")

	if userID == uuid.Nil {
		err.add("userID", fmt.Errorf("userID is empty"))
	}

	if !err.Empty() {
		return nil, err
	}

	attemptID := uuid.New()

	answers := make([]IAnswer, 0, len(quiz.Questions))
	maker := &answerMaker{attemptID: attemptID}
	for _, question := range quiz.Questions {
		if err := question.Accept(maker); err != nil {
			return nil, err
		}
		answers = append(answers, maker.maked)
	}

	return &Attempt{
		ID:              attemptID,
		QuizID:          quiz.ID,
		UserID:          userID,
		SelectedAnswers: answers,
		StartedAt:       time.Now().UTC(),
	}, nil
}

func (a *Attempt) Finish() {
	endedAt := time.Now().UTC()
	a.EndedAt = &endedAt
}

func (a *Attempt) Check(q *Quiz) error {
	domainErr := NewError("attempt")

	if a.QuizID != q.ID {
		domainErr.add("quiz_id", fmt.Errorf("invalid selected quiz"))
	}

	questionMap := make(map[uuid.UUID]IQuestion)
	for i := range q.Questions {
		questionMap[q.Questions[i].ID()] = q.Questions[i]
	}

	a.TotalScore = 0
	for _, answer := range a.SelectedAnswers {
		if err := answer.Check(questionMap[answer.QuestionID()]); err != nil {
			domainErr.add("selected_answer", err)
		}

		if answer.IsCorrect() {
			a.TotalScore += answer.Score()
		}
	}

	if !domainErr.Empty() {
		return domainErr
	}

	return nil
}

type answerMaker struct {
	attemptID uuid.UUID
	maked     IAnswer
}

func (m *answerMaker) VisitSingleQuestion(q *SingleQuestion) error {
	m.maked = &SingleAnswer{
		Answer: Answer{
			id:         uuid.New(),
			attemptID:  m.attemptID,
			questionID: q.id,
		},
	}
	return nil
}

func (m *answerMaker) VisitMultipleQuestion(q *MultipleQuestion) error {
	m.maked = &MultipleAnswer{
		Answer: Answer{
			id:         uuid.New(),
			attemptID:  m.attemptID,
			questionID: q.id,
		},
		Selected: make([]string, 0),
	}
	return nil
}

func (m *answerMaker) VisitNumericQuestion(q *NumericQuestion) error {
	m.maked = &NumericAnswer{
		Answer: Answer{
			id:         uuid.New(),
			attemptID:  m.attemptID,
			questionID: q.id,
		},
	}
	return nil
}

type IAnswer interface {
	ID() uuid.UUID
	AttemptID() uuid.UUID
	QuestionID() uuid.UUID
	Score() int
	IsCorrect() bool
	Check(IQuestion) error
	Accept(AnswerVisitor) error
}

type AnswerVisitor interface {
	VisitSingleAnswer(*SingleAnswer) error
	VisitMultipleAnswer(*MultipleAnswer) error
	VisitNumericAnswer(*NumericAnswer) error
}

type Answer struct {
	id         uuid.UUID
	attemptID  uuid.UUID
	questionID uuid.UUID
	score      int
	isCorrect  bool
}

func RestoreAnswer(id, attemptID, questionID uuid.UUID, score int, isCorrect bool) Answer {
	return Answer{
		id:         id,
		attemptID:  attemptID,
		questionID: questionID,
		score:      score,
		isCorrect:  isCorrect,
	}
}

func (a *Answer) ID() uuid.UUID {
	return a.id
}

func (a *Answer) AttemptID() uuid.UUID {
	return a.attemptID
}

func (a *Answer) QuestionID() uuid.UUID {
	return a.questionID
}

func (a *Answer) Score() int {
	return a.score
}

func (a *Answer) IsCorrect() bool {
	return a.isCorrect
}

func (a *Answer) questionIDIsCorrect(questionID uuid.UUID) error {
	if a.questionID != questionID {
		return fmt.Errorf("incorrect question id")
	}
	return nil
}

type SingleAnswer struct {
	Answer
	Selected string
}

func (a *SingleAnswer) VisitSingleQuestion(q *SingleQuestion) error {
	if err := a.questionIDIsCorrect(q.id); err != nil {
		return err
	}

	domainErr := NewError("single_answer")

	if !slices.Contains(q.Options, a.Selected) {
		domainErr.add(a.Selected, fmt.Errorf("%s is not contains in options", a.Selected))
	}

	if !domainErr.Empty() {
		return domainErr
	}

	if q.Correct == a.Selected {
		a.isCorrect = true
		a.score = q.score
	} else {
		a.isCorrect = false
		a.score = 0
	}
	return nil
}

func (a *SingleAnswer) VisitMultipleQuestion(q *MultipleQuestion) error {
	return nil
}

func (a *SingleAnswer) VisitNumericQuestion(q *NumericQuestion) error {
	return nil
}

func (a *SingleAnswer) Check(q IQuestion) error {
	return q.Accept(a)
}

func (a *SingleAnswer) Accept(visitor AnswerVisitor) error {
	return visitor.VisitSingleAnswer(a)
}

type MultipleAnswer struct {
	Answer
	Selected []string
}

func (a *MultipleAnswer) VisitSingleQuestion(q *SingleQuestion) error {
	return nil
}

func (a *MultipleAnswer) VisitMultipleQuestion(q *MultipleQuestion) error {
	if err := a.questionIDIsCorrect(q.id); err != nil {
		return err
	}

	domainErr := NewError("multiple_answer")

	if !a.selectedAnswerIsUnique() {
		domainErr.add("answers", fmt.Errorf("answers are not unique"))
	}

	if err := a.allSelectedContainsInOptions(q.Options); err != nil {
		domainErr.add("contains", err)
	}

	if !domainErr.Empty() {
		return domainErr
	}

	if a.answerIsCorrect(q.Correct) {
		a.score = q.score
		a.isCorrect = true
	} else {
		a.isCorrect = false
		a.score = 0
	}

	return nil
}

func (a *MultipleAnswer) selectedAnswerIsUnique() bool {
	selectedMap := make(map[string]struct{})
	for i := range a.Selected {
		selectedMap[a.Selected[i]] = struct{}{}
	}

	return len(selectedMap) == len(a.Selected)
}

func (a *MultipleAnswer) allSelectedContainsInOptions(options []string) error {
	domainErr := NewError("selected_answers")
	for _, selected := range a.Selected {
		if !slices.Contains(options, selected) {
			domainErr.add(selected, fmt.Errorf("%s is not contains in options", selected))
		}
	}
	if !domainErr.Empty() {
		return domainErr
	}
	return nil
}

func (a *MultipleAnswer) answerIsCorrect(correct []string) bool {
	var numberCorrect int

	for _, selected := range a.Selected {
		if slices.Contains(correct, selected) {
			numberCorrect++
		}
	}

	return len(a.Selected) == len(correct) && len(correct) == numberCorrect
}

func (a *MultipleAnswer) VisitNumericQuestion(q *NumericQuestion) error {
	return nil
}

func (a *MultipleAnswer) Check(q IQuestion) error {
	return q.Accept(a)
}

func (a *MultipleAnswer) Accept(visitor AnswerVisitor) error {
	return visitor.VisitMultipleAnswer(a)
}

type NumericAnswer struct {
	Answer
	Selected float32
}

func (a *NumericAnswer) VisitSingleQuestion(q *SingleQuestion) error {
	return nil
}

func (a *NumericAnswer) VisitMultipleQuestion(q *MultipleQuestion) error {
	return nil
}

func (a *NumericAnswer) VisitNumericQuestion(q *NumericQuestion) error {
	if err := a.questionIDIsCorrect(q.id); err != nil {
		return err
	}

	if a.Selected == q.Correct {
		a.score = q.score
		a.isCorrect = true
	} else {
		a.isCorrect = false
		a.score = 0
	}
	return nil
}

func (a *NumericAnswer) Check(q IQuestion) error {
	return q.Accept(a)
}

func (a *NumericAnswer) Accept(visitor AnswerVisitor) error {
	return visitor.VisitNumericAnswer(a)
}
