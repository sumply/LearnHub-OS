package repository

import (
	"context"
	"fmt"
	"server/internal/common"
	"server/internal/domain"
	"server/internal/logger"
	"time"
)

type Storage struct {
	credential map[common.ID]*domain.Credential
	user       map[common.ID]*domain.User
	group      map[common.ID]*domain.Group
	option     map[common.ID]*domain.Option
	question   map[common.ID]*domain.Question
	quiz       map[common.ID]*domain.Quiz
	subject    map[common.ID]*domain.Subject
}

func NewStorage() *Storage {
	return &Storage{
		credential: make(map[common.ID]*domain.Credential),
		user:       make(map[common.ID]*domain.User),
		group:      make(map[common.ID]*domain.Group),
		option:     make(map[common.ID]*domain.Option),
		question:   make(map[common.ID]*domain.Question),
		quiz:       make(map[common.ID]*domain.Quiz),
		subject:    make(map[common.ID]*domain.Subject),
	}
}

var id common.ID

func nextID() common.ID {
	id++
	return id
}

func PrepareStorage(s *Storage) {
	s.credential[0] = &domain.Credential{
		ID:        0,
		Login:     "admin",
		PwdHashed: "admin",
		Email:     "admin@admin.ru",
	}
	s.user[0] = &domain.User{
		ID:         0,
		FirstName:  "admin",
		LastName:   "admin",
		MiddleName: "admin",
		Role:       domain.UserRoot,
		Credential: s.credential[0],
		CreatedAt:  time.Now(),
	}
}

type UserMemory struct {
	storage *Storage
}

func NewUserMemory(s *Storage) *UserMemory {
	return &UserMemory{
		storage: s,
	}
}

func (m *UserMemory) Save(ctx context.Context, user *domain.User) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(user),
	)
	log.Debug("Called a save repository method")
	user.ID = nextID()
	user.Credential.ID = user.ID
	for _, saved := range m.storage.credential {
		if saved.Login == user.Credential.Login {
			return fmt.Errorf("%w: login=%s", ErrCollision, user.Credential.Login)
		}
	}
	m.storage.credential[user.Credential.ID] = user.Credential
	if _, ok := m.storage.user[common.ID(user.ID)]; ok {
		return fmt.Errorf("%w: user_id=%d", ErrCollision, user.ID)
	}
	m.storage.user[user.ID] = user
	return nil
}

func (m *UserMemory) GetByID(ctx context.Context, id common.ID) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(id),
	)
	log.Debug("Called a save repository method")
	user, ok := m.storage.user[id]
	if !ok {
		return nil, fmt.Errorf("%w: user_id=%d", ErrNotFound, id)
	}
	return user, nil
}

func (m *UserMemory) GetByLogin(ctx context.Context, login domain.Login) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(login),
	)
	log.Debug("Called a getByLogin repository method")
	for _, user := range m.storage.user {
		if user.Credential.Login == login {
			return user, nil
		}
	}
	return nil, fmt.Errorf("%w: login='%s'", ErrNotFound, login)
}

func (m *UserMemory) GetAll(ctx context.Context) ([]*domain.User, error) {
	log := logger.FromCtx(ctx)
	log.Debug("Called a getAll repository method")

	if len(m.storage.user) == 0 {
		return nil, ErrNotFound
	}

	users := make([]*domain.User, len(m.storage.user))
	for i, user := range m.storage.user {
		u := *user
		u.Credential = nil
		users[i] = &u
	}

	return users, nil
}

type GroupMemory struct {
	storage *Storage
}

func NewGroupMemory(s *Storage) *GroupMemory {
	return &GroupMemory{
		storage: s,
	}
}

func (m *GroupMemory) Save(ctx context.Context, group *domain.Group) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(group),
	)
	log.Debug("Called a save repotiroy method")

	if _, ok := m.storage.group[group.ID]; ok {
		return fmt.Errorf("%w: group_id=%d", ErrCollision, group.ID)
	}

	for _, g := range m.storage.group {
		if g.Name == group.Name {
			return fmt.Errorf(`%w: name="%s"`, ErrCollision, group.Name)
		}
	}

	curator, ok := m.storage.user[group.Curator.ID]
	if !ok {
		return fmt.Errorf("%w: curator_id=%d", ErrDependensy, group.Curator.ID)
	}

	new := *group
	new.ID = nextID()
	new.Curator = curator

	m.storage.group[new.ID] = &new
	return nil
}

func (m *GroupMemory) GetAll(ctx context.Context) ([]*domain.Group, error) {
	log := logger.FromCtx(ctx)
	log.Debug("Called a getAll repository method")

	if len(m.storage.group) == 0 {
		return nil, ErrNotFound
	}

	groups := make([]*domain.Group, len(m.storage.group))
	i := 0
	for _, g := range m.storage.group {
		new := *g
		log.With(logger.TraceFieldFromAny(new)).Debug("new")
		groups[i] = &new
		i++
	}
	return groups, nil
}

