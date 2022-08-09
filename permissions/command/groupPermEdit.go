package command

import (
	"github.com/MajestikButter/DF-MC_Permissions/permissions"
	"github.com/df-mc/dragonfly/server/cmd"
)

type GroupPermAdd struct {
	Group       cmd.SubCommand `cmd:"group"`
	Identifier  string         `cmd:"group"`
	Permissions cmd.SubCommand `cmd:"permissions"`
	Add         cmd.SubCommand `cmd:"add"`
	Permission  string         `cmd:"permission"`
}

func (t GroupPermAdd) Run(source cmd.Source, output *cmd.Output) {
	system, err := permissions.GetSystem()
	if err != nil {
		output.Errorf("Failed to get permission system")
	}

	if !system.HasGroup(t.Identifier) {
		output.Errorf("Group with identifier %s doesn't exist", t.Identifier)
		return
	}
	g := system.GetGroup(t.Identifier)

	erra := g.Perms.AddPermission(t.Permission)
	if erra != nil {
		output.Errorf("An error occured while adding permission to group %s:\n%s", t.Identifier, err.Error())
		return
	}

	output.Printf("Added %s permission to group %s", t.Permission, t.Identifier)
}

func (t GroupPermAdd) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.group.permissions.add")
}

type GroupPermRemove struct {
	Group       cmd.SubCommand `cmd:"group"`
	Identifier  string         `cmd:"group"`
	Permissions cmd.SubCommand `cmd:"permissions"`
	Remove      cmd.SubCommand `cmd:"remove"`
	Permission  string         `cmd:"permission"`
}

func (t GroupPermRemove) Run(source cmd.Source, output *cmd.Output) {
	system, err := permissions.GetSystem()
	if err != nil {
		output.Errorf("Failed to get permission system")
	}

	if !system.HasGroup(t.Identifier) {
		output.Errorf("Group with identifier %s doesn't exist", t.Identifier)
		return
	}
	g := system.GetGroup(t.Identifier)

	erra := g.Perms.RemovePermission(t.Permission)
	if erra != nil {
		output.Errorf("An error occured while adding permission to group %s:\n%s", t.Identifier, err.Error())
		return
	}

	output.Printf("Removed %s permission from group %s", t.Permission, t.Identifier)
}

func (t GroupPermRemove) Allow(source cmd.Source) bool {
	system, err := permissions.GetSystem()
	if err != nil {
		return false
	}

	return system.GetCommandPermission(source, "minecraft.chat.command.permissions.group.permissions.remove")
}
