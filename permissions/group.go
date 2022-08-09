package permissions

import (
	"encoding/json"
	"errors"
	"os"
)

type group struct {
	Perms  *Permissions
	system *PermissionSystem
}

func newGroup(perms *Permissions, system *PermissionSystem) *group {
	return &group{perms, system}
}

func loadGroupFile(path string, system *PermissionSystem) (map[string]*group, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("error when opening file")
	}

	groups := map[string]*group{}
	err = json.Unmarshal(content, &groups)
	if err != nil {
		return nil, errors.New("error during Unmarshal()")
	}

	for _, g := range groups {
		g.system = system
	}

	return groups, nil
}
