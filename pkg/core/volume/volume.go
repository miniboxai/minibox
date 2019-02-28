package volume

type Sourcer interface {
	Source() string
}

type LocalVolume struct {
	Path string
}

func (vol *LocalVolume) Source() string {
	return vol.Path
}

type ProjectVolume struct {
	Name     string
	RootPath string
}

func (vol *ProjectVolume) Source() string {
	return vol.RootPath
}

type DatasetVolume struct {
	Type int
	Path string
}

func (vol *DatasetVolume) Source() string {
	return vol.Path
}
