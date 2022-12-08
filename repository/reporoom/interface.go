package reporoom

import "context"

type RoomFinder interface {
	IsAvaiable(ctx context.Context, room string) (bool, error)
}
