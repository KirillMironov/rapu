package mocks

type TokenManagerMock struct{}

func (TokenManagerMock) GenerateAuthToken(userId string) (string, error) {
	return "token", nil
}

func (TokenManagerMock) VerifyAuthToken(token string) (string, error) {
	return "token", nil
}
