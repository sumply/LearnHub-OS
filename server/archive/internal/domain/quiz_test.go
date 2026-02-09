package domain

import (
	"server/archive/internal/common"
	"testing"
)

func TestNewOption(t *testing.T) {
	t.Log("creating a valid option")
	currOpt, err := NewOption("currect", false)
	if err != nil {
		t.Errorf("unexcepted error: %s", err.Error())
	}
	if currOpt.IsCorrect != false || currOpt.Text != "currect" {
		t.Errorf("option fields is not valid: %v", currOpt)
	}
	t.Log("creating a invalid option")
	_, err = NewOption("", false)
	if err == nil {
		t.Errorf("excepted error validation")
	}
}

func TestNewQuestion(t *testing.T) {
	options := []*Option{
		{
			ID:        1,
			Text:      "false",
			IsCorrect: false,
		},
		{
			ID:        2,
			Text:      "true",
			IsCorrect: true,
		},
	}
	t.Logf("creating a question")
	valid, err := NewQuestion("question", options)
	if err != nil {
		t.Errorf("unexcepted error: %v", err)
	}
	t.Logf("creating a question without options")
	valid, err = NewQuestion("question", []*Option{})
	if err != nil {
		t.Errorf("unexcepted error: %v", err)
	}
	if !valid.IsWritten {
		t.Errorf("a field IsWritten is false")
	}
	t.Log("creating a question with a invalid text")
	valid, err = NewQuestion("", options)
	if err == nil {
		t.Errorf("excepted a empty error: %v", valid)
	}
	t.Log("creating a question without a current option")
	options = []*Option{
		{
			ID:        1,
			Text:      "false",
			IsCorrect: false,
		},
		{
			ID:        2,
			Text:      "true",
			IsCorrect: false,
		},
	}
	_, err = NewQuestion("question", options)
	if err == nil {
		t.Logf("a question has been created without a current option")
	}
}

func TestNewQuiz(t *testing.T) {
	questions := []*Question{
		{
			ID:        1,
			Text:      "Question",
			IsWritten: false,
			Options: []*Option{
				{
					ID:        1,
					Text:      "Option",
					IsCorrect: false,
				},
				{
					ID:        1,
					Text:      "Option",
					IsCorrect: true,
				},
			},
		},
	}
	t.Logf("creating a quiz without groupIDs")
	quiz, err := NewQuiz("quiz", "summary", questions, 1, 1, []common.ID{})
	if err != nil {
		t.Errorf("unexcepted error: %v", err)
	}
	if !quiz.IsForEveryone {
		t.Errorf("a field IsForEveryone is false")
	}

	t.Logf("creating a quiz with groupIDs")
	quiz, err = NewQuiz("quiz", "summary", questions, 1, 1, []common.ID{1, 2, 3})
	if err != nil {
		t.Errorf("unexcepted error: %v", err)
	}
	if quiz.IsForEveryone {
		t.Errorf("a field IsForEveryone is true")
	}

	t.Logf("creating a quiz with invalid title")
	_, err = NewQuiz("", "summary", questions, 1, 1, []common.ID{1, 2, 3})
	if err == nil {
		t.Errorf("excepted a error")
	}

	t.Logf("creating a quiz with invalid summary")
	_, err = NewQuiz("title", "", questions, 1, 1, []common.ID{1, 2, 3})
	if err == nil {
		t.Errorf("excepted a error")
	}

	t.Logf("creating a quiz with invalid questions")
	_, err = NewQuiz("title", "summary", nil, 1, 1, []common.ID{1, 2, 3})
	if err == nil {
		t.Errorf("excepted a error")
	}
}
