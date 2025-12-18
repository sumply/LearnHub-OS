package usecase

import (
	"context"
	"fmt"
	"time"
)

type StubGroup struct{}

func NewStubGroup() *StubGroup {
	return &StubGroup{}
}

func (f *StubGroup) Create(ctx context.Context, name string) error {
	return nil
}

func (f *StubGroup) Get(ctx context.Context) ([]GroupDomain, error) {
	var groups []GroupDomain
	groups = initStubGroupSlices(groups, "А")
	groups = initStubGroupSlices(groups, "Б")
	groups = initStubGroupSlices(groups, "В")
	return groups, nil
}

func initStubGroupSlices(groups []GroupDomain, word string) []GroupDomain {
	for i := 1; i < 11; i++ {
		g := GroupDomain{
			ID:        ID(i),
			Name:      fmt.Sprintf("%d%s", i, word),
			CreatedAt: time.Now(),
		}
		groups = append(groups, g)
	}
	return groups
}
