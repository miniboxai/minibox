package acl

func (r *rule) SameRole(role *Role) bool {
	return r.Role.Same(role)
}

func (r *rule) Match(role *Role, action Action, deny bool) bool {

	if r.SameRole(role) && action == r.Grant.Action {
		if deny {
			return r.Grant.Effect == Deny
		} else {
			return r.Grant.Effect == Allow
		}
	} else {
		return false
	}
}
