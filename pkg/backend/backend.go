package backend

import (
	"context"

	"minibox.ai/pkg/backend/option"
	"minibox.ai/pkg/core/job"
)

type Executor interface {
	Execute(*job.Job, context.Context, *option.ExecuteOption) error
}
