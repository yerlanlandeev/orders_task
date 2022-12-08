package reporder

import (
	"context"
	"errors"
	"sync"
	"time"
)

// map[room][]Order
// map[string][]Order
type RoomOrder struct {
	mutex  *sync.Mutex
	Orders []Order
}

func NewRoomOrder() *RoomOrder {
	return &RoomOrder{
		mutex: &sync.Mutex{},
	}
}

var _ OrderMaker = (*RoomOrder)(nil)

type ParamsOrder struct {
	Room      string    `json:"room"`
	UserEmail string    `json:"user_email"`
	From      time.Time `json:"from"`
	To        time.Time `json:"to"`
}

type ResponseOrder struct{}

func (r *RoomOrder) Order(ctx context.Context, args Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.Orders = append(r.Orders, args)
	return nil
}

type ParamsGetAll struct {
	Offset int
	Limit  int
}

func (r *RoomOrder) GetAll(ctx context.Context, args ParamsGetAll) ([]Order, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if args.Offset > len(r.Orders) {
		return nil, errors.New("no room from this offset")
	}
	if args.Offset+args.Limit > cap(r.Orders) {
		args.Limit = cap(r.Orders) - args.Offset
	}
	return r.Orders[args.Offset : args.Offset+args.Limit], nil
}

func (r *RoomOrder) Delete(ctx context.Context, room string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// O(n) time com
	// space O(n)
	var orders []Order
	for _, order := range r.Orders {
		if order.Room == room {
			continue
		}
		orders = append(orders, order)
	}

	r.Orders = orders

	return nil
}
