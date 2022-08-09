# **DF-MC_Permissions**

## Installation

```
go get github.com/MajestikButter/DF-MC_Permissions@latest
```

## Usage

### Setting up

```go
import (
  "github.com/MajestikButter/DF-MC_Permissions/permissions"
  "github.com/MajestikButter/DF-MC_Permissions/permissions/permissionscmd"
)

func main() {
  // Loads the permission system to be used in commands and elsewhere
  permissions.LoadSystem("permissions/users.json", "permissions/groups.json", false)

  // Loads the permissions command to be used ingame. Use permissionscmd.LoadCustom() to customize the command name, description, and aliases
  permissionscmd.LoadDefault()
}
```

### Checking permissions

```go
func CanEat(player *player.Player) bool {
  system := permissions.getSystem()
  user := system.GetUser(player)
  user.HasPermission("example.eat")
}
```
