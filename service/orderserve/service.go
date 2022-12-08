package orderserve

import (
	"applicationDesignTest/repository/reporder"
	"applicationDesignTest/repository/reporoom"
	"context"
	"encoding/json"
	"fmt"
)

type OrderProcess struct {
	OrderRepo reporder.OrderMaker
	RoomRepo  reporoom.RoomFinder
}

func NewOrder(orderRepo reporder.OrderMaker, roomRepo reporoom.RoomFinder) *OrderProcess {
	return &OrderProcess{
		OrderRepo: orderRepo,
		RoomRepo:  roomRepo,
	}
}

var _ OrderProcessor = (*OrderProcess)(nil)

type ResponseOrder struct {
	Successfull bool
	Msg         string
}

func (order *OrderProcess) Order(ctx context.Context, args Order) (ResponseOrder, error) {
	// check whether room is available
	ok, err := order.RoomRepo.IsAvaiable(ctx, args.Room)
	switch {
	case err != nil:
		return ResponseOrder{}, fmt.Errorf("failed to check availability. err: %v", err)
	case !ok:
		return ResponseOrder{}, fmt.Errorf("such room is not available. err: %v", err)
	}

	// order the room
	err = order.OrderRepo.Order(ctx, reporder.Order{
		Room:      args.Room,
		UserEmail: args.UserEmail,
		From:      args.From,
		To:        args.To,
	})
	if err != nil {
		return ResponseOrder{}, fmt.Errorf("failed to order the room. err: %v", err)
	}

	return ResponseOrder{
		Successfull: true,
		Msg:         fmt.Sprintf("successfully booked room, type of %s", args.Room),
	}, err
}

type ParamsGetAll struct {
	Offset int
	Limit  int
}

func (order *OrderProcess) GetAll(ctx context.Context, args ParamsGetAll) ([]Order, error) {
	orders, err := order.OrderRepo.GetAll(ctx, reporder.ParamsGetAll{
		Offset: args.Offset,
		Limit:  args.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all orders")
	}

	b, err := json.Marshal(orders)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all orders")
	}

	var ordersResponse = make([]Order, 0, len(orders))
	if err := json.Unmarshal(b, &ordersResponse); err != nil {
		return nil, fmt.Errorf("failed to fetch all orders")
	}

	return ordersResponse, nil
}
