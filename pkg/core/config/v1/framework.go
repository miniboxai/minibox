package v1

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	gerr "github.com/pkg/errors"

	u "minibox.ai/minibox/pkg/core/utils"
	"minibox.ai/minibox/pkg/utils/semver"
)

type Framework struct {
	Lang     string
	Name     string
	Version  string
	Sense    string
	Requires []Require

	val       interface{}
	tag       string
	full      string
	cacheFull bool
	image     string
}

func errorHandle(handle func(error)) {
	if err, ok := recover().(error); ok {
		handle(err)
	}
}

func ParseFramework(val interface{}) (*Framework, error) {
	fram := &Framework{val: val}
	err := fram.Parse()
	return fram, err
}

func (frm *Framework) defaultVersion(name string) string {
	full, _ := frm.lookupFramework(name)
	return u.VMap(defaultLanguages).Sub(full+".tags", ".").String("default")
}

func (frm *Framework) lookupFramework(name string) (string, bool) {
	if !frm.cacheFull {
		full, ok := lookupFramework(name)
		if ok {
			frm.cacheFull = true
			frm.full = full
		}
	}

	return frm.full, frm.cacheFull
}

func (frm *Framework) Eval(ctx context.Context) (reflect.Value, error) {

	return reflect.Value{}, nil
}

func (frm *Framework) Parse() (err error) {
	defer errorHandle(func(_err error) {
		err = gerr.Wrap(_err, "Framework Parse error")
	})

	switch v := frm.val.(type) {
	case string:
		_, tag := frm.parseString(v)
		frm.tag = tag
	case Map:
		_, tag := frm.parseMap(v)
		frm.tag = tag
	default:
		return &ErrInvalidStruct{"train.framework", frm.val, "string or object"}
	}

	// m := u.MustMap(defaultLanguages, frm.full)
	frm.Lang = u.VMap(defaultLanguages).Map(frm.full).String("lang")
	return nil
}

func (frm *Framework) parseSense(raw string) (string, string, bool) {
	var (
		sects = strings.SplitN(raw, "#", 2)
	)

	if len(sects) == 2 {
		return sects[0], sects[1], true
	} else {
		return sects[0], "", false
	}
}

func (frm *Framework) parseString(conf string) (ver string, tag string) {
	var (
		frame, sense, full, image string
		ok                        bool
	)

	frame, ver, sense = frm.parseTagString(conf)
	full, ok = frm.lookupFramework(frame)

	if !ok {
		panic(fmt.Errorf("invalid framework: %s", frame))
	}

	frm.Name = full

	image, tag, ok = frm.adjustTag(ver, sense)

	if ok && (tag == "unavailable" || empty(tag)) {
		panic(fmt.Errorf("this is '%s' version is unavailable.", ver))
	} else if ok {
		frm.image = image
		return ver, tag
	}

	return
}

func (frm *Framework) parseTagString(conf string) (frame, ver, sense string) {
	var (
		sects = strings.SplitN(conf, "@", 2)
	)

	if len(sects) == 2 {
		frame = sects[0]
		ver, sense, _ = frm.parseSense(sects[1])
	} else if len(sects) == 1 {
		frame = sects[0]
		frame, sense, _ = frm.parseSense(sects[0])
		ver = frm.defaultVersion(frame)
	} else {
		panic(fmt.Errorf("invalid framework settings"))
	}

	return
}

func (frm *Framework) parseMap(m Map) (sver, tag string) {
	var (
		v                    = u.VMap(m)
		name                 = v.String("name")
		full, version, sense string
		ok                   bool
	)

	if full, ok = frm.lookupFramework(name); !ok {
		panic(&ErrNotSupportFramework{name})
	} else {
		frm.Name = full
		version, _ = v["version"].(string)
		sense, _ = v["sense"].(string)

		if requires, ok := v["requires"]; ok {
			if reqs, ok := frm.parseRequires(requires); ok {
				frm.Requires = reqs
			}
		}
		return frm.parseVersion(version, sense)
	}
	return
}

