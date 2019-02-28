package acl

import (
	"fmt"
)

func buildRole(role string) *Role {
	if r, err := RegisterRole(role); err != nil {
		panic(err)
	} else {
		return r
	}
}

func (ro *Role) Same(role *Role) bool {
	if ro.Instance > 0 {
		return ro.Role == role.Role &&
			ro.Instance == role.Instance
	} else {
		return ro.Role == role.Role
	}
}

func (ro *Role) String() string {
	if ro.Instance > 0 {
		return fmt.Sprintf("%s->%d", ro.Role, ro.Instance)
	} else {
		return fmt.Sprintf("%s", ro.Role)
	}
}
