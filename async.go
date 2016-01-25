package main

import (
	"net/http"
)

type ASync struct {
}

func NewASync() *ASync {
	return &ASync{}
}

func (a *ASync) ServeHTTP(ow http.ResponseWriter, r *http.Request) {
}
