package mocks

type LoggerMock struct{}

func (LoggerMock) Info(args ...interface{}) {}

func (LoggerMock) Infof(template string, args ...interface{}) {}

func (LoggerMock) Error(args ...interface{}) {}

func (LoggerMock) Errorf(template string, args ...interface{}) {}
