/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"server/internal/adapter/postgres"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/feature/user/create_user"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type SystemIdentity struct{}

func (r *SystemIdentity) ID() uuid.UUID {
	return uuid.Nil
}

func (r *SystemIdentity) Role() domain.UserRole {
	return domain.RoleAdmin
}

// superuserCmd represents the superuser command
var superuserCmd = &cobra.Command{
	Use:   "user",
	Short: "Creating new user.",
	RunE: func(cmd *cobra.Command, args []string) error {
		creator := new(config.Env)
		pOpt := creator.CreatePostgresConnection()
		p, err := postgres.New(pOpt.CreateOptions())
		if err != nil {
			return err
		}
		uc := create_user.New(
			postgres.NewUser(p),
			postgres.NewTeacher(p),
			postgres.NewStudent(p),
		)
		input, err := makeCreateUserInput(cmd)
		if err != nil {
			return err
		}
		_, err = uc.CreateUser(context.Background(), new(SystemIdentity), input)
		return err
	},
}

func makeCreateUserInput(cmd *cobra.Command) (create_user.Input, error) {
	fn, err := cmd.Flags().GetString("first_name")
	if err != nil {
		return create_user.Input{}, err
	}
	ln, err := cmd.Flags().GetString("last_name")
	if err != nil {
		return create_user.Input{}, err
	}
	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return create_user.Input{}, err
	}
	role, err := cmd.Flags().GetString("role")
	return create_user.Input{
		FirstName: fn,
		LastName:  ln,
		Email:     email,
		Role:      role,
	}, nil
}

func init() {
	rootCmd.AddCommand(superuserCmd)

	superuserCmd.Flags().String("first_name", "default", "user's first name")
	superuserCmd.Flags().String("last_name", "default", "user's last name")
	superuserCmd.Flags().String("email", "", "user's email address")
	superuserCmd.Flags().String("role", "admin", "user's role")
	superuserCmd.MarkFlagRequired("email")
}
