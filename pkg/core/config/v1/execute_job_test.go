package v1

import (
	"testing"

	"github.com/kr/pretty"
	"minibox.ai/pkg/core/option"
)

type testTrainer struct {
	ProjectName string
	Content     []byte
	t           *testing.T
}

func TestTrainer(t *testing.T) {
	testTrain(testTrainer{
		"project",
		[]byte(`---
train:
  framework: torch#gpu
`), t})

	testTrain(testTrainer{
		"project",
		[]byte(`---
train:
  framework: tensorflow@1.8.0#py2
`), t})
}

func testTrain(t testTrainer) {
	cfg, err := Read(t.Content, option.CProject(t.ProjectName))
	if err != nil {
		t.t.Fatalf("error %s", err)
	}

	ctx := NewContext(cfg.opts.ProjectName, nil)
	job, err := cfg.Start(ctx, nil)

	// cfg.Eval(ctx)

	t.t.Logf("job:\n % #v\n", pretty.Formatter(job))
}
