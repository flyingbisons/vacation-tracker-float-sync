package integrator

type RequestRepository interface {
	GetRequest(vtRequestID string) (Request, error)
	CreateRequest(vtRequestID string, floatTimeOffID int64) error
}
