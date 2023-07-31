package order

type OrderServicer interface {
	GetById(id string) (*Order, error)
	GetIdsList() []string
	Insert(Order) error
}
