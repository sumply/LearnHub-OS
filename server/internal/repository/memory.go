package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	progress   map[common.ID]*domain.QuizProgress
	answer     map[common.ID]*domain.Answer
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
		progress:   make(map[common.ID]*domain.QuizProgress),
		answer:     make(map[common.ID]*domain.Answer),
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
	if err := s.joinDepencyAnswer(); err != nil {
		return err
	}
	if err := s.joinDepencyProgress(); err != nil {
		return err
	}
	fmt.Println("Data has been loaded!")
	return nil
}

func (s *Storage) joinDepencyUser() error {
	for _, user := range s.user {
		cred, ok := s.credential[user.ID]
		if !ok {
			return fmt.Errorf("%w: credential_id=%d", ErrDependence, user.ID)
		}
		user.Credential = cred
	}
	return nil
}

func (s *Storage) joinDepencyGroup() error {
	for _, group := range s.group {
		curator, ok := s.user[group.Curator.ID]
		if !ok {
			return fmt.Errorf("%w: curator_id=%d", ErrDependence, group.Curator.ID)
		}
		group.Curator = curator
		for i := range group.Students {
			student, ok := s.user[group.Students[i].ID]
			if !ok {
				return fmt.Errorf("%w: student_id=%d", ErrDependence, group.Students[i].ID)
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
				return fmt.Errorf("%w: option_id=%d", ErrDependence, question.Options[i].ID)
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
			return fmt.Errorf("%w: owner_id=%d", ErrDependence, quiz.Owner.ID)
		}
		quiz.Owner = owner
		subject, ok := s.subject[quiz.Subject.ID]
		if !ok {
			return fmt.Errorf("%w: subject_id=%d", ErrDependence, quiz.Subject.ID)
		}
		quiz.Subject = subject
		for i := range quiz.Groups {
			group, ok := s.group[quiz.Groups[i].ID]
			if !ok {
				return fmt.Errorf("%w: group_id=%d", ErrDependence, quiz.Groups[i].ID)
			}
			quiz.Groups[i] = group
		}
		for i := range quiz.Questions {
			question, ok := s.question[quiz.Questions[i].ID]
			if !ok {
				return fmt.Errorf("%w: question_id=%d", ErrDependence, quiz.Questions[i].ID)
			}
			quiz.Questions[i] = question
		}
	}
	return nil
}

func (s *Storage) joinDepencyProgress() error {
	for _, progress := range s.progress {
		quiz, ok := s.quiz[progress.Quiz.ID]
		if !ok {
			return fmt.Errorf("%w: quiz_id=%d", ErrDependence, progress.Quiz.ID)
		}
		progress.Quiz = quiz
		user, ok := s.user[progress.User.ID]
		if !ok {
			return fmt.Errorf("%w: user_id=%d", ErrDependence, progress.User.ID)
		}
		progress.User = user
		for i, answer := range progress.Answers {
			a, ok := s.answer[answer.ID]
			if !ok {
				return fmt.Errorf("%w: answer_id=%d", ErrDependence, answer.ID)
			}
			progress.Answers[i] = a
		}
	}
	return nil
}

func (s *Storage) joinDepencyAnswer() error {
	for _, answer := range s.answer {
		question, ok := s.question[answer.Question.ID]
		if !ok {
			return fmt.Errorf("%w: question_id=%d", ErrDependence, answer.Question.ID)
		}
		answer.Question = question
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
	if err := f(&s.progress, "progress"); err != nil {
		return err
	}
	if err := f(&s.answer, "answer"); err != nil {
		return err
	}
	return nil
}

func (s *Storage) saveMap(m any, pathName string) error {
	data, err := json.MarshalIndent(m, "\t", "\t")
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	err = os.Mkdir("save", 0777)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("storage: %w", err)
	}
	file, err := os.Create("save/" + pathName + ".json")
	if err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("storage: %w", err)
	}
	return nil
}

