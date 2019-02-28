package v1

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
	"minibox.ai/minibox/pkg/core/option"
)

var defaultLanguages map[interface{}]interface{}

func Read(buf []byte, opts ...option.ParseOpt) (*Config, error) {
	var cfg = &Config{}

	m := make(map[interface{}]interface{})

	err := yaml.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}

	if err = cfg.Parse(m, opts...); err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return Read(buf)
}

//go:generate go-bindata -pkg v1 assets/
func loadLanguagesConfig() {
	m := make(map[interface{}]interface{})

	// _, filename, _, _ := runtime.Caller(1)
	// dir := filepath.Dir(filename)

	f := bytes.NewBuffer(MustAsset("assets/framework.defaults.yaml"))
	// f, err := os.Open(path.Join(dir, "framework.defaults.yaml"))
	// if err != nil {
	// 	panic(err)
	// }

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}

	langs, ok := m["languages"].(map[interface{}]interface{})
	if !ok {
		panic(errors.New("don't have the 'languages'; section"))
	}

	defaultLanguages = langs
}

func (cfg *Config) loadOpts(opts ...option.ParseOpt) {
	cfg.opts = new(option.ConfigOpt)
	for _, op := range opts {
		op(cfg.opts)
	}
}

func init() {
	loadLanguagesConfig()
}
