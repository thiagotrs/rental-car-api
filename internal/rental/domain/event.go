package domain

type Event interface {
	Name() string
}

type OpenedOrder struct {
	ID        string `json:"id"`
	CarId     string `json:"carId"`
	StationId string `json:"stationId"`
}

func (c OpenedOrder) Name() string {
	return "order.opened"
}

type ConfirmedOrder struct {
	ID    string `json:"id"`
	CarId string `json:"carId"`
}

func (c ConfirmedOrder) Name() string {
	return "order.confirmed"
}

type ClosedOrder struct {
	ID        string `json:"id"`
	CarId     string `json:"carId"`
	StationId string `json:"stationId"`
	FinalKM   uint64 `json:"finalKM"`
}

func (c ClosedOrder) Name() string {
	return "order.closed"
}

type CanceledOrder struct {
	ID        string `json:"id"`
	CarId     string `json:"carId"`
	StationId string `json:"stationId"`
	FinalKM   uint64 `json:"finalKM"`
}

func (c CanceledOrder) Name() string {
	return "order.canceled"
}
