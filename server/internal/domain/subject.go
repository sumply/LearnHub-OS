package domain

type SubjectID uint64

type SubjectName string

type Subject struct {
	ID   SubjectID
	Name SubjectName
}
