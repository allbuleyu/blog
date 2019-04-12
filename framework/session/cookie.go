package session

import "net/http"

func newCookieFromOptions(name,value string, options *Options) *http.Cookie {
	return &http.Cookie{
		Name:name,
		Value:value,

		Domain:options.Domain,
		MaxAge:options.MaxAge,
		Path:options.Path,
		Secure:options.Secure,
		HttpOnly:options.HttpOnly,
	}
}
