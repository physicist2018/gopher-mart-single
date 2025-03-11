package user

type RegisterUseCase interface {
	Execute(login, password string) (string, error)
}
