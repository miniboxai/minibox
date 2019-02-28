package option

type ConfigOpt struct {
	Strict      bool
	Slient      bool
	RootPath    string
	ProjectName string
}

type ParseOpt func(*ConfigOpt)

type JobOpt struct {
	Wait bool
}

type ExecuteOpt struct {
	Log bool
}

func CProject(name string) ParseOpt {
	return func(opt *ConfigOpt) {
		opt.ProjectName = name
	}
}
