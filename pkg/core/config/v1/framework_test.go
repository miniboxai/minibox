package v1

import (
	"testing"
)

type testframe struct {
	val    string
	failed bool
}

func (tt testframe) Exec(t *testing.T) {

	fram, err := ParseFramework(tt.val)
	if err != nil {
		if tt.failed {
			t.Logf("Test framework must failed: %s", err)
			return
		}
		t.Fatalf("Test framework failed: %s", err)
	}
	t.Logf("framework %#v", fram)
}

var tests = []testframe{
	testframe{
		"tensorflow@1.0.0",
		true,
	},
	testframe{
		"tensorflow@1.5.1",
		false,
	},
	testframe{
		"tensorflow@1.7.5",
		false,
	},
	testframe{
		"tensorflow@1.8.5",
		false,
	},
	testframe{
		"tensorflow@1.8.0#py2",
		false,
	},
	testframe{
		"torch",
		false,
	},
	testframe{
		"torch#gpu",
		false,
	},
}

type testframeDetail struct {
	name    string
	version string
	failed  bool
}

var testDetails = []testframeDetail{
	testframeDetail{
		"tensorflow",
		"1.0.0",
		true,
	},
	testframeDetail{
		"tensorflow",
		"1.5.1",
		false,
	},
	testframeDetail{
		"tensorflow",
		"1.7.5",
		false,
	},
	testframeDetail{
		"tensorflow",
		"1.8.5",
		false,
	},
	testframeDetail{
		"tensorflow",
		"1.8.0#py2",
		false,
	},
	testframeDetail{
		"torch",
		"",
		false,
	},
}

func (tt testframeDetail) Exec(t *testing.T) {
	frame, err := ParseFramework(Map{"name": tt.name, "version": tt.version})
	if err != nil {
		if tt.failed {
			t.Logf("Test framework must failed: %s", err)
			return
		}
		t.Fatalf("Test framework failed: %s", err)
	}
	t.Logf("Test framework %#v", frame)
}

func TestFrameworkString(t *testing.T) {
	for _, tt := range tests {
		tt.Exec(t)
	}
}

func TestFrameworkDetail(t *testing.T) {
	for _, tt := range testDetails {
		tt.Exec(t)
	}
}

func TestFrameworkImage(t *testing.T) {
	frame, _ := ParseFramework("tensorflow@1.8.0#py2")
	t.Logf("parse \"tensorflow@1.8.0#py2\" to Image: %s", frame.Image())

	frame, _ = ParseFramework("tensorflow@1.8.5")
	t.Logf("parse \"tensorflow@1.8.5\" to Image: %s", frame.Image())

	m := Map{"name": "tensorflow", "version": "1.8.0"}
	frame, _ = ParseFramework(m)
	t.Logf("parse map config `%v` to Image: %s", m, frame.Image())

	frame, _ = ParseFramework("torch#gpu")
	t.Logf("parse \"torch#gpu\" to Image: %s", frame.Image())

	frame, _ = ParseFramework("torch")
	t.Logf("parse \"torchu\" to Image: %s", frame.Image())
}

func TestFrameworkError(t *testing.T) {
	_, err := ParseFramework([]string{"tensorflow@1.8.0#py2"})
	if err == nil {
		t.Fatalf("[]string must be failed")
	}
	t.Logf("parse []string error %s", err)
}
