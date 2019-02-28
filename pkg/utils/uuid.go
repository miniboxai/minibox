package utils

import (
	uuid "github.com/satori/go.uuid"
)

func UUID() string {
	uid := uuid.NewV4()
	return uid.String()
}