func (m *GroupMemory) AddStudent(ctx context.Context, groupID common.ID, userIDs []common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(groupID),
		logger.TraceFieldFromAny(userIDs),
	)
	log.Debug("Called a addStudent repository method")

	group, ok := m.storage.group[groupID]
	if !ok {
		return fmt.Errorf("%w: group_id=%d", ErrNotFound, groupID)
	}

	if len(userIDs) == 0 {
		return fmt.Errorf("%w: length a userIDs is 0", ErrInvalid)
	}

	uniqueIDs := make(map[common.ID]bool)
	for _, student := range group.Students {
		uniqueIDs[student.ID] = true
	}

	for _, id := range userIDs {
		if _, ok := uniqueIDs[id]; ok {
			return fmt.Errorf("%w: user_id=%d", ErrCollision, id)
		}
		if _, ok := m.storage.user[id]; !ok {
			return fmt.Errorf("%w: user_id=%d", ErrNotFound, id)
		}
		uniqueIDs[id] = true
	}

	newStudents := make([]*domain.User, len(uniqueIDs))
	i := 0
	for id := range uniqueIDs {
		newStudents[i] = m.storage.user[id]
		i++
	}

	group.Students = newStudents
	return nil
}

func (m *GroupMemory) RemoveStudent(ctx context.Context, groupID common.ID, userIDs []common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(groupID),
		logger.TraceFieldFromAny(userIDs),
	)
	log.Debug("Called a removeStudent repository method")

	group, ok := m.storage.group[groupID]
	if !ok {
		return fmt.Errorf("%w: group_id=%d", ErrNotFound, groupID)
	}

	setUserIDs := make(map[common.ID]bool)
	for _, user := range group.Students {
		setUserIDs[user.ID] = true
	}

	for _, id := range userIDs {
		if _, ok := setUserIDs[id]; !ok {
			return fmt.Errorf("%w: user_id=%d", ErrNotFound, id)
		}
		delete(setUserIDs, id)
	}

	newStudents := make([]*domain.User, len(setUserIDs))
	i := 0
	for id := range setUserIDs {
		newStudents[i] = m.storage.user[id]
	}

	group.Students = newStudents
	return nil
}

type SubjectMemory struct {
	storage *Storage
}

func NewSubjectMemory(s *Storage) *SubjectMemory {
	return &SubjectMemory{
		storage: s,
	}
}

func (m *SubjectMemory) Save(ctx context.Context, subject *domain.Subject) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(subject),
	)
	log.Debug("Called a save SubjectMemory method")

	for _, saved := range m.storage.subject {
		if saved.Name == subject.Name {
			return fmt.Errorf("%w: subject_name='%s'", ErrCollision, subject.Name)
		}
	}

	new := *subject
	new.ID = nextID()
	m.storage.subject[new.ID] = &new

	return nil
}

func (m *SubjectMemory) GetAll(ctx context.Context) ([]*domain.Subject, error) {
	log := logger.FromCtx(ctx)
	log.Debug("Called a getAll SubjectMemory method")

	if len(m.storage.subject) == 0 {
		return nil, ErrNotFound
	}

	subjects := make([]*domain.Subject, len(m.storage.subject))
	i := 0
	for _, saved := range m.storage.subject {
		new := *saved
		subjects[i] = &new
		i++
	}
	return subjects, nil
}

type QuizMemory struct {
	storage *Storage
}

func NewQuizMemory(s *Storage) *QuizMemory {
	return &QuizMemory{
		storage: s,
	}
}

func (m *QuizMemory) Save(ctx context.Context, quiz *domain.Quiz) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(quiz),
	)
	log.Debug("Called quizMemory method save")

	if _, ok := m.storage.user[quiz.Owner.ID]; !ok {
		return fmt.Errorf("%w: owner_id=%d", ErrDependensy, quiz.Owner.ID)
	}

	if _, ok := m.storage.subject[quiz.Subject.ID]; !ok {
		return fmt.Errorf("%w: subject_id=%d", ErrDependensy, quiz.Subject.ID)
	}

	for _, group := range quiz.Groups {
		if _, ok := m.storage.group[group.ID]; !ok {
			return fmt.Errorf("%w: group_id=%d", ErrDependensy, group.ID)
		}
	}

	new := *quiz
	new.ID = nextID()
	new.Questions = m.saveQuestions(new.Questions)

	m.storage.quiz[new.ID] = &new

	return nil
}

func (m *QuizMemory) saveQuestions(questions []*domain.Question) []*domain.Question {
	copied := make([]*domain.Question, len(questions))
	for i, quiestion := range questions {
		new := *quiestion
		copied[i] = &new
	}

	for _, question := range copied {
		question.ID = nextID()
		m.saveOptions(question.Options)
		m.storage.question[question.ID] = question
	}

	return copied
}

func (m *QuizMemory) saveOptions(options []*domain.Option) []*domain.Option {
	copied := make([]*domain.Option, len(options))
	for i, opt := range options {
		new := *opt
		copied[i] = &new
	}

	for _, opt := range copied {
		opt.ID = common.ID(nextID())
		m.storage.option[opt.ID] = opt
	}

	return copied
}
