package acl

import (
	"encoding/json"
	"time"
)

// 作用
type Effect int
type Action int
type RoleType int
type GroupRole int
type GRoleType int

//go:generate stringer -type RoleType
//go:generate stringer -type Action
//go:generate stringer -type Effect
//go:generate stringer -type GRoleType
const (
	NotApplicable Effect = iota // 无设置
	Allow                       // 允许
	Deny                        // 禁止
)

const (
	None Action = iota
	ListUsers
	AddUser
	UpdateUser
	DeleteUser
	ListProjects
	ListPrivateDatasets

	CreateDataset
	ListDatasets
	DeleteDataset
)

const (
	User RoleType = iota
	Organization
	Staff
	Admin
)

const (
	Authors GRoleType = iota
	Contributors
)

type ACL struct {
	// sync.RWMutex
	Rules      _rules
	GroupRules []_group_rules
	Timestamp  time.Time
}

// 角色
type Role struct {
	Role     RoleType
	Instance int
}

// 组角色
type GRole struct {
	Role     GRoleType
	Instance int
}

// 组
type Group struct {
	Role  Role
	Roles []Role
}

type Grant struct {
	Action   Action // 权限类型
	Effect   Effect
	Resource []string // 作用资源
}

type rule struct {
	Role  Role
	Grant Grant
}

type groupRule struct {
	Role   GRole
	Grants []Grant
}

type _rules struct {
	Allows []*rule
	Denies []*rule
}

type _grules struct {
	Allows []Action
	Denies []Action
}

type _group_rules struct {
	Role GRole
	// Members []Role
	Rules _grules
}

type CanOption func(*Option)

// var globalRules = &_rules{}

// func Can(role *Role, action Action, opts ...CanOption) bool {
// 	for _, act := range globalRules.Denies {
// 		if act == action {
// 			return false
// 		}
// 	}

// 	for _, act := range globalRules.Allows {
// 		if act == action {
// 			return true
// 		}
// 	}

// 	return false
// }

func (a *ACL) Rule(role Role, grant Grant) {
	var rule = &rule{Role: role, Grant: grant}

	if grant.Effect == Deny {
		a.Rules.Denies = append(a.Rules.Denies, rule)
	} else {
		a.Rules.Allows = append(a.Rules.Allows, rule)
	}
}

func (a *ACL) RuleS(role string, grant Grant) {
	r := buildRole(role)
	a.Rule(*r, grant)
}

func (a *ACL) GroupRule(role GRole, grant Grant) {
	// var rule = &rule{Role: role, Grant: grant}

	for i, rules := range a.GroupRules {
		if rules.Role.Same(&role) {
			if grant.Effect == Deny {
				a.GroupRules[i].Rules.Denies = append(a.GroupRules[i].Rules.Denies, grant.Action)
			} else {
				a.GroupRules[i].Rules.Allows = append(a.GroupRules[i].Rules.Allows, grant.Action)
			}
			return
		}
	}
	var grules _grules

	if grant.Effect == Deny {
		grules.Denies = append(grules.Denies, grant.Action)
	} else {
		grules.Allows = append(grules.Allows, grant.Action)
	}

	a.GroupRules = append(a.GroupRules,
		_group_rules{Role: role, Rules: grules})
}

func (a *ACL) GroupRuleS(role string, grant Grant) {
	r := buildGRole(role)
	a.GroupRule(*r, grant)
}

func (a *ACL) Can(role *Role, action Action, opts ...CanOption) bool {
	allow, _ := a.Check(role, action)
	return allow
}

func (a *ACL) Check(role *Role, action Action) (bool, bool) {

	for _, rule := range a.Rules.Denies {
		if rule.Match(role, action, true) {
			return false, true
		}
	}

	for _, rule := range a.Rules.Allows {
		if rule.Match(role, action, false) {
			return true, true
		}
	}

	return false, false
}

func (a *ACL) CanGroup(role *GRole, action Action, opts ...CanOption) bool {
	allow, _ := a.CheckGroup(role, action)
	return allow
}

func (a *ACL) CheckGroup(role *GRole, action Action) (bool, bool) {
	for _, grp := range a.GroupRules {
		if grp.Role.Same(role) {
			for _, act := range grp.Rules.Denies {
				if act == action {
					return false, true
				}
			}

			for _, act := range grp.Rules.Allows {
				if act == action {
					return true, true
				}
			}
		}
	}

	return false, false
}

func (a *ACL) applyOpt(opts []CanOption) *Option {
	var opt = new(Option)
	for _, op := range opts {
		op(opt)
	}

	return opt
}

func (a *ACL) Load(buf []byte, timestamp time.Time) error {
	var m map[string]interface{}
	if err := json.Unmarshal(buf, &m); err != nil {
		return err
	}

	a.Timestamp = timestamp

	return nil
}

func (a *ACL) Save() ([]byte, error) {
	return nil, nil
}

// func Rule(role Role, grant Grant) {
// 	var rule = &rule{Role: role, Grant: grant}

// 	if grant.Effect == Deny {
// 		globalRules.Denies = append(globalRules.Denies, rule)
// 	} else {
// 		globalRules.Allows = append(globalRules.Allows, rule)
// 	}
// }

// func RuleS(role string, grant Grant) {
// 	r := buildRole(role)
// 	Rule(*r, grant)
// }

// func init() {

// 	RuleS("Admin", Grant{Action: AddUser, Effect: Allow})       // 管理员, 可以访问 ListUsers
// 	RuleS("User-1234", Grant{Action: ListUsers, Effect: Allow}) // 特定用户角色, 可以访问 ListUsers
// 	// RuleGroupS("Sales-1234", Grant{Action: ListUsers, Effect: Allow})
// 	// log.Printf("rules :%# v", pretty.Formatter(globalRules))
// }
