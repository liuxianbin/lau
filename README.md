# Lau Web Framework

Lau is a web framework written in Golang.


## Installation

To install package

```sh
$ go get -u gitee.com/lauj/lau
```

Import it in your code:

```go
import "gitee.com/lauj/lau"
```


## Quick start

```go
package main

import "gitee.com/lauj/lau"

func main() {
    l := lau.Default()
    l.Static("/static", "static")
    l.GET("/hello/:one/:two/info", func(c *lau.Context) {
        c.HTML(200, "views/two.tmpl", nil)
    })
    l.GET("/hello", func(c *lau.Context) {
        c.String(200, "%v", "hello")
    })
    l.GET("/demo/*aaa", func(c *lau.Context) {
        c.JSON(200, lau.H{
            "info": "info",
        })
        fmt.Println(c.Params)
    })
    l.Use(func(c *lau.Context) {
        c.String(200, "%v\n", "middleware")
    })
    l.Run(":8000")
}
```
