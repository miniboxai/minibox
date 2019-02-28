package apiserver

import (
	validator "gopkg.in/go-playground/validator.v9"
)

// Form 渲染辅助类
type CreateNamespaceForm struct {
	Action     string
	Namespace  string `schema:"namespace,required" validate:"alphanum,required,min=5"`
	Name       string `schema:"name,required"`
	Email      string `schema:"email,required"`
	Profile    string `schema:"profile"`
	Picture    string `schema:"picture,required"`
	Sub        string `schema:"sub,required"`
	FamilyName string `schema:"family_name"`
	GivenName  string `schema:"given_name"`
	Gender     string `schema:"gender"`
}

func validateForm(obj interface{}) error {
	validate := validator.New()
	return validate.Struct(obj)
}

func rangeValidationErrors(err error) []validator.FieldError {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	return errs
}
