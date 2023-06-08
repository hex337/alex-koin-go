package command

import (
	"fmt"
)

type Command interface {
	Run(msg string, event *CoinEvent) (BotResponse, error)
}

func RunCommand(name string, event *CoinEvent) (BotResponse, error) {
	registry := NewRegistry()
	registry.Register("all_nfts", &AllNftsCommand{})
	registry.Register("balance", &BalanceCommand{})
	registry.Register("create_coin", &CreateCoinCommand{})
	registry.Register("create_nft", &CreateNftCommand{})
	registry.Register("destroy_coin", &DestroyCoinCommand{})
	registry.Register("my_nfts", &MyNftsCommand{})
	registry.Register("stats", &StatsCommand{})
	registry.Register("transfer_coin", &TransferCoinCommand{})
	registry.Register("transfer_nft", &TransferNftCommand{})
	registry.Register("what_am_i", &IdentityCommand{})
	registry.Register("what_are_you", &SelfIdentityCommand{})
	cmd, err := registry.Lookup(name)

	if err != nil {
		return BotResponse{Text: ""}, err
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
