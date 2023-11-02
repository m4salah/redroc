package pubsub

import "context"

type PubSub interface {
	Pub(ctx context.Context, message any)
	Sub(ctx context.Context, message any)
}
