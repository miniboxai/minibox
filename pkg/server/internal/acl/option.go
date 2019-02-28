package acl

type Option struct {
	Resources []string
	Check     bool
}

func Check() CanOption {
	return func(opt *Option) {
		opt.Check = true
	}
}
