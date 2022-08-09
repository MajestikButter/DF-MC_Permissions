package permissionscmd

import (
	"github.com/MajestikButter/DF-MC_Permissions/permissions"
	"github.com/df-mc/dragonfly/server/cmd"
)

type GroupCreate struct {
	Group      cmd.SubCommand `cmd:"group"`
	Create     cmd.SubCommand `cmd:"create"`
	Identifier string         `cmd:"identifier"`
}

func (t GroupCreate) Run(source cmd.Source, output *cmd.Output) {
	system, err := permissions.GetSystem()
	if err != nil {
		output.Errorf("Failed to get permission system")
	}

	if system.HasGroup(t.Identifier) {
		output.Errorf("Group with identifier %s already exists", t.Identifier)
		return
	}
	g := system.CreateGroup(system.CreatePermissions([]string{}))
	system.AddGroup(t.Identifier, g)
	output.Printf("Created group with identifier %s", t.Identifier)
}

func (t GroupCreate) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.group.create")
}

type GroupDelete struct {
	Group      cmd.SubCommand `cmd:"group"`
	Delete     cmd.SubCommand `cmd:"delete"`
	Identifier string         `cmd:"identifier"`
}

func (t GroupDelete) Run(source cmd.Source, output *cmd.Output) {
	system, err := permissions.GetSystem()
	if err != nil {
		output.Errorf("Failed to get permission system")
	}

	if !system.HasGroup(t.Identifier) {
		output.Errorf("Group with identifier %s doesn't exist", t.Identifier)
		return
	}
	system.DelGroup(t.Identifier)
	output.Printf("Deleted group with identifier %s", t.Identifier)
}

func (t GroupDelete) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.group.delete")
}
