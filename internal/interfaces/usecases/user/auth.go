package user

type AuthUseCase interface {
	Execute(login, password string) (string, error)
}
