package controllers

import (
	"github.com/gorilla/schema"
	"net/http"
)

func parse(r *http.Request, dst interface{}) error {
	decoder := schema.NewDecoder()
	if err := decoder.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
