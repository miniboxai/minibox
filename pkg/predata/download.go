package predata

import (
	"os"

	getter "github.com/hashicorp/go-getter"
)

var pwd string

func Setwd(p string) {
	pwd = p
}

func Download(dest, src string) error {
	src, err := getter.Detect(src, pwd, getter.Detectors)
	if err != nil {
		return err
	}

	return getter.Get(dest, src)
}

func init() {
	p, _ := os.Getwd()
	Setwd(p)
}
