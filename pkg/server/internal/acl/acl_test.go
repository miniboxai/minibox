package acl

import "testing"

func buildAcl() *ACL {
	var acl = &ACL{}

	acl.RuleS("Admin", Grant{Action: AddUser, Effect: Allow})       // 管理员, 可以访问 ListUsers
	acl.RuleS("User-1234", Grant{Action: ListUsers, Effect: Allow}) // 特定用户角色, 可以访问 ListUsers

	acl.GroupRuleS("Authors", Grant{Action: ListProjects, Effect: Allow})
	acl.GroupRuleS("Contributors", Grant{Action: AddUser, Effect: Deny})

	return acl
}

func TestAcl(t *testing.T) {
	acl := buildAcl()

	if acl.Can(&Role{Role: User}, ListUsers) {
		t.Fatal("User ListUsers can't pass")
	}

	t.Logf("User can't ListUsers")

	if !acl.Can(&Role{Admin, 12}, AddUser) {
		t.Fatal("Admin AddUser can't pass")
	}

	t.Logf("Admin can AddUser")

	// RuleS("Admin", Grant{Action: UpdateUser, Effect: Allow})
	// RuleS("Admin", Grant{Action: UpdateUser, Effect: Deny})

	if acl.Can(&Role{Admin, 12}, UpdateUser) {
		t.Fatal("Admin UpdateUser not have permission")
	}

	t.Logf("Admin not have UpdateUser permission")

	acl.RuleS("Admin", Grant{Action: UpdateUser, Effect: Allow})
	// RuleS("Admin", Grant{Action: UpdateUser, Effect: Deny})

	if !acl.Can(&Role{Admin, 12}, UpdateUser) {
		t.Fatal("Admin UpdateUser should be allow")
	}

	t.Logf("Admin allow UpdateUser")

	acl.RuleS("Admin", Grant{Action: UpdateUser, Effect: Deny})

	if acl.Can(&Role{Admin, 12}, UpdateUser) {
		t.Fatal("Admin UpdateUser should be deny")
	}

	t.Logf("Admin deny UpdateUser")
}
