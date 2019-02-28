package cmd

import "github.com/parnurzeal/gorequest"

type Request struct {
	*gorequest.SuperAgent
}

func NewRequest() *Request {
	return &Request{
		SuperAgent: gorequest.New(),
	}
}
