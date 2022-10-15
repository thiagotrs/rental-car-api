package ipc

type CarData struct {
	ID        string `json:"id"`
	Age       uint16 `json:"age"`
	Plate     string `json:"plate"`
	Document  string `json:"document"`
	Model     string `json:"model"`
	Make      string `json:"make"`
	StationId string `json:"stationId"`
	KM        uint64 `json:"km"`
	Status    uint   `json:"status"`
}

type PolicyData struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Price   float32 `json:"price"`
	Unit    uint    `json:"unit"`
	MinUnit uint    `json:"minUnit"`
}

type LogisticsIPC interface {
	GetCar(stationId, carModel string) (*CarData, error)
}

type PricingIPC interface {
	GetPolicy(categoryId, carModel, policyId string) (*PolicyData, error)
}
