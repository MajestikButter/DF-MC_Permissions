package permissions

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"golang.org/x/exp/slices"
)

type PermissionSystem struct {
	users  map[string]*User
	groups map[string]*group

	userFile  string
	groupFile string

	defaultPermission bool
}

var system *PermissionSystem

// Create a new permission system struct, returns a pointer to the new PermissionSystem
// and an error if the system has already been loaded.
func LoadSystem(userFile string, groupFile string, defaultPermission bool) (*PermissionSystem, error) {
	if system != nil {
		return system, errors.New("permission system has already been loaded")
	}

	s := &PermissionSystem{map[string]*User{}, map[string]*group{}, userFile, groupFile, defaultPermission}
	s.Load()

	if !s.HasGroup("default") {
		s.AddGroup("default", s.CreateGroup(s.CreatePermissions([]string{})))
	}

	s.Save()

	go func() {
		saveTicker := time.NewTicker(time.Minute * 5)
		for {
			<-saveTicker.C
			s.Save()
		}
	}()

	system = s
	return s, nil
}

// Gets the loaded permission system, returns an error if no system is loaded
func GetSystem() (*PermissionSystem, error) {
	if system == nil {
		return nil, errors.New("no permission system is currently loaded")
	}
	return system, nil
}

// Create a new permissions struct
func (s *PermissionSystem) CreatePermissions(permissions []string) *Permissions {
	return newPermissions(permissions)
}

func (s *PermissionSystem) GetCommandPermission(source cmd.Source, perm string) bool {
	plr, isPlr := source.(*player.Player)
	if !isPlr {
		return true
	}
	return s.GetUser(plr).HasPermission(perm)
}

// Get permissions user by server player struct
func (s *PermissionSystem) GetUser(player *player.Player) *User {
	x := player.XUID()
	if s.users[x] != nil {
		return s.users[x]
	}
	user := newUser(s.CreatePermissions([]string{}), []string{"default"}, s)
	s.users[x] = user
	return user
}

// Create a new group (won't be automatically added)
func (s *PermissionSystem) CreateGroup(permissions *Permissions) *group {
	return newGroup(permissions, s)
}

// Add a group to the system
func (s *PermissionSystem) AddGroup(id string, g *group) {
	s.groups[id] = g
}

// Check if a group is in the system
func (s *PermissionSystem) HasGroup(id string) bool {
	return s.groups[id] != nil
}

// Delete an existing group using it's identifier
func (s *PermissionSystem) DelGroup(id string) {
	for _, v := range s.users {
		index := slices.IndexFunc(v.Groups, func(group string) bool {
			return group == id
		})
		v.Groups = shiftIndex(v.Groups, index)
	}
	delete(s.groups, id)
}

// Get an existing group by it's identifier
func (s *PermissionSystem) GetGroup(id string) *group {
	return s.groups[id]
}

// Load stuff

func (s *PermissionSystem) Load() {
	s.loadUsers()
	s.loadGroups()
}

// Load users save file
func (s *PermissionSystem) loadUsers() {
	contents, err := loadUserFile(s.userFile, s)
	if contents == nil || err != nil {
		return
	}
	s.users = contents
}

// Load groups save file
func (s *PermissionSystem) loadGroups() {
	contents, err := loadGroupFile(s.groupFile, s)
	if contents == nil || err != nil {
		return
	}
	s.groups = contents
}

// Save stuff

func (s *PermissionSystem) Save() (error, error) {
	return s.saveUsers(), s.saveGroups()
}

// Write to users file
func (s *PermissionSystem) saveUsers() error {
	contents, err := json.MarshalIndent(s.users, "", "  ")
	if err != nil {
		return err
	}
	return s.writeFile(s.userFile, contents)
}

// Write to groups file
func (s *PermissionSystem) saveGroups() error {
	contents, err := json.MarshalIndent(s.groups, "", "  ")
	if err != nil {
		return err
	}
	return s.writeFile(s.groupFile, contents)
}

func (s *PermissionSystem) writeFile(path string, data []byte) error {
	path = strings.ReplaceAll(path, "\\", "/")
	split := strings.Split(path, "/")
	if len(split) > 1 {
		parent := split[:len(split)-1]
		os.Mkdir(strings.Join(parent, "/"), 0666)
	}
	return os.WriteFile(path, data, 0666)
}
