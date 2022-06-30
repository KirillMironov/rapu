package mock

type Logger struct{}

func (Logger) Error(...interface{}) {}

func (Logger) Errorf(string, ...interface{}) {}
