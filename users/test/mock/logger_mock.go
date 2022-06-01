package mock

type LoggerMock struct{}

func (LoggerMock) Info(...interface{}) {}

func (LoggerMock) Error(...interface{}) {}
