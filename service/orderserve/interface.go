package orderserve

import "context"

type OrderProcessor interface {
	Order(ctx context.Context, args Order) (ResponseOrder, error)
	GetAll(ctx context.Context, args ParamsGetAll) ([]Order, error)
}
