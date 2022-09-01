package commandProcessor

import "context"

type ICommandProcessor interface {
	ProcessCommand(ctx context.Context, command, params string) string
}
