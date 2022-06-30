package mock

type JWT struct{}

func (JWT) Generate(string) (string, error) {
	return "token", nil
}

func (JWT) Verify(string) (string, error) {
	return "token", nil
}