func (s *Storage) loadMap(m any, pathName string) error {
	if err := os.Mkdir("save", 0777); err != nil && !os.IsExist(err) {
		return fmt.Errorf("%w: make directory", err)
	}
	file, err := os.Open("save/" + pathName + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create("save/" + pathName + ".json")
			if err != nil {
				return fmt.Errorf("%w: open file", err)
			}
			return file.Close()
		}
		return fmt.Errorf("%w: open file", err)
	}
	if err := json.NewDecoder(file).Decode(m); err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("%w: decoding file info", err)
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
	i := 0
	for _, user := range m.storage.user {
		u := user.Copy()
		u.Credential = nil
		users[i] = u
		i++
	}

	return users, nil
}

func (m *UserMemory) DeleteByID(ctx context.Context, userID common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.NewTracedField("userID", userID),
	)
	log.Debug("Called a deleteByID user repository")

	delete(m.storage.user, userID)

	return nil
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
		return fmt.Errorf("%w: curator_id=%d", ErrDependence, group.Curator.ID)
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

	for _, quiz := range m.storage.quiz {
		for _, group := range quiz.Groups {
			if group.ID == groupID {
				m.storage.saveProgress(quiz)
			}
		}
	}

	return nil
}

func (m *GroupMemory) RemoveStudent(ctx context.Context, groupID common.ID, studentID common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(groupID),
		logger.TraceFieldFromAny(studentID),
	)
	log.Debug("Called a removeStudent repository method")

	group, ok := m.storage.group[groupID]
	if !ok {
		return fmt.Errorf("%w: group_id=%d", ErrNotFound, groupID)
	}

	for i := range group.Students {
		if group.Students[i].ID == studentID {
			group.Students = slices.Delete(group.Students, i, i+1)
			break
		}
	}

	return nil
}

func (m *GroupMemory) GetByID(ctx context.Context, groupID common.ID) (*domain.Group, error) {
	log := logger.FromCtx(ctx).With(
		logger.NewTracedField("groupID", groupID),
	)
	log.Debug("Called a getByID group repository")

	group, ok := m.storage.group[groupID]
	if !ok {
		return nil, fmt.Errorf("%w: group_id=%d", ErrNotFound, groupID)
	}
	return group, nil
}

func (m *GroupMemory) DeleteByID(ctx context.Context, groupID common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.NewTracedField("groupID", groupID),
	)
	log.Debug("Called a DeleteByID group repository")

	group, ok := m.storage.group[groupID]
	if !ok {
		return nil
	}

	if len(group.Students) != 0 {
		return fmt.Errorf("%w: group has students", ErrDependence)
	}

	delete(m.storage.group, groupID)

	return nil
}

func (m *GroupMemory) GetByStudentID(ctx context.Context, studentID common.ID) (*domain.Group, error) {
	log := logger.FromCtx(ctx).With(
		logger.NewTracedField("studentID", studentID),
	)
	log.Debug("Called a getByStudent group memory")

	var found *domain.Group
	for _, group := range m.storage.group {
		ok := slices.ContainsFunc(group.Students, func(u *domain.User) bool {
			if u.ID == studentID {
				return true
			}
			return false
		})
		if ok {
			found = group
			break
		}
	}

	return found, nil
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

	new := quiz.Copy()

	if err := m.prepareQuiz(new); err != nil {
		return err
	}

	m.saveQuiz(new)

	if err := m.storage.saveProgress(new); err != nil {
		return err
	}

	return nil
}

func (m *QuizMemory) prepareQuiz(quiz *domain.Quiz) error {
	if err := m.checkQuizDepensy(quiz); err != nil {
		return err
	}

	m.joinQuizDepensy(quiz)

	return nil
}

