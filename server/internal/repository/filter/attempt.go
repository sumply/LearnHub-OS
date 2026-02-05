package filter

import "github.com/google/uuid"

type Attempt struct {
	UserID      *uuid.UUID
	IsCompleted *bool
}
