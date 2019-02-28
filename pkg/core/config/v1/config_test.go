package v1

import (
	"testing"

	"github.com/kr/pretty"
	"minibox.ai/pkg/core/option"
)

type testConfig struct {
	ProjectName string
	Content     []byte
	t           *testing.T
}

func TestConfig(t *testing.T) {
	testCfg(testConfig{
		ProjectName: "test",
		Content: []byte(`---
train:
  framework: torch@2.1.7
`),
		t: t,
	})

	testCfg(testConfig{
		"test2",
		[]byte(`---
train:
  framework: tf@1.8.0
  notebook: true
  workdir: /usr/src/{projectName}
  dirs:
    data: data/
    log: logdir/
  dataset:
    - mnist/test.tfrecords
    - mnist/validation.tfrecords    
    - mnist/train.tfrecords
  pretrain:
    - mnist/models.ckpt
  dashboard: true
  cmd: "scripts/start-dist-mnist.sh {datadir} {logdir}"
`), t})

	testCfg(testConfig{
		"test3",
		[]byte(`---
train:
  framework: tf@1.8.0
  cmd:
    - python
    - models.py
`), t})
}

func testCfg(t testConfig) {
	cfg, err := Read(t.Content, option.CProject(t.ProjectName))
	if err != nil {
		t.t.Fatalf("error %s", err)
	}

	ctx := NewContext(cfg.opts.ProjectName, nil)
	cfg.Eval(ctx)

	t.t.Logf("content:\n%s\n", t.Content)
	t.t.Logf("filename: %s", cfg.opts.ProjectName)
	t.t.Logf("config: % #v\n", pretty.Formatter(cfg))
}
