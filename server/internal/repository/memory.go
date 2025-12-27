package repository

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/logger"
	"time"
)

var credStorage = make(map[domain.Login]*domain.Credential)
var userStorage = make(map[domain.UserID]*domain.User)
var specStorage = make(map[domain.SpecialityID]*domain.Speciality)
var groupStorage = make(map[domain.GroupID]*domain.Group)

var id uint64

func nextID() uint64 {
	id++
	return id
}

func curID() uint64 {
	return id
}

func init() {
	PrepareMemory()
}

func PrepareMemory() {
	now := time.Now().UTC()
	credStorage = map[domain.Login]*domain.Credential{
		"root": {
			Login:     "root",
			PwdHashed: domain.HashPassword("root"),
			Email:     "root@root.root",
		},
		"admin": {
			Login:     "admin",
			PwdHashed: domain.HashPassword("admin"),
			Email:     "admin@admin.admin",
		},
		"teacher": {
			Login:     "teacher",
			PwdHashed: domain.HashPassword("teacher"),
			Email:     "teacher@teacher.teacher",
		},
		"student": {
			Login:     "student",
			PwdHashed: domain.HashPassword("student"),
			Email:     "student@student.student",
		},
	}
	userStorage = map[domain.UserID]*domain.User{
		domain.UserID(nextID()): {
			ID:         domain.UserID(curID()),
			FirstName:  "root",
			LastName:   "root",
			MiddleName: "root",
			Role:       domain.UserRoot,
			Credential: credStorage["root"],
			CreatedAt:  now,
		},
		domain.UserID(nextID()): {
			ID:         domain.UserID(curID()),
			FirstName:  "admin",
			LastName:   "admin",
			MiddleName: "admin",
			Role:       domain.UserAdmin,
			Credential: credStorage["admin"],
			CreatedAt:  now,
		},
		domain.UserID(nextID()): {
			ID:         domain.UserID(curID()),
			FirstName:  "teacher",
			LastName:   "teacher",
			MiddleName: "teacher",
			Role:       domain.UserTeacher,
			Credential: credStorage["teacher"],
			CreatedAt:  now,
		},
		domain.UserID(nextID()): {
			ID:         domain.UserID(curID()),
			FirstName:  "student",
			LastName:   "student",
			MiddleName: "student",
			Role:       domain.UserStudent,
			Credential: credStorage["student"],
			CreatedAt:  now,
		},
	}
	specStorage = map[domain.SpecialityID]*domain.Speciality{
		domain.SpecialityID(nextID()): {
			ID:   domain.SpecialityID(curID()),
			Name: "Common",
		},
		domain.SpecialityID(nextID()): {
			ID:   domain.SpecialityID(curID()),
			Name: "English",
		},
		domain.SpecialityID(nextID()): {
			ID:   domain.SpecialityID(curID()),
			Name: "Math",
		},
	}
	groupStorage = map[domain.GroupID]*domain.Group{
		domain.GroupID(nextID()): {
			ID:         domain.GroupID(curID()),
			Name:       "9-A",
			Curator:    userStorage[3],
			Speciality: specStorage[5],
			Students:   []*domain.User{userStorage[4]},
		},
		domain.GroupID(nextID()): {
			ID:         domain.GroupID(curID()),
			Name:       "10-A",
			Curator:    userStorage[2],
			Speciality: specStorage[6],
			Students:   []*domain.User{},
		},
		domain.GroupID(nextID()): {
			ID:         domain.GroupID(curID()),
			Name:       "11-A",
			Curator:    userStorage[3],
			Speciality: specStorage[7],
			Students:   []*domain.User{userStorage[1]},
		},
	}
}

type UserMemory struct{}

func NewUserMemory() User {
	return &UserMemory{}
}

func (u *UserMemory) Save(ctx context.Context, user *domain.User) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(user),
	)
	log.Debug("Called a save repository method")
	user.ID = domain.UserID(nextID())
	if _, ok := credStorage[user.Credential.Login]; ok {
		return fmt.Errorf("%w: login=%s", ErrCollision, user.Credential.Login)
	}
	credStorage[user.Credential.Login] = user.Credential
	if _, ok := userStorage[user.ID]; ok {
		return fmt.Errorf("%w: user_id=%d", ErrCollision, user.ID)
	}
	userStorage[user.ID] = user
	return nil
}

func (u *UserMemory) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(id),
	)
	log.Debug("Called a save repository method")
	user, ok := userStorage[id]
	if !ok {
		return nil, fmt.Errorf("%w: user_id=%d", ErrNotFound, id)
	}
	return user, nil
}
func (u *UserMemory) GetByLogin(ctx context.Context, login domain.Login) (*domain.User, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(login),
	)
	log.Debug("Called a getByLogin repository method")
	for _, user := range userStorage {
		if user.Credential.Login == login {
			return user, nil
		}
	}
	return nil, fmt.Errorf(`%w: login="%s"`, ErrNotFound, login)
}
func (u *UserMemory) GetAll(ctx context.Context) ([]*domain.User, error) {
	log := logger.FromCtx(ctx)
	log.Debug("Called a getAll repository method")

	if len(userStorage) == 0 {
		return nil, ErrNotFound
	}

	users := make([]*domain.User, len(userStorage))
	for i, user := range userStorage {
		u := *user
		u.Credential = nil
		users[i] = &u
	}

	return users, nil
}

