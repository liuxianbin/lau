// Copyright (c) 2020 Lau All rights reserved.
// Use of this source code is governed by MIT License that can be found in the LICENSE file.
// Author: Lau <lauj@foxmail.com>
package helper

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"sort"
	"strings"
	"time"
)

type jwtData map[string]string

func (a jwtData) bytes() []byte {
	var keys []string
	var content string
	for k, _ := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sep string
	for _, k := range keys {
		content += sep + k + "=" + a[k]
		sep = ","
	}
	return []byte(content)
}

func parseAccessData(data string) jwtData {
	result := make(jwtData)
	parts := strings.Split(data, ",")
	for _, item := range parts {
		item := strings.Split(item, "=")
		result[item[0]] = item[1]
	}
	return result
}

type jwt struct {
	key    string
	header jwtData
}

func newAccess(key string) *jwt {
	return &jwt{key: key, header: jwtData{"alg": "sha256"}}
}

func (a jwt) signature(data string) string {
	key := []byte(a.key)
	var f func() hash.Hash
	switch a.header["alg"] {
	case "sha256":
		f = sha256.New
	case "md5":
		f = md5.New
	default:
		f = sha256.New
	}
	mac := hmac.New(f, key)
	mac.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

func (a jwt) getToken(payload jwtData) string {
	b64Header := base64.URLEncoding.EncodeToString(a.header.bytes())
	b64Payload := base64.URLEncoding.EncodeToString(payload.bytes())
	data := b64Header + "." + b64Payload
	return data + "." + a.signature(data)
}

func VerifyJWT(key, token string) (jwtData, bool) {
	tokens := strings.Split(token, ".")
	if len(tokens) != 3 {
		return nil, false
	}
	var b64DecodeHeader, decodePayload []byte
	var err error
	if b64DecodeHeader, err = base64.URLEncoding.DecodeString(tokens[0]); err != nil {
		return nil, false
	}
	b64Header := parseAccessData(string(b64DecodeHeader))
	if _, ok := b64Header["alg"]; !ok {
		return nil, false
	}
	if decodePayload, err = base64.URLEncoding.DecodeString(tokens[1]); err != nil {
		return nil, false
	}
	payload := parseAccessData(string(decodePayload))
	now := time.Now().Format("2006-01-02 15:04:05")
	// create time
	if iat, ok := payload["createTime"]; ok {
		if iat > now {
			return nil, false
		}
	}
	// expire time
	if exp, ok := payload["expireTime"]; ok {
		if exp < now {
			return nil, false
		}
	}
	a := newAccess(key)
	if a.signature(tokens[0]+"."+tokens[1]) != tokens[2] {
		return nil, false
	}
	return payload, true
}

// JWT
func JWT(key string, payload jwtData) string {
	a := newAccess(key)
	return a.getToken(payload)
}
