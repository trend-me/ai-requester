package interfaces

import "context"

type Ai interface {
	Prompt(ctx context.Context, prompt string) (string, error)
}
