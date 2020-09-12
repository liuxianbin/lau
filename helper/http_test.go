package helper

import (
	"fmt"
	url2 "net/url"
	"testing"
)

func TestHttpDo(t *testing.T) {
	url := "http://www.01happy.com/demo/accept.php"
	data := url2.Values{
		"lang": {"golang"},
		"name": {"lau"},
	}
	header := map[string]string{
		"cookie": "name=lau;lang=golang",
	}
	fmt.Println(data.Encode())
	content, err := HttpDo("POST", url, data, header)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(content))
}
