package domain

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	ID          uuid.UUID
	Title       string
	Summary     string
	OwnerID     uuid.UUID
	SubjectID   uuid.UUID
	Questions   []Question
	Deadline    *time.Time
	MaxAttempts int
	TotalScore  int
	CreatedAt   time.Time
}

func NewQuiz(
	ownerID, subjectID uuid.UUID,
	title, summary string,
	questions []Question,
	maxAttempts int,
	deadline *time.Time,
) (Quiz, error) {
	quizID := uuid.New()

	err := NewError("quiz")
	if ownerID == uuid.Nil {
		err.Add("ownerID", fmt.Errorf("ownerID is empty"))
	}

	if subjectID == uuid.Nil {
		err.Add("subjectID", fmt.Errorf("subjectID is empty"))
	}

	if title == "" {
		err.Add("title", fmt.Errorf("title is empty"))
	}

	if summary == "" {
		err.Add("summary", fmt.Errorf("summary is empty"))
	}

	if len(questions) == 0 {
		err.Add("questions", fmt.Errorf("questions is empty"))
	}

	totalCount := 0
	for i := range questions {
		totalCount += questions[i].Score
		questions[i].QuizID = quizID
	}

	if maxAttempts <= 0 {
		err.Add("max_attempts", fmt.Errorf("max_attempts less 0"))
	}

	if deadline != nil && deadline.Before(time.Now().UTC()) {
		err.Add("deadline", fmt.Errorf("deadline is before now"))
	}

	if !err.Empty() {
		return Quiz{}, err
	}

	return Quiz{
		ID:          quizID,
		Title:       title,
		Summary:     summary,
		OwnerID:     ownerID,
		SubjectID:   subjectID,
		Questions:   questions,
		Deadline:    deadline,
		MaxAttempts: maxAttempts,
		TotalScore:  totalCount,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func (q *Quiz) CheckAttempt(attempt *Attempt) error {
	if attempt == nil {
		return fmt.Errorf("attempt is nil")
	}

	attempt.Score = 0

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

	return nil
}

type Question struct {
	ID      uuid.UUID
	QuizID  uuid.UUID
	Text    string
	Details QuestionDetails
	Score   int
}

func NewQuestion(text string, details QuestionDetails, score int) (Question, error) {
	err := NewError("question")
	if text == "" {
		err.Add("text", fmt.Errorf("text is empty"))
	}

	if score <= 0 {
		err.Add("score", fmt.Errorf("score less 0"))
	}

	if details == nil {
		err.Add("details", fmt.Errorf("details is nil"))
	}

	if !err.Empty() {
		return Question{}, err
	}

	if err := details.Validate(); err != nil {
		return Question{}, err
	}

	return Question{
		ID:      uuid.New(),
		Text:    text,
		Details: details,
		Score:   score,
	}, nil
}

type QuestionDetails interface {
	Validate() error
	Variant() QuestionType
	ReviewAnswer(answer any) (bool, error)
}

type SingleChoiceQuestion struct {
	Options []string
	Correct string
}

func NewSingleChoiceQuestion(options []string, correct string) (SingleChoiceQuestion, error) {
	single := SingleChoiceQuestion{
		Options: slices.Clone(options),
		Correct: correct,
	}
	if err := single.Validate(); err != nil {
		return SingleChoiceQuestion{}, err
	}
	return single, nil
}

func (s *SingleChoiceQuestion) Validate() error {
	if len(s.Options) == 0 {
		return fmt.Errorf("options is nil")
	}
	if !slices.Contains(s.Options, s.Correct) {
		return fmt.Errorf("options do not contain the correct")
	}
	return nil
}

func (s *SingleChoiceQuestion) ReviewAnswer(answer any) (bool, error) {
	_answer, ok := answer.(string)
	if !ok {
		return false, fmt.Errorf("answer is incorrect type (%T)", answer)
	}
	if s.Correct != _answer {
		return false, nil
	}
	return true, nil
}

func (s *SingleChoiceQuestion) Variant() QuestionType {
	return TypeSingleChoice
}

type MultipleChoiceQuestion struct {
	Options []string `json:"options"`
	Correct []string `json:"correct"`
}

func NewMultipleChoiceQuestion(options, correct []string) (MultipleChoiceQuestion, error) {
	multiple := MultipleChoiceQuestion{
		Options: slices.Clone(options),
		Correct: slices.Clone(correct),
	}
	if err := multiple.Validate(); err != nil {
		return MultipleChoiceQuestion{}, err
	}
	return multiple, nil
}

func (m *MultipleChoiceQuestion) Validate() error {
	if len(m.Options) == 0 {
		return fmt.Errorf("options is empty")
	}
	if len(m.Correct) == 0 {
		return fmt.Errorf("correct is empty")
	}
	if len(m.Correct) > len(m.Options) {
		return fmt.Errorf("more correct answers than options")
	}
	for i := range m.Correct {
		if !slices.Contains(m.Options, m.Correct[i]) {
			return fmt.Errorf("options do not contain the correct")
		}
	}
	return nil
}

func (m *MultipleChoiceQuestion) ReviewAnswer(answer any) (bool, error) {
	var _answer []string

	if a, ok := answer.([]string); ok {
		answer = a
	} else if a, ok := answer.([]any); ok {
		_answer = make([]string, len(a))
		for i := range _answer {
			s, ok := a[i].(string)
			if !ok {
				return false, fmt.Errorf("answer is incorrect type (%T)", answer)
			}
			_answer[i] = s
		}
	} else {
		return false, fmt.Errorf("answer is incorrect type (%T)", answer)
	}

	for i := range m.Correct {
		if m.Correct[i] != _answer[i] {
			return false, nil
		}
	}
	return true, nil
}

func (m *MultipleChoiceQuestion) Variant() QuestionType {
	return TypeMultipleChoice
}

type NumericQuestion struct {
	Correct float64 `json:"correct"`
}

func NewNumericQuestion(correct float64) (NumericQuestion, error) {
	numeric := NumericQuestion{
		Correct: correct,
	}
	if err := numeric.Validate(); err != nil {
		return NumericQuestion{}, err
	}
	return numeric, nil
}

func (n *NumericQuestion) Validate() error {
	return nil
}

func (n *NumericQuestion) ReviewAnswer(answer any) (bool, error) {
	_answer, ok := n.castAnyToFloat64(answer)
	if !ok {
		return false, fmt.Errorf("answer is incorrect type (%T)", answer)
	}
	if n.Correct != _answer {
		return false, nil
	}
	return true, nil
}

func (n *NumericQuestion) castAnyToFloat64(a any) (float64, bool) {
	switch number := a.(type) {
	case
		int:
		return float64(number), true
	case int64:
		return float64(number), true
	case float32:
		return float64(number), true
	case float64:
		return float64(number), true
	default:
		return 0, false
	}
}

func (n *NumericQuestion) Variant() QuestionType {
	return TypeNumeric
}

type QuestionType string

const (
	TypeSingleChoice   QuestionType = "single"
	TypeMultipleChoice QuestionType = "multiple"
	TypeNumeric        QuestionType = "numeric"
)
