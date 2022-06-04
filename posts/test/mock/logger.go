package mock

type Logger struct{}

func (Logger) Error(...interface{}) {}