func (frm *Framework) parseVersion(vers, sense string) (sver, tag string) {
	var (
		image string
		ok    bool
	)

	if empty(vers) {
		vers = frm.defaultVersion(frm.Name)
	}

	image, tag, ok = frm.adjustTag(vers, sense)
	if ok && (tag == "unavailable" || empty(tag)) {
		panic(fmt.Errorf("this is '%s' version is unavailable.", vers))
	} else if ok {
		frm.image = image
		return vers, tag
	}

	return
}

func (frm *Framework) parseRequires(val interface{}) ([]Require, bool) {
	return nil, false
}

func (frm *Framework) adjustTag(sver, sense string) (image, tag string, ok bool) {
	var (
		tagconf  string
		versions = loadVersions(frm.Name, sense)
	)
	tagconf, ok = frm.validVersions(sver, versions)
	if !ok {
		return "", "", false
	}

	image, tag = frm.splitTag(tagconf)
	return
}

func (frm *Framework) validVersions(sver string, versions []string) (tagconf string, ok bool) {
	var (
		ver semver.Version
		cmp = func(rng semver.Range) bool { return rng(ver) }
		err error
	)

	if sver == "latest" {
		tagconf = versions[len(versions)-1]
		ss := strings.SplitN(tagconf, "=>", 2)
		if len(ss) == 2 {
			return ss[1], true
		}

		return "", false
	} else if ver, err = semver.Make(sver); err != nil {
		return "", false
	}

	for _, rv := range versions {
		ss := strings.SplitN(rv, "=>", 2)

		if tagconf, ok = rangec(ss, cmp); ok {
			return tagconf, true
		}
	}

	return "", false
}

func (frm *Framework) splitTag(stag string) (string, string) {
	ss := strings.SplitN(stag, ":", 2)
	if len(ss) == 2 {
		return ss[0], ss[1]
	} else {
		return "", ss[0]
	}
}

func (frm *Framework) Image() string {
	var img string
	if empty(frm.image) {
		frame := defaultLanguages[frm.Name].(Map)
		img = u.MustString(frame["image"], "image")
	} else {
		img = frm.image
	}

	return fmt.Sprintf("%s:%s", img, frm.tag)
}

func (frm *Framework) Defaults() Map {
	return defaultLanguages[frm.Name].(Map)
}

func pureName(raw string) string {
	ss := strings.SplitN(raw, "#", 2)
	if len(ss) == 2 {
		return ss[0]
	} else if len(ss) == 1 {
		return ss[0]
	} else {
		return raw
	}
}

func lookupFramework(fw string) (string, bool) {
	fw = pureName(fw)
	for fram, val := range defaultLanguages {
		fram := u.MustString(fram, "fram")

		if strings.EqualFold(fram, fw) {
			return fram, true
		}

		if v, ok := val.(Map); ok {
			if alias, ok := v["alias"].([]interface{}); ok {
				for _, nam := range alias {
					nam, _ := nam.(string)
					if strings.EqualFold(nam, fw) {
						return fram, true
					}
				}
			}
		}
	}

	return "", false
}

type validRange func(semver.Range) bool

func convImgAndTag(beconv string) (img, tag string) {
	if strings.Index(beconv, ":") > 0 {
		ss := strings.SplitN(beconv, ":", 2)
		img = ss[0]
		tag = ss[1]
	} else {
		tag = beconv
	}

	return
}

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func loadVersions(fram string, sense string) []string {
	var versions []string
	frame := u.VMap(defaultLanguages[fram].(Map))
	if empty(sense) {
		versions = frame.StringSlice("versions")
	} else {
		versions = frame.StringSlice("versions_" + sense)
	}
	return versions
}

func parseRange(s string) semver.Range {
	if rn, err := semver.ParseRange(s); err != nil {
		panic(err)
	} else {
		return rn
	}
}

func rangec(secs []string, check validRange) (string, bool) {
	if len(secs) == 2 {
		if check(parseRange(secs[0])) {
			return secs[1], true
		} else {
			return "", false
		}
	} else {
		panic(errors.New("[]string slice length must equal 2"))
	}
	return "", false
}
