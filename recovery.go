// Copyright (c) 2020 Lau All rights reserved.
// Use of this source code is governed by MIT License that can be found in the LICENSE file.
// Author: Lau <lauj@foxmail.com>
package lau

import (
	"log"
)

func Recovery() handleFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				c.mIndex = len(c.mHandlers)
				log.Println(err)
				c.Error(500, err)
			}
		}()
		c.Next()
	}
}
