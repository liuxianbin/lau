// Copyright (c) 2020 Lau All rights reserved.
// Use of this source code is governed by MIT License that can be found in the LICENSE file.
// Author: Lau <lauj@foxmail.com>
package lau

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestRoute_GET(t *testing.T) {
	r := NewRoute()
	r.GET("/hello/:name/show", func(c *Context) {
		c.String(200, "%v\n", "name="+c.Param("name"))
	})
	r.GET("/static/*filepath", func(c *Context) {
		c.String(200, "%v\n", "filepath="+c.Param("filepath"))
	})
	for _, item := range []string{"/hello/go/show", "/static/a/b/c"} {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", item, nil)
		c := NewContext(res, req)
		r.handler(c)
		if res.Code != 200 {
			t.Errorf("Response code is %v\n", res.Code)
			return
		}
		fmt.Println(res.Body.String())
	}
}

func ExampleSplitattern() {
	fmt.Println(splitPattern("/hello"))
	fmt.Println(splitPattern("/hello/"))
	fmt.Println(splitPattern("hello"))
	fmt.Println(splitPattern("/hello/:name/aa/:info"))
	fmt.Println(splitPattern("/hello/*"))
	fmt.Println(splitPattern("/hello/*abc"))
	// output:
	// [hello]
	// [hello]
	// [hello]
	// [hello :name aa :info]
	// [hello *]
	// [hello *abc]
}
