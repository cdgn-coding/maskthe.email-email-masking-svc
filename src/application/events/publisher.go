package events

type Publisher interface {
	Dispatch(string) error
}
