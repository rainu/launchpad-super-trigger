package actor

import "context"

type Actor interface {
	Do(ctx context.Context) error
}
