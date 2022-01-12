package mocks

type TokenManagerMock struct{}

func (TokenManagerMock) Generate(string) (string, error) {
	return "token", nil
}

func (TokenManagerMock) Verify(string) (string, error) {
	return "token", nil
}
