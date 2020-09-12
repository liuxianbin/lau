package helper

import (
	"fmt"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	key := "key"
	str := JWT(key, map[string]string{
		"lang": "golang",
		"exp":  time.Now().Format("2006-01-02 15:04:05"),
	})
	fmt.Println(str)
	b, ok := VerifyJWT(key, str)
	if !ok {
		t.Errorf("%s", "check error")
	} else {
		fmt.Println(b, ok)
	}
}
