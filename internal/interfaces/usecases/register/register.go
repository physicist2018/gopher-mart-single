package register

type UseCase interface {
	Execute(login, password string) (string, error)
}
