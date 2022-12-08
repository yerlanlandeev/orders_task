package reporder

import "context"

type OrderMaker interface {
	Order(ctx context.Context, args Order) error
	GetAll(ctx context.Context, args ParamsGetAll) ([]Order, error)
	Delete(ctx context.Context, room string) error
}
