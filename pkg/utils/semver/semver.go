package semver

import (
	sver "github.com/blang/semver"
)

type Version = sver.Version
type Range = sver.Range

func Make(s string) (Version, error) {
	return sver.Make(s)
}

func New(s string) (vp *Version, err error) {
	return sver.New(s)
}

func ParseRange(s string) (Range, error) {
	return sver.ParseRange(s)
}