type SpecialityMemory struct{}

func (s *SpecialityMemory) Save(ctx context.Context, spec *domain.Speciality) error {
	spec.ID = domain.SpecialityID(nextID())

	log := logger.FromCtx(ctx)
	log.Debug("Called a save repository method")

	if _, ok := specStorage[spec.ID]; ok {
		return fmt.Errorf("%w: speciality_id=%d", ErrCollision, spec.ID)
	}

	for _, specRange := range specStorage {
		if specRange.Name == spec.Name {
			return fmt.Errorf(`%w: name="%s"`, ErrCollision, spec.Name)
		}
	}

	specStorage[spec.ID] = spec
	return nil
}

func (s *SpecialityMemory) GetAll(ctx context.Context) ([]*domain.Speciality, error) {
	log := logger.FromCtx(ctx)
	log.Debug("Called a getAll repository method")

	if len(specStorage) == 0 {
		return nil, ErrNotFound
	}

	domains := make([]*domain.Speciality, len(specStorage))
	i := 0
	for _, domain := range specStorage {
		new := *domain
		domains[i] = &new
		i++
	}
	return domains, nil
}

type GroupMemory struct{}

func (g *GroupMemory) Save(ctx context.Context, group *domain.Group) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(group),
	)
	log.Debug("Called a save repotiroy method")

	if _, ok := groupStorage[group.ID]; ok {
		return fmt.Errorf("%w: group_id=%d", ErrCollision, group.ID)
	}

	for _, g := range groupStorage {
		if g.Name == group.Name {
			return fmt.Errorf(`%w: name="%s"`, ErrCollision, group.Name)
		}
	}

	curator, ok := userStorage[group.Curator.ID]
	if !ok {
		return fmt.Errorf("%w: curator_id=%d", ErrDependensy, group.Curator.ID)
	}

	speciality, ok := specStorage[group.Speciality.ID]
	if !ok {
		return fmt.Errorf("%w: speciality_id=%d", ErrDependensy, group.Speciality.ID)
	}

	new := *group
	new.ID = domain.GroupID(nextID())
	new.Curator = curator
	new.Speciality = speciality

	groupStorage[new.ID] = &new
	return nil
}
func (g *GroupMemory) GetAll(ctx context.Context) ([]*domain.Group, error) {
	log := logger.FromCtx(ctx)
	log.Debug("Called a getAll repository method")

	if len(groupStorage) == 0 {
		return nil, ErrNotFound
	}

	groups := make([]*domain.Group, len(groupStorage))
	i := 0
	for _, g := range groupStorage {
		new := *g
		log.With(logger.TraceFieldFromAny(new)).Debug("new")
		groups[i] = &new
		i++
	}
	return groups, nil
}

func (g *GroupMemory) AddStudent(ctx context.Context, groupID domain.GroupID, userIDs []domain.UserID) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(groupID),
		logger.TraceFieldFromAny(userIDs),
	)
	log.Debug("Called a addStudent repository method")

	group, ok := groupStorage[groupID]
	if !ok {
		return fmt.Errorf("%w: group_id=%d", ErrNotFound, groupID)
	}

	if len(userIDs) == 0 {
		return fmt.Errorf("%w: length a userIDs is 0", ErrInvalid)
	}

	uniqueIDs := make(map[domain.UserID]bool)
	for _, student := range group.Students {
		uniqueIDs[student.ID] = true
	}

	for _, id := range userIDs {
		if _, ok := uniqueIDs[id]; ok {
			return fmt.Errorf("%w: user_id=%d", ErrCollision, id)
		}
		if _, ok := userStorage[id]; !ok {
			return fmt.Errorf("%w: user_id=%d", ErrNotFound, id)
		}
		uniqueIDs[id] = true
	}

	newStudents := make([]*domain.User, len(uniqueIDs))
	i := 0
	for id := range uniqueIDs {
		newStudents[i] = userStorage[id]
		i++
	}

	group.Students = newStudents
	return nil
}

func (g *GroupMemory) RemoveStudent(ctx context.Context, groupID domain.GroupID, userIDs []domain.UserID) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(groupID),
		logger.TraceFieldFromAny(userIDs),
	)
	log.Debug("Called a removeStudent repository method")

	group, ok := groupStorage[groupID]
	if !ok {
		return fmt.Errorf("%w: group_id=%d", ErrNotFound, groupID)
	}

	setUserIDs := make(map[domain.UserID]bool)
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
		newStudents[i] = userStorage[id]
	}

	group.Students = newStudents
	return nil
}
