// Copyright (c) 2020 Lau All rights reserved.
// Use of this source code is governed by MIT License that can be found in the LICENSE file.
// Author: Lau <lauj@foxmail.com>
package lau

import (
	"html/template"
	"net/http"
)

type H map[string]interface{}

type Engine struct {
	*Route
	HTMLTemplate *template.Template
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	c.Engine = e
	c.mHandlers = e.mHandlers
	c.mHandlers = append(c.mHandlers, e.handler)
	c.Next()
}

// New returns a new initialized Engine
func New() *Engine {
	return &Engine{
		Route:        NewRoute(),
		HTMLTemplate: nil,
	}
}

func Default() *Engine {
	l := New()
	l.Use(Recovery())
	return l
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}

// ParseFiles from the named files
func (e *Engine) LoadHTMLFiles(filenames ...string) {
	t := template.Must(template.ParseFiles(filenames...))
	e.HTMLTemplate = t
}

// ParseFiles with the list of files matched by the pattern
func (e *Engine) LoadHTMLGlob(pattern string) {
	t := template.Must(template.ParseGlob(pattern))
	e.HTMLTemplate = t
}