func (m *QuizMemory) joinQuizDepensy(quiz *domain.Quiz) {
	quiz.ID = m.storage.nextID()

	owner := m.storage.user[quiz.Owner.ID]
	quiz.Owner = owner

	subject := m.storage.subject[quiz.Subject.ID]
	quiz.Subject = subject

	for i, group := range quiz.Groups {
		g := m.storage.group[group.ID]
		quiz.Groups[i] = g
	}
}

func (m *QuizMemory) checkQuizDepensy(quiz *domain.Quiz) error {
	if _, ok := m.storage.user[quiz.Owner.ID]; !ok {
		return fmt.Errorf("%w: owner_id=%d", ErrDependence, quiz.Owner.ID)
	}

	if _, ok := m.storage.subject[quiz.Subject.ID]; !ok {
		return fmt.Errorf("%w: subject_id=%d", ErrDependence, quiz.Subject.ID)
	}

	for _, group := range quiz.Groups {
		if _, ok := m.storage.group[group.ID]; !ok {
			return fmt.Errorf("%w: group_id=%d", ErrDependence, group.ID)
		}
	}
	return nil
}

func (m *QuizMemory) saveQuiz(quiz *domain.Quiz) {
	quiz.Questions = m.saveQuestions(quiz.Questions)
	m.storage.quiz[quiz.ID] = quiz
}

func (s *Storage) saveProgress(quiz *domain.Quiz) error {
	for _, group := range quiz.Groups {
		for _, student := range group.Students {
			if !s.progressIsExists(student.ID, quiz.ID) {
				progress := domain.NewQuizProgress(student.ID, quiz)
				progress.ID = s.nextID()

				s.saveAnswers(progress.Answers)

				if err := s.joinProgressDependency(progress); err != nil {
					return err
				}
				s.progress[progress.ID] = progress
			}
		}
	}
	return nil
}

func (s *Storage) progressIsExists(studentID, quizID common.ID) bool {
	for _, saved := range s.progress {
		if saved.Quiz.ID == quizID && saved.User.ID == studentID {
			return true
		}
	}
	return false
}

func (s *Storage) saveAnswers(answers []*domain.Answer) {
	for _, answer := range answers {
		answer.ID = s.nextID()
		s.answer[answer.ID] = answer
	}
}

func (s *Storage) joinProgressDependency(progress *domain.QuizProgress) error {
	user, ok := s.user[progress.User.ID]
	if !ok {
		return fmt.Errorf("%w: user_id=%d", ErrDependence, progress.User.ID)
	}
	quiz, ok := s.quiz[progress.Quiz.ID]
	if !ok {
		return fmt.Errorf("%w: quiz_id=%d", ErrDependence, progress.Quiz.ID)
	}
	progress.User = user
	progress.Quiz = quiz
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

func (m *QuizMemory) GetWithFilter(ctx context.Context, filter *QuizFilter) ([]*domain.Quiz, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(filter),
	)
	log.Debug("Called a find quizMemory method")

	var quizzes []*domain.Quiz
	for _, saved := range m.storage.quiz {
		if satisfiesQuizFilter(filter, saved) {
			quizzes = append(quizzes, saved.Copy())
		}
	}
	return quizzes, nil
}

func satisfiesQuizFilter(filter *QuizFilter, quiz *domain.Quiz) bool {
	if filter == nil {
		return true
	}
	if quiz.IsForEveryone {
		return true
	}
	if filter.OwnerID != nil && *filter.OwnerID != quiz.Owner.ID {
		return false
	}
	if filter.Group != nil {
		ok := slices.ContainsFunc(quiz.Groups, func(g *domain.Group) bool {
			return satisfiesGroupFilter(filter.Group, g)
		})
		if !ok {
			return false
		}
	}
	return true
}

func satisfiesGroupFilter(filter *GroupFilter, group *domain.Group) bool {
	if filter == nil {
		return true
	}
	if filter.StudentID != nil {
		ok := slices.ContainsFunc(group.Students, func(u *domain.User) bool {
			return u.ID == *filter.StudentID
		})
		if !ok {
			return false
		}
	}
	return true
}

