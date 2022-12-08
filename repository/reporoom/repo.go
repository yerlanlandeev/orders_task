package reporoom

import "context"

type Room struct {
	availableRooms map[string]struct{}
}

func NewRoom() *Room {
	return &Room{
		availableRooms: map[string]struct{}{"econom": {}, "standart": {}, "lux": {}},
	}
}

var _ RoomFinder = (*Room)(nil)

func (r Room) IsAvaiable(ctx context.Context, room string) (bool, error) {
	_, ok := r.availableRooms[room]
	return ok, nil
}
