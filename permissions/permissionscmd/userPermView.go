package permissionscmd

import (
	"github.com/MajestikButter/DF-MC_Permissions/permissions"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type UserPermView struct {
	User        cmd.SubCommand `cmd:"user"`
	Target      []cmd.Target   `cmd:"user"`
	Permissions cmd.SubCommand `cmd:"permissions"`
	View        cmd.SubCommand `cmd:"view"`
}

func (t UserPermView) Run(source cmd.Source, output *cmd.Output) {
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
		output.Printf("User %s permissions\n%s", plr.Name(), user.Perms.String())
	}
}

func (t UserPermView) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.user.permissions.view")
}
