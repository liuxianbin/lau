// Copyright (c) 2020 Lau All rights reserved.
// Use of this source code is governed by MIT License that can be found in the LICENSE file.
// Author: Lau <lauj@foxmail.com>
package lau

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Method string
	Path   string
	Code   int
	*Engine
	Params map[string]string
	// middleware
	mHandlers []handleFunc
	mIndex    int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{W: w, R: r, Method: r.Method, Path: r.URL.Path, mIndex: -1}
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("content-type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.W, err.Error(), 500)
	}
}

func (c *Context) Error(code int, err interface{}) {
	c.Code = code
	http.Error(c.W, fmt.Sprintf("%s\n", err), code)
}

func (c *Context) SetStatus(code int) {
	c.Code = code
	c.W.WriteHeader(code)
}

func (c *Context) SetHeader(k, v string) {
	c.W.Header().Set(k, v)
}

func (c *Context) String(code int, format string, val string) {
	c.SetHeader("content-type", "text/plain")
	c.SetStatus(code)
	fmt.Fprintf(c.W, format, val)
}

func (c *Context) HTML(code int, filename string, obj interface{}) {
	c.SetHeader("content-type", "text/html")
	c.SetStatus(code)
	if err := c.HTMLTemplate.ExecuteTemplate(c.W, filename, obj); err != nil {
		http.Error(c.W, err.Error(), 500)
	}
}

// get bind parameters by dynamic route
func (c *Context) Param(k string) string {
	v, _ := c.Params[k]
	return v
}

func (c *Context) FormValue(k string) string {
	return c.R.FormValue(k)
}

func (c *Context) PostFormValue(k string) string {
	return c.R.PostFormValue(k)
}

func (c *Context) Next() {
	c.mIndex++
	for ; c.mIndex < len(c.mHandlers); {
		c.mHandlers[c.mIndex](c)
		c.mIndex++
	}
}
