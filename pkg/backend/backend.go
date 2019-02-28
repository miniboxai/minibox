package backend

import (
	"context"

	"minibox.ai/minibox/pkg/backend/option"
	"minibox.ai/minibox/pkg/core/job"
)

type Executor interface {
	Execute(*job.Job, context.Context, *option.ExecuteOption) error
}
