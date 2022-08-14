package auth

type Service interface {
	Login(username, password string) string
	Register(username, password string)
	GetUserId(sessionKey string) (int, error)
}

type service struct {
	users map[string]int
}
func (s service) Login(username string, password string) string {
	panic("not implemented") // TODO: Implement
}

func (s service) Register(username string, password string) {
	panic("not implemented") // TODO: Implement
}

func (s service) GetUserId(sessionKey string) (int, error) {
	panic("not implemented") // TODO: Implement
}

// s service