func (m *QuizMemory) GetByID(ctx context.Context, quizID common.ID) (*domain.Quiz, error) {
	log := logger.FromCtx(ctx).With(
		logger.NewTracedField("quizID", quizID),
	)
	log.Debug("Called findWithFilter")

	quiz, ok := m.storage.quiz[quizID]
	if !ok {
		return nil, fmt.Errorf("%w: quiz_id=%d", ErrNotFound, quizID)
	}

	return quiz.Copy(), nil
}

func (m *QuizMemory) Delete(ctx context.Context, quizID common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.NewTracedField("quizID", quizID),
	)
	log.Debug("Called delete")

	return m.deleteQuiz(quizID)
}

func (m *QuizMemory) deleteQuiz(quizID common.ID) error {
	quiz, ok := m.storage.quiz[quizID]
	if !ok {
		return fmt.Errorf("%w: quiz_id=%d", ErrNotFound, quizID)
	}
	for _, question := range quiz.Questions {
		m.deleteQuestion(question)
	}
	delete(m.storage.quiz, quizID)
	return nil
}

func (m *QuizMemory) deleteQuestion(question *domain.Question) {
	if _, ok := m.storage.question[question.ID]; !ok {
		return
	}
	for _, option := range question.Options {
		m.deleteOption(option)
	}
	delete(m.storage.question, question.ID)
}

func (m *QuizMemory) deleteOption(option *domain.Option) {
	if _, ok := m.storage.option[option.ID]; !ok {
		return
	}
	delete(m.storage.option, option.ID)
}

type ProgressMemory struct {
	storage *Storage
}

func NewProgressMemory(storage *Storage) *ProgressMemory {
	return &ProgressMemory{
		storage: storage,
	}
}

func (m *ProgressMemory) GetByID(ctx context.Context, progressID common.ID) (*domain.QuizProgress, error) {
	log := logger.FromCtx(ctx).With(
		logger.NewTracedField("progress_id", progressID),
	)
	log.Debug("Called a getByID repository")

	progress, ok := m.storage.progress[progressID]
	if !ok {
		return nil, fmt.Errorf("%w: progress_id=%d", ErrNotFound, progressID)
	}
	return progress, nil
}

func (m *ProgressMemory) UpdateAnswer(ctx context.Context, answer *domain.Answer) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(answer),
	)
	log.Debug("Called a updateAnswer repository")

	if _, ok := m.storage.answer[answer.ID]; !ok {
		return fmt.Errorf("%w: answer_id=%d", ErrNotFound, answer.ID)
	}

	m.storage.answer[answer.ID] = answer
	return nil
}

func (m *ProgressMemory) Get(ctx context.Context, filter *ProgressFilter) ([]*domain.QuizProgress, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(filter),
	)
	log.Debug("Called a get repository")

	var progresses []*domain.QuizProgress
	for _, progress := range m.storage.progress {
		if satisfiesProgressFilter(filter, progress) {
			progresses = append(progresses, progress.Copy())
		}
	}

	log.With(
		logger.NewTracedField("result", progresses),
		logger.NewTracedField("data", m.storage.progress),
	).Debug("Result function")
	return progresses, nil
}

func satisfiesProgressFilter(filter *ProgressFilter, progress *domain.QuizProgress) bool {
	if filter == nil {
		return true
	}
	if filter.UserID != nil && *filter.UserID != progress.User.ID {
		return false
	}
	if filter.Quiz != nil && !satisfiesQuizFilter(filter.Quiz, progress.Quiz) {
		return false
	}
	return true
}

func (m *ProgressMemory) Update(progress *domain.QuizProgress) error {
	if _, ok := m.storage.progress[progress.ID]; !ok {
		return fmt.Errorf("%w: progress_id=%d", ErrNotFound, progress.ID)
	}
	m.storage.progress[progress.ID] = progress
	return nil
}
