package order

type Service interface {
	CheckOrdersForAccrual() error
}
