package mocks

type TokenManagerMock struct{}

func (TokenManagerMock) Generate(userId string) (string, error) {
	return "token", nil
}

func (TokenManagerMock) Verify(token string) (string, error) {
	return "token", nil
}
