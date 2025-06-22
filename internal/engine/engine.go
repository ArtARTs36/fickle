package engine

import "context"

type Engine interface {
	Find(ctx context.Context, serviceName string) (*Container, error)
	Stop(ctx context.Context, id string) error
	Start(ctx context.Context, id string) error
}
