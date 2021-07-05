package command

import (
	"fmt"

	"github.com/slack-go/slack/slackevents"
)

type Command interface {
	Run(msg string, event *slackevents.AppMentionEvent) (string, error)
}

func RunCommand(name string, event *slackevents.AppMentionEvent) (string, error) {
	registry := NewRegistry()
	registry.Register("balance", &BalanceCommand{})
	cmd, err := registry.Lookup(name)

	if err != nil {
		return "", err
	}

	return cmd.Run(name, event)
}

// CommandRegistry contains Command implementations that implement custom
// behaviors for each supported koin command
type Registry struct {
	impls map[string]Command
}

func NewRegistry() *Registry {
	return &Registry{impls: map[string]Command{}}
}

// Register registers impl for commands. It will be called by ProcessMessage()
func (r *Registry) Register(name string, impl Command) {
	r.impls[name] = impl
}

// Lookup returns the Command implementation or an error if one can't be found.
func (r *Registry) Lookup(name string) (Command, error) {
	impl, ok := r.impls[name]
	if !ok {
		var cmd Command
		return cmd, fmt.Errorf("unknown command: %s", name)
	}
	return impl, nil
}
