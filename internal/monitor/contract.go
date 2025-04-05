package monitor

type Store interface {
	SaveJSON(v interface{}) error
	LoadJSON(v interface{}) error
}
