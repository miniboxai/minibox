package acl

import (
	"strconv"
	"strings"
)

var (
	roleExtensions   []string
	actionExtensions []string
	groupExtensions  []string
)

func RegisterRole(roleStr string) (role *Role, err error) {
	var (
		ss       = strings.SplitN(roleStr, "-", 2)
		roleName string
		instance int
	)

	if len(ss) == 1 { //without Instance
		roleName = ss[0]
	} else if len(ss) == 2 {
		roleName = ss[0]
		if instance, err = strconv.Atoi(ss[1]); err != nil {
			return nil, err
		}
	} else {
		return nil, ErrInvalidRoleString
	}

	l := len(_RoleType_index) - 1

	for i := range _RoleType_index[:l] {
		name := _RoleType_name[_RoleType_index[i]:_RoleType_index[i+1]]
		if roleName == name {
			role = new(Role)
			role.Role = RoleType(i)
			role.Instance = instance
			return
		}
	}

	var (
		j   int
		ext string
	)

	for j, ext = range roleExtensions {
		if roleName == ext && instance == 0 {
			return nil, ErrAlwaysRegisterRole
		} else if roleName == ext {
			role = new(Role)
			role.Role = RoleType(l + j)
			role.Instance = instance
			return
		}
	}

	roleExtensions = append(roleExtensions, roleName)
	role = new(Role)
	role.Role = RoleType(l + j)
	role.Instance = instance
	return role, nil
}

func RegisterAction(actStr string) (Action, error) {
	var (
		actName string
	)

	actName = actStr

	l := len(_Action_index) - 1

	for i := range _Action_index[:l] {
		name := _Action_name[_Action_index[i]:_Action_index[i+1]]
		if actName == name {
			return Action(i), nil
		}
	}

	var (
		j   int
		ext string
	)

	for j, ext = range actionExtensions {
		if actName == ext {
			return None, ErrAlwaysRegisterRole
		}
	}

	actionExtensions = append(actionExtensions, actName)
	return (Action)(l + j), nil
}

func RegisterGRole(roleStr string) (role *GRole, err error) {
	var (
		ss       = strings.SplitN(roleStr, "-", 2)
		roleName string
		instance int
	)

	if len(ss) == 1 { //without Instance
		roleName = ss[0]
	} else if len(ss) == 2 {
		roleName = ss[0]
		if instance, err = strconv.Atoi(ss[1]); err != nil {
			return nil, err
		}
	} else {
		return nil, ErrInvalidRoleString
	}

	l := len(_GRoleType_index) - 1

	for i := range _GRoleType_index[:l] {
		name := _GRoleType_name[_GRoleType_index[i]:_GRoleType_index[i+1]]
		if roleName == name {
			role = new(GRole)
			role.Role = GRoleType(i)
			role.Instance = instance
			return
		}
	}

	var (
		j   int
		ext string
	)

	for j, ext = range groupExtensions {
		if roleName == ext && instance == 0 {
			return nil, ErrAlwaysRegisterRole
		} else if roleName == ext {
			role = new(GRole)
			role.Role = GRoleType(l + j)
			role.Instance = instance
			return
		}
	}

	groupExtensions = append(groupExtensions, roleName)
	role = new(GRole)
	role.Role = GRoleType(l + j)
	role.Instance = instance
	return role, nil
}
