package prompt

import "context"

func Always(ctx context.Context) bool {
	return true
}

func ContextValueIs(ctx context.Context, key any, value any) bool {
	return ctx.Value(key) == value
}
