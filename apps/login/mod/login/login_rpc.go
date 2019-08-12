package login

type LoginServer struct {
}

func (l *LoginServer) Login(user string, pass string) (bool, error) {
	if pass == "123" {
		return true, nil
	}
	return false, nil
}
