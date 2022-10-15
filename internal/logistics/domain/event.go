package domain

type Event interface {
	Name() string
}

type CarAdded struct {
	ID        string `json:"id"`
	StationId string `json:"stationId"`
}

func (c CarAdded) Name() string {
	return "car.added"
}

type CarUnderMaintenance struct {
	ID        string    `json:"id"`
	StationId string    `json:"stationId"`
	CarStatus CarStatus `json:"carStatus"`
}

func (c CarUnderMaintenance) Name() string {
	return "car.under-maintenance"
}

type CarInTransfer struct {
	ID            string `json:"id"`
	StationIdFrom string `json:"stationIdFrom"`
	StationIdTo   string `json:"stationIdTo"`
}

func (c CarInTransfer) Name() string {
	return "car.in-transfer"
}

type CarParked struct {
	ID        string `json:"id"`
	StationId string `json:"stationId"`
	KM        uint64 `json:"km"`
}

func (c CarParked) Name() string {
	return "car.parked"
}

type SyncCarParked struct {
	ID        string `json:"id"`
	StationId string `json:"stationId"`
	KM        uint64 `json:"km"`
}

func (c SyncCarParked) Name() string {
	return "sync.car.parked"
}

type SyncCarReserved struct {
	ID        string `json:"id"`
	StationId string `json:"stationId"`
}

func (c SyncCarReserved) Name() string {
	return "sync.car.reserved"
}

type SyncCarInTransit struct {
	ID string `json:"id"`
}

func (c SyncCarInTransit) Name() string {
	return "sync.car.in-transit"
}
