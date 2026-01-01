package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"server/internal/common"
	"server/internal/domain"
	"server/internal/logger"
	"slices"
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
	id         common.ID
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

func (s *Storage) Load() error {
	if err := s.mapFields(s.loadMap); err != nil {
		return err
	}
	if err := s.joinDepencyUser(); err != nil {
		return err
	}
	if err := s.joinDepencyGroup(); err != nil {
		return err
	}
	if err := s.joinDepencyQuestion(); err != nil {
		return err
	}
	if err := s.joinDepencyQuiz(); err != nil {
		return err
	}
	fmt.Println("Data has been loaded!")
	return nil
}

func (s *Storage) joinDepencyUser() error {
	for _, user := range s.user {
		cred, ok := s.credential[user.ID]
		if !ok {
			return fmt.Errorf("%w: credential_id=%d", ErrDependensy, user.ID)
		}
		user.Credential = cred
	}
	return nil
}

func (s *Storage) joinDepencyGroup() error {
	for _, group := range s.group {
		curator, ok := s.user[group.Curator.ID]
		if !ok {
			return fmt.Errorf("%w: curator_id=%d", ErrDependensy, group.Curator.ID)
		}
		group.Curator = curator
		for i := range group.Students {
			student, ok := s.user[group.Students[i].ID]
			if !ok {
				return fmt.Errorf("%w: student_id=%d", ErrDependensy, group.Students[i].ID)
			}
			group.Students[i] = student
		}
	}
	return nil
}

func (s *Storage) joinDepencyQuestion() error {
	for _, question := range s.question {
		for i := range question.Options {
			option, ok := s.option[question.Options[i].ID]
			if !ok {
				return fmt.Errorf("%w: option_id=%d", ErrDependensy, question.Options[i].ID)
			}
			question.Options[i] = option
		}
	}
	return nil
}

func (s *Storage) joinDepencyQuiz() error {
	for _, quiz := range s.quiz {
		owner, ok := s.user[quiz.Owner.ID]
		if !ok {
			return fmt.Errorf("%w: owner_id=%d", ErrDependensy, quiz.Owner.ID)
		}
		quiz.Owner = owner
		subject, ok := s.subject[quiz.Subject.ID]
		if !ok {
			return fmt.Errorf("%w: subject_id=%d", ErrDependensy, quiz.Subject.ID)
		}
		quiz.Subject = subject
		for i := range quiz.Groups {
			group, ok := s.group[quiz.Groups[i].ID]
			if !ok {
				return fmt.Errorf("%w: group_id=%d", ErrDependensy, quiz.Groups[i].ID)
			}
			quiz.Groups[i] = group
		}
		for i := range quiz.Questions {
			question, ok := s.question[quiz.Questions[i].ID]
			if !ok {
				return fmt.Errorf("%w: question_id=%d", ErrDependensy, quiz.Questions[i].ID)
			}
			quiz.Questions[i] = question
		}
	}
	return nil
}

func (s *Storage) Save() error {
	return s.mapFields(s.saveMap)
}

func (s *Storage) mapFields(f func(any, string) error) error {
	if err := f(&s.credential, "credential"); err != nil {
		return err
	}
	if err := f(&s.user, "user"); err != nil {
		return err
	}
	if len(s.user) == 0 {
		PrepareStorage(s)
	}
	if err := f(&s.group, "group"); err != nil {
		return err
	}
	if err := f(&s.option, "option"); err != nil {
		return err
	}
	if err := f(&s.question, "question"); err != nil {
		return err
	}
	if err := f(&s.quiz, "quiz"); err != nil {
		return err
	}
	if err := f(&s.subject, "subject"); err != nil {
		return err
	}
	if err := f(&s.id, "id"); err != nil {
		return err
	}
	return nil
}

func (s *Storage) saveMap(m any, pathName string) error {
	data, err := json.MarshalIndent(m, "\t", "\t")
	if err != nil {
		return err
	}
	err = os.Mkdir("save", 0777)
	if err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.Create("save/" + pathName + ".json")
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}

