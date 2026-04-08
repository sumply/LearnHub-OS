/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"server/internal/adapter/postgres"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/feature/attempt/start_attempt"
	"server/internal/feature/group/create_group"
	"server/internal/feature/quiz/create_quiz"
	"server/internal/feature/subject/create_subject"
	"server/internal/feature/user/create_user"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// prepareCmd represents the prepare command
var prepareCmd = &cobra.Command{
	Use: "prepare",
	RunE: func(cmd *cobra.Command, args []string) error {
		creator := new(config.Env)
		conf := creator.CreatePostgresConnection()
		p, err := postgres.New(conf.CreateOptions())
		if err != nil {
			return err
		}

		createUserUC := create_user.New(
			postgres.NewUser(p),
			postgres.NewTeacher(p),
			postgres.NewStudent(p),
		)

		cmd.Println("Creating teacher...")
		resp, err := createUserUC.CreateUser(cmd.Context(), &SystemIdentity{}, create_user.Input{
			FirstName: "Имя",
			LastName:  "Фамилия",
			Email:     "teacher@example.ru",
			Role:      "teacher",
		})
		if err != nil {
			return err
		}
		teacherID := resp.ID
		cmd.Printf("Teacher has been created (id=%s)!\n", teacherID)

		cmd.Println("Creating student...")
		resp, err = createUserUC.CreateUser(cmd.Context(), &SystemIdentity{}, create_user.Input{
			FirstName: "Имя",
			LastName:  "Фамилия",
			Email:     "student@example.ru",
			Role:      "student",
		})
		if err != nil {
			return err
		}
		studentID := resp.ID
		cmd.Printf("Student has been created (id=%s)!\n", studentID)

		createGroupUC := create_group.New(
			postgres.NewGroup(p),
			postgres.NewUser(p),
		)

		cmd.Println("Creating group...")
		createGroupResp, err := createGroupUC.CreateGroup(cmd.Context(), &SystemIdentity{}, &create_group.Input{
			Name:       "11-A",
			CuratorID:  teacherID,
			StudentIDs: uuid.UUIDs{studentID},
		})
		if err != nil {
			return err
		}
		groupID := createGroupResp.ID
		cmd.Printf("Group has been created (id=%s)!\n", groupID)

		cmd.Println("Creating subject...")
		createSubjectUC := create_subject.New(
			postgres.NewSubject(p),
		)
		createSubjectResp, err := createSubjectUC.CreateSubject(cmd.Context(), &SystemIdentity{}, create_subject.Input{
			Name: "Математика",
		})
		if err != nil {
			return err
		}
		subjectID := createSubjectResp.ID
		cmd.Printf("Subject has been created (id=%s)!\n", subjectID)

		createQuizUC := create_quiz.New(
			postgres.NewQuiz(p),
		)
		cmd.Println("Creating quiz...")
		createQuizResp, err := createQuizUC.CreateQuiz(cmd.Context(), create_quiz.Request{
			Title:     "Тестовый квиз",
			Summary:   "Тест сгенерирован для проверки приложения",
			OwnerID:   teacherID,
			SubjectID: subjectID,
			GroupIDs:  uuid.UUIDs{groupID},
			Deadline: func() *time.Time {
				t := time.Now().Add(time.Hour * 24 * 7).UTC()
				return &t
			}(),
			MaxAttempts: 3,
			Questions: []create_quiz.RequestQuestion{
				{
					Domain: func() domain.IQuestion {
						question, _ := domain.NewSingleQuestion(
							"Сколько будет 2 + 2?",
							"4",
							[]string{"1", "2", "3", "4"},
							10,
						)
						return question
					}(),
				},
				{
					Domain: func() domain.IQuestion {
						question, _ := domain.NewMultipleQuestion(
							"Что является математическими операциями?",
							[]string{"+", "-"},
							[]string{"+", "-", "|", "&&"},
							5,
						)
						return question
					}(),
				},
				{
					Domain: func() domain.IQuestion {
						question, _ := domain.NewNumericQuestion(
							"2 + 2 * 2",
							6,
							2,
						)
						return question
					}(),
				},
			},
		})
		if err != nil {
			return err
		}
		quizID := createQuizResp.ID
		cmd.Printf("Quiz has been created (id=%s)!\n", quizID)

		startAttemptUC := start_attempt.New(
			postgres.NewQuiz(p),
			postgres.NewAttempt(p),
		)

		cmd.Println("Creating attempt...")
		startAttemptResp, err := startAttemptUC.StartAttempt(cmd.Context(), start_attempt.Request{
			QuizID: quizID,
			UserID: studentID,
		})
		if err != nil {
			return err
		}
		attemptID := startAttemptResp.ID
		cmd.Printf("Attempt has been created (id=%s)!\n", attemptID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(prepareCmd)
}
