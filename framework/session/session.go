package session

import (
	"net/http"
	"time"
)

type Session struct {
	ID string

	Values map[interface{}]interface{}

	Options *Options

	IsNew bool

	store Store

	name string
}

func NewSession(store Store, name string) *Session {
	session := &Session{
		name:name,
		store:store,
		IsNew:true,
		Options:new(Options),
	}

	return session
}

func (s *Session) Name() string {
	return s.name
}

// Registry -------------------------------------------------------------------

// sessionInfo stores a session tracked by the registry.
// sessionInfo 存储由registry跟踪的session
type sessionInfo struct {
	s *Session
	e error
}

// contextKey is the type used to store the registry in the context.
// contextKey 是用于在context 中存储registry的类型
type contextKey int


// registryKey is the key used to store the registry in the context.
// registryKey
const registryKey contextKey = 0

// Registry stores sessions used during a request.
type Registry struct {
	request  *http.Request
	sessions map[string]sessionInfo
}

// NewCookie returns an http.Cookie with the options set. It also sets
// the Expires field calculated based on the MaxAge value, for Internet
// Explorer compatibility.
func NewCookie(name, value string, options *Options) *http.Cookie {
	cookie := newCookieFromOptions(name,value, options)
	if cookie.MaxAge > 0 {
		d := time.Duration(cookie.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	}else {
		cookie.Expires = time.Unix(1, 0)
	}


	return cookie
}