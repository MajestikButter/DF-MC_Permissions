package command

import (
	"github.com/MajestikButter/DF-MC_Permissions/permissions"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type UserPermAdd struct {
	User        cmd.SubCommand `cmd:"user"`
	Target      []cmd.Target   `cmd:"user"`
	Permissions cmd.SubCommand `cmd:"permissions"`
	Add         cmd.SubCommand `cmd:"add"`
	Permission  string         `cmd:"permission"`
}

func (t UserPermAdd) Run(source cmd.Source, output *cmd.Output) {
	system, err := permissions.GetSystem()
	if err != nil {
		output.Errorf("Failed to get permission system")
	}

	players := []*player.Player{}
	users := []*permissions.User{}

	for _, target := range t.Target {
		plr, ok := target.(*player.Player)
		if !ok {
			output.Error("Invalid selector")
			return
		}
		players = append(players, plr)
		users = append(users, system.GetUser(plr))
	}

	for i, user := range users {
		plr := players[i]

		err := user.Perms.AddPermission(t.Permission)
		if err != nil {
			output.Errorf("An error occured while adding permission to user %s:\n%s", plr.Name(), err.Error())
			return
		}
		output.Printf("Added %s permission to user %s", t.Permission, plr.Name())
	}
}

func (t UserPermAdd) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.user.permissions.add")
}

type UserPermRemove struct {
	User        cmd.SubCommand `cmd:"user"`
	Target      []cmd.Target   `cmd:"user"`
	Permissions cmd.SubCommand `cmd:"permissions"`
	Remove      cmd.SubCommand `cmd:"remove"`
	Permission  string         `cmd:"permission"`
}

func (t UserPermRemove) Run(source cmd.Source, output *cmd.Output) {
	system, err := permissions.GetSystem()
	if err != nil {
		output.Errorf("Failed to get permission system")
	}

	players := []*player.Player{}
	users := []*permissions.User{}

	for _, target := range t.Target {
		plr, ok := target.(*player.Player)
		if !ok {
			output.Error("Invalid selector")
			return
		}
		players = append(players, plr)
		users = append(users, system.GetUser(plr))
	}

	for i, user := range users {
		plr := players[i]

		err := user.Perms.RemovePermission(t.Permission)
		if err != nil {
			output.Errorf("An error occured while removing permission from user %s:\n%s", plr.Name(), err.Error())
			return
		}
		output.Printf("Added %s permission to user %s", t.Permission, plr.Name())
	}
}

func (t UserPermRemove) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.user.permissions.remove")
}
