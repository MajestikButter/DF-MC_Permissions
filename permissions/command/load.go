package command

import (
	"errors"

	"github.com/df-mc/dragonfly/server/cmd"
)

var loaded = false

// Loads the permission command using a custom name, description, and aliases
// Returns error if the command has already been loaded
func LoadCustom(name string, description string, aliases []string) error {
	if loaded {
		return errors.New("command has already been loaded")
	}
	cmd.Register(cmd.New(name, description, aliases, GroupCreate{}, GroupDelete{}, GroupPermView{}, GroupPermAdd{}, GroupPermRemove{}, GroupUserAdd{}, GroupUserRemove{}, UserPermAdd{}, UserPermRemove{}, UserPermView{}))
	loaded = true
	return nil
}

// Loads the permission command using the default name, description, and aliases
// Returns error if the command has already been loaded
func LoadDefault() error {
	return LoadCustom("permissions", "Handle user and group permissions", []string{"permissions", "perms"})
}
