package acl

import "testing"

func TestExtensionRole(t *testing.T) {
	role, err := RegisterRole("Engine")
	if err != nil {
		t.Fatalf("register role %s", err)
	}
	t.Logf("Engine\t\trole: %#v", role)

	role, err = RegisterRole("Engine-1234")
	if err != nil {
		t.Fatalf("register role %s", err)
	}
	t.Logf("Engine\t\trole: %#v", role)

	role, err = RegisterRole("Admin-321")
	if err != nil {
		t.Fatalf("register role %s", err)
	}
	t.Logf("Admin\t\trole: %#v", role)

	role, err = RegisterRole("User-321")
	if err != nil {
		t.Fatalf("register role %s", err)
	}
	t.Logf("User\t\trole: %#v", role)

	role, err = RegisterRole("Organization")
	if err != nil {
		t.Fatalf("register role %s", err)
	}
	t.Logf("Organization\t\trole: %#v", role)

	role, err = RegisterRole("Staff")
	if err != nil {
		t.Fatalf("register role %s", err)
	}
	t.Logf("Staff\t\trole: %#v", role)

	role, err = RegisterRole("Admin")
	if err != nil {
		t.Fatalf("register role %s", err)
	}
	t.Logf("Admin\t\trole: %#v", role)
}

func TestExtensionAction(t *testing.T) {
	action, err := RegisterAction("OpenBook")
	if err != nil {
		t.Fatalf("register action %s", err)
	}
	t.Logf("OpenBook\t\taction: %s", action)

	action, err = RegisterAction("ListUsers")
	if err != nil {
		t.Fatalf("register action %s", err)
	}
	t.Logf("ListUsers\t\taction: %s", action)

	action, err = RegisterAction("AddUser")
	if err != nil {
		t.Fatalf("register action %s", err)
	}
	t.Logf("AddUser\t\taction: %s", action)
}
