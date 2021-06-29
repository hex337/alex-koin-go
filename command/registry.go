package command

import (
	"fmt"
)

type Command interface {
	Run(msg string) (string, error)
}

func RunCommand(name string) (string, error) {
	registry := NewRegistry()
	registry.Register("balance", &BalanceCommand{})
	cmd, _ := registry.Lookup(name)
	//TODO do something with err
	return cmd.Run(name)
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
