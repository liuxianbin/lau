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
