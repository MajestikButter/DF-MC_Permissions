package permissions

import (
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/exp/slices"
)

type Permissions struct {
	allow      [][]string
	disallow   [][]string
	superAllow [][]string
}

type permissionsJSON struct {
	Allow      []string `json:"allow"`
	Disallow   []string `json:"disallow"`
	SuperAllow []string `json:"superAllow"`
}

func ValidPermissionString(perm string) bool {
	return strings.HasPrefix(perm, "+") || strings.HasPrefix(perm, "-") || strings.HasPrefix(perm, "*")
}

func newPermissions(perms []string) *Permissions {
	p := &Permissions{[][]string{}, [][]string{}, [][]string{}}
	for _, v := range perms {
		if err := p.AddPermission(v); err != nil {
			panic(err)
		}
	}
	return p
}

func (p *Permissions) compare(check []string, slice []string) (bool, int) {
	cLen := len(check)
	if len(slice) > cLen {
		return false, -1
	}

	slen := len(slice)
	for i, v := range check {
		if v != slice[i] {
			return false, -1
		}
		if slen-1 <= i {
			return true, i
		}
	}
	return false, -1
}

func (p *Permissions) checkPermArray(checker []string, compareTo [][]string, currDepth int) (bool, int) {
	found := false
	for _, v := range compareTo {
		if len(v) > len(checker) {
			continue
		}
		res, depth := p.compare(checker, v)
		if res && (depth >= currDepth) {
			found = true
			currDepth = depth
		}
	}
	return found, currDepth
}

// Returns whether the permission is allowed and the depth it was found at
func (p *Permissions) Get(perm string, defaultPermission bool) (bool, int) {
	split := strings.Split(perm, ".")

	allowed := defaultPermission
	depth := -1
	found := false

	found, depth = p.checkPermArray(split, p.allow, depth)
	if found {
		allowed = true
	}

	found, depth = p.checkPermArray(split, p.disallow, depth)
	if found {
		allowed = false
	}

	found, depth = p.checkPermArray(split, p.superAllow, depth)
	if found {
		allowed = true
	}

	return allowed, depth
}

func (p *Permissions) permExists(split []string) bool {
	_, aD := p.checkPermArray(split, p.allow, -1)
	if aD > -1 {
		return true
	}
	_, dD := p.checkPermArray(split, p.disallow, -1)
	if dD > -1 {
		return true
	}
	_, sD := p.checkPermArray(split, p.superAllow, -1)
	return sD > -1
}

func (p *Permissions) AddPermission(perm string) error {
	split := strings.Split(perm[1:], ".")
	if p.permExists(split) {
		return errors.New("permission already exists")
	}

	if strings.HasPrefix(perm, "-") {
		p.disallow = append(p.disallow, split)
	} else if strings.HasPrefix(perm, "+") {
		p.allow = append(p.allow, split)
	} else if strings.HasPrefix(perm, "*") {
		p.superAllow = append(p.superAllow, split)
	} else {
		return errors.New("permission string invalid, must be prefixed with +, -, or *")
	}
	return nil
}

func findIndex(slice [][]string, perm []string) int {
	return slices.IndexFunc(slice, func(s []string) bool {
		if len(s) != len(perm) {
			return false
		}

		for i, v := range perm {
			if s[i] != v {
				return false
			}
		}
		return true
	})
}

func (p *Permissions) RemovePermission(perm string) error {
	split := strings.Split(perm, ".")
	if !p.permExists(split) {
		return errors.New("permission doesn't exist")
	}

	aI := findIndex(p.allow, split)
	if aI >= 0 {
		p.allow = shiftIndex(p.allow, aI)
	}

	dI := findIndex(p.disallow, split)
	if dI >= 0 {
		p.disallow = shiftIndex(p.disallow, dI)
	}

	sI := findIndex(p.superAllow, split)
	if sI >= 0 {
		p.superAllow = shiftIndex(p.superAllow, sI)
	}

	return nil
}

// Returns a string of the permissions
func (p *Permissions) String() string {
	res := ""
	for _, v := range p.allow {
		s := ", +" + strings.Join(v, ".")
		res += s
	}
	for _, v := range p.disallow {
		s := ", -" + strings.Join(v, ".")
		res += s
	}
	for _, v := range p.superAllow {
		s := ", *" + strings.Join(v, ".")
		res += s
	}
	return res[2:]
}

func (p *Permissions) MarshalJSON() ([]byte, error) {
	allow := []string{}
	for _, v := range p.allow {
		allow = append(allow, strings.Join(v, "."))
	}

	disallow := []string{}
	for _, v := range p.disallow {
		disallow = append(disallow, strings.Join(v, "."))
	}

	superAllow := []string{}
	for _, v := range p.superAllow {
		superAllow = append(superAllow, strings.Join(v, "."))
	}

	return json.Marshal(&permissionsJSON{allow, disallow, superAllow})
}

func (p *Permissions) UnmarshalJSON(data []byte) error {
	jsonStruct := permissionsJSON{}
	err := json.Unmarshal(data, &jsonStruct)
	if err != nil {
		return err
	}

	p.allow = [][]string{}
	for _, v := range jsonStruct.Allow {
		p.allow = append(p.allow, strings.Split(v, "."))
	}

	p.disallow = [][]string{}
	for _, v := range jsonStruct.Disallow {
		p.disallow = append(p.disallow, strings.Split(v, "."))
	}

	p.superAllow = [][]string{}
	for _, v := range jsonStruct.SuperAllow {
		p.superAllow = append(p.superAllow, strings.Split(v, "."))
	}

	return nil
}
