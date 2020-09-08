// | Copyright (c) 2020 Lau All rights reserved.
// | Use of this source code is governed by MIT License that can be found in the LICENSE file.
// | Author: Lau <lauj@foxmail.com>
package lau

import (
	"net/http"
	"strings"
)

type handleFunc func(*Context)

type Route struct {
	handlers map[string]handleFunc
	trees    map[string]*node
}

func NewRoute() *Route {
	return &Route{handlers: make(map[string]handleFunc), trees: make(map[string]*node)}
}

func (r *Route) GET(path string, handler handleFunc) {
	r.addRoute("GET", path, handler)
}

func (r *Route) POST(path string, handler handleFunc) {
	r.addRoute("POST", path, handler)
}

func (r *Route) addRoute(method, path string, handler handleFunc) {
	parts := splitPattern(path)
	root := r.trees[method]
	if root == nil {
		root = &node{}
		r.trees[method] = root
	}
	root.insertChild(path, parts, 0, handler)
}

// FileServer returns a handler that serves HTTP requests with the contents of the file system rooted at root.
func (r *Route) Static(prefix, root string) {
	fileHandler := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
	r.GET(prefix+"/*", func(c *Context) {
		fileHandler.ServeHTTP(c.W, c.R)
	})
}

// MUX Handler
func (r *Route) handler(c *Context) {
	root := r.trees[c.Method]
	parts := splitPattern(c.Path)
	handler, params := root.getHandler(parts, 0)
	if handler != nil {
		c.Params = params
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

// node type
const (
	NODEDEFAULT  uint8 = iota // normal string
	NODEPARAM                 // bind param
	NODEWILDCARD              // wildcard
)

type node struct {
	path     string
	section  string
	nType    uint8
	children []*node
	handler  handleFunc
}

//split url path into slice
func splitPattern(pattern string) []string {
	parts := make([]string, 0)
	for _, item := range strings.Split(strings.TrimLeft(pattern, "/"), "/") {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func getNodeType(part string) uint8 {
	if len(part) == 0 {
		return NODEDEFAULT
	}
	switch part[0] {
	case ':':
		return NODEPARAM
	case '*':
		return NODEWILDCARD
	default:
		return NODEDEFAULT
	}
}

// find matched route node
func (n *node) getHandler(parts []string, index int) (handleFunc, map[string]string) {
	params := make(map[string]string)
	if index == len(parts) || n.nType == NODEWILDCARD {
		params = getBindParams(n.path, parts)
		return n.handler, params
	}
	section := parts[index]
	childs := n.matchChilds(section)
	for _, child := range childs {
		handler, p := child.getHandler(parts, index+1)
		if handler != nil {
			return handler, p
		}
	}
	return nil, nil
}

func getBindParams(path string, parts []string) map[string]string {
	sections := splitPattern(path)
	params := make(map[string]string)
	for i, r := range sections {
		nType := getNodeType(r)
		switch nType {
		case NODEPARAM:
			params[r[1:]] = parts[i]
		case NODEWILDCARD:
			params[r[1:]] = strings.Join(parts[i:], "/")
		}
	}
	return params
}

// insert trie tree node
func (n *node) insertChild(pattern string, parts []string, index int, handler handleFunc) {
	section := parts[index]
	child := n.matchChild(section)
	if child == nil {
		child = &node{
			section: section,
			nType:   getNodeType(section),
		}
		n.children = append(n.children, child)
	}
	if index == len(parts)-1 {
		child.path = pattern
		child.handler = handler
		return
	}
	child.insertChild(pattern, parts, index+1, handler)
}

// match one
func (n *node) matchChild(section string) *node {
	for _, child := range n.children {
		if child.section == section || child.nType == NODEPARAM || child.nType == NODEWILDCARD {
			return child
		}
	}
	return nil
}

// match all
func (n *node) matchChilds(section string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.section == section || child.nType == NODEPARAM || child.nType == NODEWILDCARD {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
