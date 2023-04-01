package http

import "net/http"

func CheckRequestContentType(headers http.Header, wantType string) bool {
	ctVal, ok := headers["Content-Type"]
	if !ok {
		return false
	}
	for _, v := range ctVal {
		if v == wantType {
			return true
		}
	}
	return false
}
