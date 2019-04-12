package session

import (
	"context"
	"encoding/gob"
	"fmt"
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

// NewSession is called by session stores to create a new session instance.
func NewSession(store Store, name string) *Session {
	return &Session{
		Values:  make(map[interface{}]interface{}),
		store:   store,
		name:    name,
		Options: new(Options),
	}
}

// Default flashes key.
const flashesKey = "_flash"


// Flashes returns a slice of flash messages from the session.
//
// A single variadic argument is accepted, and it is optional: it defines
// the flash key. If not defined "_flash" is used by default.
func (s *Session) Flashes(vars ...string) []interface{} {
	key := flashesKey
	if len(vars) > 0 {
		key = vars[0]
	}

	var flashes []interface{}
	if v, ok := s.Values[key]; ok {

		delete(s.Values, key)
		flashes = v.([]interface{})
	}

	return flashes
}

// AddFlash adds a flash message to the session.
//
// A single variadic argument is accepted, and it is optional: it defines
// the flash key. If not defined "_flash" is used by default.
func (s *Session) AddFlash(value interface{}, vars ...string) []interface{} {
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

// Save is a convenience method to save this session. It is the same as calling
// store.Save(request, response, session). You should call Save before writing to
// the response or returning from the handler.
func (s *Session) Save(r *http.Request, w http.ResponseWriter) error {
	return s.store.Save(r, w, s)
}

// Name returns the name used to register the session.
func (s *Session) Name() string {
	return s.name
}


// Store returns the session store used to register the session.
func (s *Session) Store() Store {
	return s.store
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

// GetRegistry 在当前请求r返回一个*Registry实例
func GetRegistry(r *http.Request) *Registry {
	var ctx = r.Context()
	registry := ctx.Value(registryKey)
	if registry != nil {
		return registry.(*Registry)		// miracle???
	}

	newRegistry := &Registry{
		request:r,
		sessions:make(map[string]sessionInfo),
	}

	*r = *r.WithContext(context.WithValue(ctx, registryKey, newRegistry))

	return newRegistry
}

// Get  注册以及返回一个指定name与store的session,
//
// 如果指定名字的session不存在,则按照指定的name创建一个
func (s *Registry) Get(store Store, name string) (session *Session, err error) {
	if !isCookieNameValid(name) {
		return nil, fmt.Errorf("sessions: invalid character in cookie name: %s", name)
	}

	if info, ok := s.sessions[name]; ok {
		session , err = info.s, info.e
	}else {
		session, err = store.New(s.request, name)
		session.name = name
		s.sessions[name] = sessionInfo{s:session, e: nil}
	}
	session.store = store

	return
}

// Save saves all sessions registered for the current request.
// 保存当前请求中所有注册的session
func (s *Registry) Save(w http.ResponseWriter) error {
	var errMulti MultiError
	for name, info := range s.sessions {
		session := info.s
		if session.store == nil {
			errMulti = append(errMulti, fmt.Errorf("sessions: missing store for session %q", name))
		}else if err := session.store.Save(s.request, w, session); err != nil {
			errMulti = append(errMulti, fmt.Errorf(
				"sessions: error saving session %q -- %v", name, err))
		}
	}

	return errMulti
}


// Helpers --------------------------------------------------------------------

func Init() {
	gob.Register([]interface{}{})
}

// Save saves all sessions used during the current request.
func Save(r *http.Request, w http.ResponseWriter) error {
	return GetRegistry(r).Save(w)
}

// NewCookie returns an http.Cookie with the options set. It also sets
// the Expires field calculated based on the MaxAge value, for Internet
// Explorer compatibility.
func NewCookie(name, value string, options *Options) *http.Cookie {
	cookie := newCookieFromOptions(name, value, options)
	if cookie.MaxAge > 0 {
		d := time.Duration(cookie.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	}else if cookie.MaxAge < 0 {
		// Set it to the past to expire now.
		cookie.Expires = time.Unix(1, 0)
	}

	return cookie
}


// Error ----------------------------------------------------------------------

type MultiError []error

func (m MultiError) Error() string {
	s, n := "", 0
	for _, e := range m {
		if e != nil {
			if n == 0 {
				s = e.Error()
			}
			n++
		}
	}
	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, n-1)
}