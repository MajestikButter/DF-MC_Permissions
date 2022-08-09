package permissionscmd

import (
	"github.com/MajestikButter/DF-MC_Permissions/permissions"
	"github.com/df-mc/dragonfly/server/cmd"
)

type GroupPermView struct {
	Group       cmd.SubCommand `cmd:"group"`
	Identifier  string         `cmd:"group"`
	Permissions cmd.SubCommand `cmd:"permissions"`
	View        cmd.SubCommand `cmd:"view"`
}

func (t GroupPermView) Run(source cmd.Source, output *cmd.Output) {
	system, err := permissions.GetSystem()
	if err != nil {
		output.Errorf("Failed to get permission system")
	}

	if !system.HasGroup(t.Identifier) {
		output.Errorf("Group with identifier %s doesn't exist", t.Identifier)
		return
	}
	g := system.GetGroup(t.Identifier)
	output.Printf("Group %s permissions\n%s", t.Identifier, g.Perms.String())
}

func (t GroupPermView) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.group.permissions.view")
}
