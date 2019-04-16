package session

import (
	"encoding/base32"
	"github.com/gorilla/securecookie"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Store interface {
	// Get should return a cached session.
	Get(r *http.Request, name string) (*Session, error)

	// New should create and return a new session.
	//
	// Note that New should never return a nil session, even in the case of
	// an error if using the Registry infrastructure to cache the session.
	New(r *http.Request, name string) (*Session, error)

	// Save should persist session to the underlying store implementation.
	Save(r *http.Request, w http.ResponseWriter, s *Session) error
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	Codecs  []securecookie.Codec
	Options *Options // default configuration
}

func NewCookieStore(keyPairs ...[]byte) *CookieStore {
	cs := &CookieStore{
		Codecs:securecookie.CodecsFromPairs(keyPairs...),
		Options:&Options{
			MaxAge:86400 * 30,
			Path:"/",
		},
	}

	cs.MaxAge(cs.Options.MaxAge)
	return cs
}

func (s *CookieStore) Get(r *http.Request, name string) (*Session, error) {
	
	return nil, nil
}

// New returns a session for the given name without adding it to the registry.
//
// The difference between New() and Get() is that calling New() twice will
// decode the session data twice, while Get() registers and reuses the same
func (s *CookieStore) New(r *http.Request, name string) (*Session, error) {
	session := NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true

	var err error
	if c, cookieErr := r.Cookie(name); cookieErr == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.Values, s.Codecs...)
		if err == nil {
			session.IsNew = false
		}
	}

	return session, err
}

// Save adds a single session to the response.
func (s *CookieStore) Save(r *http.Request, w http.ResponseWriter, session *Session) error {
	 encode, err := securecookie.EncodeMulti(session.Name(), session.Values, s.Codecs...)
	 if err != nil {
	 	return err
	 }

	 http.SetCookie(w, NewCookie(session.Name(), encode, session.Options))

	return nil
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (s *CookieStore) MaxAge(age int) {
	s.Options.MaxAge = age

	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}




// FilesystemStore ------------------------------------------------------------



// CookieStore stores sessions using secure cookies.
type FilesystemStore struct {
	Codecs  []securecookie.Codec
	Options *Options // default configuration
	path string
}

var fileMutex sync.RWMutex

func NewFilesystemStore(path string, keyPairs ...[]byte) *FilesystemStore {
	if path == "" {
		path = os.TempDir()
	}
	fs := &FilesystemStore{
		Codecs:securecookie.CodecsFromPairs(keyPairs...),
		Options:&Options{
			MaxAge:86400 * 30,
			Path:"/",
		},
		path: path,
	}

	fs.MaxAge(fs.Options.MaxAge)
	return fs
}


func (s *FilesystemStore) Get(r *http.Request, name string) (*Session, error) {
	panic("implement me")
}


func (s *FilesystemStore) New(r *http.Request, name string) (*Session, error) {
	session := NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true

	var err error
	if c, cookieErr := r.Cookie(name); cookieErr == nil {
		err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...)
		if err == nil {

			if err = s.load(session); err == nil {
				session.IsNew = false
			}

		}
	}

	return session, err
}


func (s *FilesystemStore) Save(r *http.Request, w http.ResponseWriter, session *Session) error {
	var err error
	if session.Options.MaxAge <= 0 {
		if err = s.eras(session); err != nil {
			return err
		}

		http.SetCookie(w, NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		// Because the ID is used in the filename, encode it to
		// use alphanumeric characters only.
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}

	if err = s.save(session); err != nil {
		return err
	}

	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID,
		s.Codecs...)
	if err != nil {
		return err
	}

	http.SetCookie(w, NewCookie(session.Name(), encoded, session.Options))
	return nil
}


// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (s *FilesystemStore) MaxAge(age int) {
	s.Options.MaxAge = age

	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

// load reads a file and decodes its content into session.Values.
func (s *FilesystemStore) load(session *Session) error {
	filename := path.Join(s.path, session.ID)

	fileMutex.Lock()
	defer fileMutex.Unlock()

	fdata, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = securecookie.DecodeMulti(session.Name(), string(fdata), &session.Values, s.Codecs...)
	if err != nil {
		return err
	}

	return nil
}

// load reads a file and decodes its content into session.Values.
func (s *FilesystemStore) save(session *Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		s.Codecs...)
	if err != nil {
		return err
	}
	filename := filepath.Join(s.path, "session_"+session.ID)
	fileMutex.Lock()
	defer fileMutex.Unlock()
	return ioutil.WriteFile(filename, []byte(encoded), 0600)
}

// delete session file
func (s *FilesystemStore) eras(session *Session) error {
	filename := path.Join(s.path, session.ID)
	fileMutex.Lock()
	defer fileMutex.Unlock()

	err := os.Remove(filename)

	return err
}