package session

type Session struct {
	ID string

	Values map[interface{}]interface{}

	Options *Options

	IsNew bool

	store Store

	name string
}

func NewSession(store Store, name string) *Session {
	return &Session{
		Values:  make(map[interface{}]interface{}),
		store:   store,
		name:    name,
		Options: new(Options),
	}
}
