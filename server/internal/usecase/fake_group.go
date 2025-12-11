package usecase

import (
	"context"
	"fmt"
	"time"
)

type FakeGroup struct{}

func NewFakeGroup() *FakeGroup {
	return &FakeGroup{}
}

func (f *FakeGroup) Create(ctx context.Context, name string) error {
	return nil
}

func (f *FakeGroup) Get(ctx context.Context) ([]GroupDomain, error) {
	var groups []GroupDomain
	groups = initFakeGroupSlices(groups, "А")
	groups = initFakeGroupSlices(groups, "Б")
	groups = initFakeGroupSlices(groups, "В")
	return groups, nil
}

func initFakeGroupSlices(groups []GroupDomain, word string) []GroupDomain {
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
