package mock

type JWTManager struct{}

func (JWTManager) Generate(string) (string, error) {
	return "token", nil
}

func (JWTManager) Verify(string) (string, error) {
	return "token", nil
}
