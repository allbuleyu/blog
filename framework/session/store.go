package session

type Store interface {
	New() *Session

	Get() *Session

	Save() error
}
