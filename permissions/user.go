package permissions

import (
	"encoding/json"
	"errors"
	"os"

	"golang.org/x/exp/slices"
)

type User struct {
	Perms  *Permissions
	Groups []string
	system *PermissionSystem
}

func newUser(perms *Permissions, groups []string, system *PermissionSystem) *User {
	return &User{perms, groups, system}
}

func (p *User) GetGroups() map[string]*group {
	res := map[string]*group{}
	for _, v := range p.Groups {
		if p.system.HasGroup(v) {
			res[v] = p.system.GetGroup(v)
		}
	}
	return res
}

func (p *User) HasGroup(identifier string) bool {
	return slices.Contains(p.Groups, identifier)
}

func (p *User) AddGroup(identifier string) error {
	if p.HasGroup(identifier) {
		return errors.New("user already has this group")
	}
	p.Groups = append(p.Groups, identifier)
	return nil
}

func (p *User) RemoveGroup(identifier string) error {
	if !p.HasGroup(identifier) {
		return errors.New("user doesn't have this group")
	}

	i := slices.IndexFunc(p.Groups, func(v string) bool {
		return v == identifier
	})

	if i < 0 {
		return errors.New("failed to get index of group")
	}

	p.Groups = shiftIndex(p.Groups, i)
	return nil
}

func (p *User) HasPermission(permission string) bool {
	allowed, depth := p.Perms.Get(permission, p.system.defaultPermission)

	for _, g := range p.GetGroups() {
		a, d := g.Perms.Get(permission, g.system.defaultPermission)
		if d > depth {
			depth = d
			allowed = a
		}
	}

	return allowed
}

func loadUserFile(path string, system *PermissionSystem) (map[string]*User, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("error when opening file")
	}

	users := map[string]*User{}
	err = json.Unmarshal(content, &users)
	if err != nil {
		return nil, errors.New("error during Unmarshal()")
	}

	for _, p := range users {
		p.system = system
	}

	return users, nil
}
