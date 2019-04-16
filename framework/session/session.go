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
		Values: map[interface{}]interface{}{},
	}

	return session
}

// default session key
var flashesKey = "_session_key"

func (s *Session ) AddFlash(value interface{}, vars ...string) []interface{} {
	key := flashesKey
	if len(vars) > 0 {
		key = vars[0]
	}
	var flashes []interface{}
	if v, ok := s.Values[key]; ok {
		flashes = v.([]interface{})
	}
	s.Values[key] = append(flashes, value)

	return flashes
}

func (s *Session) Flashes(vars ...string) []interface{} {
	var flashes []interface{}
	key := flashesKey
	if len(vars) > 0 {
		key = vars[0]
	}
	if v, ok := s.Values[key]; ok {
		// Drop the flashes and return it.
		delete(s.Values, key)
		flashes = v.([]interface{})
	}
	return flashes
}

func (s *Session) Name() string {
	return s.name
}

func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	return s.store.Save(r, w, s)
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