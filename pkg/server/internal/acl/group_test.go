package acl

import (
	"testing"
)

func TestGroup(t *testing.T) {
	acl := buildAcl()

	if !acl.CanGroup(&GRole{Role: Authors}, ListProjects) {
		t.Fatal("Authors group must ListProjects")
	}

	t.Logf("Authors group can ListProjects")

	if acl.CanGroup(&GRole{Role: Authors}, ListUsers) {
		t.Fatal("Authors can't pass ListUsers")
	}

	t.Logf("Authors group can't ListUsers")

	if acl.CanGroup(&GRole{Role: Contributors}, AddUser) {
		t.Fatal("Contributors can't pass AddUser")
	}

	t.Logf("Contributors group can't AddUser")
}
