package main

import "errors"

// command represents a CLI command with its name and arguments.
type command struct {
	Name string
	Args []string
}

// commands holds a registry of all available CLI commands and their handlers.
type commands struct {
	registeredCommands map[string]func(*state, command) error
}

// register adds a new command handler to the registry.
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

// run executes a command by looking up its handler in the registry
// and calling it with the provided state and command context.
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
