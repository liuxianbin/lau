// Copyright (c) 2020 Lau All rights reserved.
// Use of this source code is governed by MIT License that can be found in the LICENSE file.
// Author: Lau <lauj@foxmail.com>
package helper

import (
	url2 "net/url"
	"testing"
)

func TestHttpDo(t *testing.T) {
	url := "https://aip.baidubce.com/rpc/2.0/nlp/v1/topic"
	data := url2.Values{
		"name": {"123"},
	}
	//header := map[string]string{
	//	"cookie": "name=lau;lang=golang",
	//}
	//content, err := HttpPost(url, data)
	_, err := HttpDo("POST", url, data, nil)
	if err != nil {
		t.Error(err)
	}
}
