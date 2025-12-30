package domain

import (
	"fmt"
	"server/internal/common"
	"strings"
	"time"
)

type Quiz struct {
	ID              common.ID
	Title           string
	Summary         string
	NumberQuestions uint8
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

	new := &Quiz{
		Title:     title,
		Summary:   summary,
		Questions: questions,
		Owner:     &User{ID: UserID(ownerID)},
		Subject:   &Subject{ID: SubjectID(subjectID)},
	}

	if len(groupIDs) == 0 {
		new.IsForEveryone = true
	} else {
		groups := make([]*Group, len(groupIDs))
		for i, id := range groupIDs {
			groups[i] = &Group{ID: GroupID(id)}
		}
		new.Groups = groups
	}
	return new, nil
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
