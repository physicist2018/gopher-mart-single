package auth

type UseCase interface {
	Execute(login, password string) (string, error)
}
