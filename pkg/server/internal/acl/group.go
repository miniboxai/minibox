package acl

func buildGRole(role string) *GRole {
	if r, err := RegisterGRole(role); err != nil {
		panic(err)
	} else {
		return r
	}
}

func (gr *GRole) Same(role *GRole) bool {
	if gr.Instance > 0 {
		return gr.Role == role.Role &&
			gr.Instance == role.Instance
	} else {
		return gr.Role == role.Role
	}
}