func (s *Storage) loadMap(m any, pathName string) error {
	if err := os.Mkdir("save", 0777); err != nil && !os.IsExist(err) {
		return err
	}
	file, err := os.Open("save/" + pathName + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create("save/" + pathName + ".json")
			if err != nil {
				return err
			}
			return file.Close()
		}
		return err
	}
	if err := json.NewDecoder(file).Decode(m); err != nil {
		return err
	}
	return nil
}

func (s *Storage) nextID() common.ID {
	s.id++
	return s.id
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
	user.ID = m.storage.nextID()
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
	new.ID = m.storage.nextID()
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
	new.ID = m.storage.nextID()
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

	log.Debug("copying a quiz")
	new := quiz.Copy()

	log.Debug("set next id")
	new.ID = m.storage.nextID()

	log.Debug("finding a owner")
	owner, ok := m.storage.user[quiz.Owner.ID]
	if !ok {
		return fmt.Errorf("%w: owner_id=%d", ErrDependensy, quiz.Owner.ID)
	}
	new.Owner = owner

	log.Debug("finding a subject")
	subject, ok := m.storage.subject[quiz.Subject.ID]
	if !ok {
		return fmt.Errorf("%w: subject_id=%d", ErrDependensy, quiz.Subject.ID)
	}
	new.Subject = subject

	log.Debug("finding groups")
	for i, group := range quiz.Groups {
		g, ok := m.storage.group[group.ID]
		if !ok {
			return fmt.Errorf("%w: group_id=%d", ErrDependensy, group.ID)
		}
		new.Groups[i] = g
	}

	log.Debug("saving questions")
	new.Questions = m.saveQuestions(new.Questions)

	log.Debug("saving a quiz")
	m.storage.quiz[new.ID] = new

	return nil
}

func (m *QuizMemory) saveQuestions(questions []*domain.Question) []*domain.Question {
	copied := make([]*domain.Question, len(questions))
	for i, quiestion := range questions {
		copied[i] = quiestion.Copy()
	}

	for _, question := range copied {
		question.ID = m.storage.nextID()
		question.Options = m.saveOptions(question.Options)
		m.storage.question[question.ID] = question
	}

	return copied
}

func (m *QuizMemory) saveOptions(options []*domain.Option) []*domain.Option {
	copied := make([]*domain.Option, len(options))
	for i, opt := range options {
		copied[i] = opt.Copy()
	}

	for _, opt := range copied {
		opt.ID = m.storage.nextID()
		m.storage.option[opt.ID] = opt
	}

	return copied
}

func (m *QuizMemory) GetAll(ctx context.Context) ([]*domain.Quiz, error) {
	log := logger.FromCtx(ctx)
	log.Debug("Called a getAll quizMemory method")

	if len(m.storage.quiz) == 0 {
		return nil, ErrNotFound
	}

	quizzes := make([]*domain.Quiz, len(m.storage.quiz))
	i := 0
	for _, saved := range m.storage.quiz {
		copied := saved.Copy()
		copied.Questions = nil
		quizzes[i] = copied
		i++
	}

	return quizzes, nil
}

func (m *QuizMemory) Find(ctx context.Context, filter *QuizFilter) ([]*domain.Quiz, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(filter),
	)
	log.Debug("Called a find quizMemory method")

	var quizzes []*domain.Quiz
	for _, saved := range m.storage.quiz {
		if m.satisfiesFilter(filter, saved) {
			quizzes = append(quizzes, saved.Copy())
		}
	}
	return quizzes, nil
}

func (m *QuizMemory) satisfiesFilter(filter *QuizFilter, quiz *domain.Quiz) bool {
	if quiz.IsForEveryone {
		return true
	}
	if filter.OwnerID != nil && *filter.OwnerID != quiz.Owner.ID {
		return false
	}
	if filter.Group != nil && !m.satisfiesGroupFilter(filter.Group, quiz.Groups) {
		return false
	}
	return true
}

func (m *QuizMemory) satisfiesGroupFilter(filter *GroupFilter, groups []*domain.Group) bool {
	if filter.StudentID == nil {
		return true
	}
	if len(groups) == 0 {
		return false
	}
	for _, group := range groups {
		ok := slices.ContainsFunc(group.Students, func(u *domain.User) bool {
			return u.ID == *filter.StudentID
		})
		if ok {
			return true
		}
	}
	return false
}
