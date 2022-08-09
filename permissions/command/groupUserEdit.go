package command

import (
	"github.com/MajestikButter/DF-MC_Permissions/permissions"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type GroupUserAdd struct {
	Group      cmd.SubCommand `cmd:"group"`
	Identifier string         `cmd:"group"`
	User       cmd.SubCommand `cmd:"user"`
	Add        cmd.SubCommand `cmd:"add"`
	Target     []cmd.Target   `cmd:"user"`
}

func (t GroupUserAdd) Run(source cmd.Source, output *cmd.Output) {
	system, err := permissions.GetSystem()
	if err != nil {
		output.Errorf("Failed to get permission system")
	}

	if !system.HasGroup(t.Identifier) {
		output.Errorf("Group with identifier %s doesn't exist", t.Identifier)
		return
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
		if user.HasGroup(t.Identifier) {
			output.Errorf("%s already has group %s\n", plr.Name(), t.Identifier)
			continue
		}
		err := user.AddGroup(t.Identifier)
		if err != nil {
			output.Errorf("An error occured while adding %s to group %s:\n%s", plr.Name(), t.Identifier, err.Error())
			return
		}
		output.Printf("Added %s to group %s", plr.Name(), t.Identifier)
	}
}

func (t GroupUserAdd) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.group.user.add")
}

type GroupUserRemove struct {
	Group      cmd.SubCommand `cmd:"group"`
	Identifier string         `cmd:"group"`
	User       cmd.SubCommand `cmd:"user"`
	Remove     cmd.SubCommand `cmd:"remove"`
	Target     []cmd.Target   `cmd:"user"`
}

func (t GroupUserRemove) Run(source cmd.Source, output *cmd.Output) {
	system, err := permissions.GetSystem()
	if err != nil {
		output.Errorf("Failed to get permission system")
	}

	if !system.HasGroup(t.Identifier) {
		output.Errorf("Group with identifier %s doesn't exist", t.Identifier)
		return
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
		if !user.HasGroup(t.Identifier) {
			output.Errorf("%s isn't in group %s\n", plr.Name(), t.Identifier)
			continue
		}
		err := user.RemoveGroup(t.Identifier)
		if err != nil {
			output.Errorf("An error occured while removing %s from group %s:\n%s", plr.Name(), t.Identifier, err.Error())
			return
		}
		output.Printf("Removed %s from group %s", plr.Name(), t.Identifier)
	}
}

func (t GroupUserRemove) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.group.user.remove")
}
