package server

import (
	lru "github.com/hashicorp/golang-lru"
	"minibox.ai/minibox/pkg/api/v1/types"
	"minibox.ai/minibox/pkg/server/internal/acl"
)

var (
	ruleTemplates map[acl.RoleType][]acl.Action
	aclCache      *lru.Cache
)

type ACL struct {
	*acl.ACL
}

func init() {
	ruleTemplates = map[acl.RoleType][]acl.Action{
		acl.User: []acl.Action{
			acl.ListUsers,
			acl.AddUser,
			acl.ListDatasets,
		},
		acl.Admin: []acl.Action{},
	}

	aclCache, _ = lru.New(102400)
}

func BuildACL(usr *types.User) *ACL {
	var ctl = ACL{ACL: &acl.ACL{}}

	for _, act := range ruleTemplates[acl.User] {
		ctl.Rule(acl.Role{Role: acl.User},
			acl.Grant{Action: act, Effect: acl.Allow})
	}

	return &ctl
}

type Role struct {
	usr *types.User
	acl *ACL
}

func (r *Role) Can(act acl.Action) (ok bool) {
	if val, ok := aclCache.Get(r.usr.ID); !ok {
		r.acl = BuildACL(r.usr)
		aclCache.Add(r.usr.ID, r.acl)
	} else if r.acl, ok = val.(*ACL); !ok {
		return false
	}

	allow, chk := r.acl.Check(&acl.Role{Role: acl.User}, act)
	if chk {
		return allow
	}

	return false
}
