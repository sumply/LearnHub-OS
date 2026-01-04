package domain

import (
	"fmt"
	"server/internal/common"
	"slices"
	"strings"
	"time"
)

type Quiz struct {
	ID              common.ID
	Title           string
	Summary         string
	NumberQuestions uint8
	HasWritten      bool
	IsForEveryone   bool
	CreatedAt       time.Time

	Owner     *User
	Subject   *Subject
	Groups    []*Group
	Questions []*Question
}

func NewQuiz(title, summary string, questions []*Question, ownerID common.ID, subjectID common.ID, groupIDs []common.ID) (*Quiz, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, fmt.Errorf("a quiz title is empty")
	}
	summary = strings.TrimSpace(summary)
	if summary == "" {
		return nil, fmt.Errorf("a quiz summary is empty")
	}
	if len(questions) == 0 {
		return nil, fmt.Errorf("a quiz does not have questions")
	}

	hasWritten := slices.ContainsFunc(questions, func(q *Question) bool {
		return q.IsWritten
	})

	new := &Quiz{
		Title:           title,
		Summary:         summary,
		Questions:       questions,
		HasWritten:      hasWritten,
		Owner:           &User{ID: ownerID},
		NumberQuestions: uint8(len(questions)),
		Subject:         &Subject{ID: subjectID},
	}

	if len(groupIDs) == 0 {
		new.IsForEveryone = true
	} else {
		groups := make([]*Group, len(groupIDs))
		for i, id := range groupIDs {
			groups[i] = &Group{ID: id}
		}
		new.Groups = groups
	}
	return new, nil
}

func (q *Quiz) IsOwner(userID common.ID) bool {
	return q.Owner.ID == userID
}

func (q Quiz) Copy() *Quiz {
	if q.Owner != nil {
		q.Owner = q.Owner.Copy()
	}
	for i, question := range q.Questions {
		if question != nil {
			q.Questions[i] = question.Copy()
		}
	}
	for i, group := range q.Groups {
		if group != nil {
			q.Groups[i] = group.Copy()
		}
	}

	if q.Subject != nil {
		q.Subject = q.Subject.Copy()
	}
	return &q
}

type Question struct {
	ID        common.ID
	Text      string
	IsWritten bool
	Options   []*Option
}

func NewQuestion(text string, options []*Option) (*Question, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty")
	}

	new := &Question{
		Text: text,
	}

	if len(options) != 0 {
		hasCorrect := false
		for _, opt := range options {
			if opt.IsCorrect {
				hasCorrect = true
			}
		}
		if !hasCorrect {
			return nil, fmt.Errorf("a question does not have a correct option")
		}
		new.Options = options
	} else {
		new.IsWritten = true
	}
	return new, nil
}

func (q Question) Copy() *Question {
	for i, opt := range q.Options {
		q.Options[i] = opt.Copy()
	}
	return &q
}

type Option struct {
	ID        common.ID
	Text      string
	IsCorrect bool
}

func NewOption(text string, correct bool) (*Option, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("a option text in empty")
	}
	return &Option{
		Text:      text,
		IsCorrect: correct,
	}, nil
}

func (o Option) Copy() *Option {
	return &o
}
