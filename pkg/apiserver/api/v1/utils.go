package v1

import (
	"context"

	"minibox.ai/pkg/models"
)

func getCurrntUser(ctx context.Context) (*models.User, bool) {
	if cur_usr, ok := ctx.Value("user").(*models.User); ok {
		return cur_usr, true
	} else {
		return nil, false
	}

}
