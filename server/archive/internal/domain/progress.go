package domain

import (
	"fmt"
	"server/archive/internal/common"
	"time"
)

type ProgressStatus common.Enum

const (
	ProgressStatusNotStarted ProgressStatus = iota
	ProgressStatusProgress
	ProgressStatusPendingReview
	ProgressStatusCompleted
)

type QuizProgress struct {
	ID            common.ID
	Quiz          *Quiz
	User          *User
	Status        ProgressStatus
	Score         int
	Answers       []*Answer
	CompletedDate *time.Time
	StartDate     *time.Time
}

func NewQuizProgress(userID common.ID, quiz *Quiz) *QuizProgress {
	return &QuizProgress{
		Quiz:    quiz,
		User:    &User{ID: userID},
		Status:  ProgressStatusNotStarted,
		Score:   0,
		Answers: initAnswers(quiz.Questions),
	}
}

func initAnswers(questions []*Question) []*Answer {
	answers := make([]*Answer, len(questions))
	for i, question := range questions {
		answers[i] = NewAnswer(question)
	}
	return answers
}

func (p *QuizProgress) InProgress() bool {
	return p.Status == ProgressStatusProgress
}

func (p *QuizProgress) IsCompleted() bool {
	return p.Status == ProgressStatusCompleted
}

func (p *QuizProgress) IsNotStarted() bool {
	return p.Status == ProgressStatusNotStarted
}

func (p *QuizProgress) IsPendingReview() bool {
	return p.Status == ProgressStatusPendingReview
}

func (p *QuizProgress) GetAnswer(answerID common.ID) (*Answer, bool) {
	for _, answer := range p.Answers {
		if answer.ID == answerID {
			return answer, true
		}
	}
	return nil, false
}

func (p *QuizProgress) Complete() error {
	if !p.InProgress() {
		return fmt.Errorf("the progress (id=%d) is not in progress", p.ID)
	}

	now := time.Now().UTC()
	p.CompletedDate = &now

	if p.Quiz.HasWritten {
		p.Status = ProgressStatusPendingReview
	} else {
		p.Status = ProgressStatusCompleted
		p.countScore()
	}
	return nil
}

func (p *QuizProgress) countScore() {
	p.Score = 0
	for i := range p.Answers {
		if p.Answers[i].Status == AnswerStatusCorrect {
			p.Score++
		}
	}
}

func (p *QuizProgress) CompleteReview() error {
	if !p.IsPendingReview() {
		return fmt.Errorf("quiz progress (id=%d) is not pending review", p.ID)
	}
	p.Status = ProgressStatusCompleted
	p.countScore()
	return nil
}

func (p *QuizProgress) Start() error {
	if !p.IsNotStarted() {
		return fmt.Errorf("the progress (id=%d) is started or completed", p.ID)
	}
	now := time.Now().UTC()
	p.StartDate = &now
	p.Status = ProgressStatusProgress
	return nil
}

func (p QuizProgress) Copy() *QuizProgress {
	if p.Quiz != nil {
		p.Quiz = p.Quiz.Copy()
	}
	if p.User != nil {
		p.User = p.User.Copy()
	}
	for i, answer := range p.Answers {
		if answer != nil {
			p.Answers[i] = answer.Copy()
		}
	}
	return &p
}

type AnswerStatus common.Enum

const (
	AnswerStatusNotAnswered AnswerStatus = iota
	AnswerStatusPendingReview
	AnswerStatusIncorrect
	AnswerStatusCorrect
)

type Answer struct {
	ID       common.ID
	Text     string
	Status   AnswerStatus
	Question *Question
}

func NewAnswer(question *Question) *Answer {
	return &Answer{
		Status:   AnswerStatusNotAnswered,
		Question: question,
	}
}

func (a *Answer) Submit(text string) error {
	if a.Question.IsWritten {
		a.Status = AnswerStatusPendingReview
	} else {
		option, ok := a.findOption(text)
		if !ok {
			return fmt.Errorf("option (text=%s) is not exists", text)
		}
		if option.IsCorrect {
			a.Status = AnswerStatusCorrect
		} else {
			a.Status = AnswerStatusIncorrect
		}
	}
	a.Text = text
	return nil
}

func (a *Answer) findOption(text string) (*Option, bool) {
	for _, opt := range a.Question.Options {
		if opt.Text == text {
			return opt, true
		}
	}
	return nil, false
}

func (a *Answer) Review(isCorrect bool) error {
	if !a.Question.IsWritten {
		return fmt.Errorf("question (id=%d) has options", a.Question.ID)
	}
	if isCorrect {
		a.Status = AnswerStatusCorrect
	} else {
		a.Status = AnswerStatusIncorrect
	}
	return nil
}

func (a Answer) Copy() *Answer {
	if a.Question != nil {
		a.Question = a.Question.Copy()
	}
	return &a
}
