package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/broker"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/config"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/database"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/events"
	"github.com/thiagotrs/rentalcar-ddd/internal/pkg/ipc"

	"github.com/thiagotrs/rentalcar-ddd/internal/logistics/adapters/consumer"
	ehLogistics "github.com/thiagotrs/rentalcar-ddd/internal/logistics/adapters/eventhandler"
	hLogistics "github.com/thiagotrs/rentalcar-ddd/internal/logistics/adapters/http"
	ipcLogistics "github.com/thiagotrs/rentalcar-ddd/internal/logistics/adapters/ipc"
	repoLogistics "github.com/thiagotrs/rentalcar-ddd/internal/logistics/adapters/repository"
	appLogistics "github.com/thiagotrs/rentalcar-ddd/internal/logistics/application"
	domainLogistics "github.com/thiagotrs/rentalcar-ddd/internal/logistics/domain"

	hPricing "github.com/thiagotrs/rentalcar-ddd/internal/pricing/adapters/http"
	ipcPricing "github.com/thiagotrs/rentalcar-ddd/internal/pricing/adapters/ipc"
	repoPricing "github.com/thiagotrs/rentalcar-ddd/internal/pricing/adapters/repository"
	appPricing "github.com/thiagotrs/rentalcar-ddd/internal/pricing/application"

	ehRental "github.com/thiagotrs/rentalcar-ddd/internal/rental/adapters/eventhandler"
	hRental "github.com/thiagotrs/rentalcar-ddd/internal/rental/adapters/http"
	repoRental "github.com/thiagotrs/rentalcar-ddd/internal/rental/adapters/repository"
	svcRental "github.com/thiagotrs/rentalcar-ddd/internal/rental/adapters/service"
	appRental "github.com/thiagotrs/rentalcar-ddd/internal/rental/application"
	domainRental "github.com/thiagotrs/rentalcar-ddd/internal/rental/domain"
)

func setupLogistics(db *sqlx.DB, r *mux.Router, e events.Dispatcher, b broker.Subscriber) ipc.LogisticsIPC {
	// LOGISTICS STATION

	stationRepo := repoLogistics.NewStationRepositorySqlx(context.Background(), db)
	stationUC := appLogistics.NewStationUseCase(stationRepo)
	stationController := hLogistics.NewStationController(stationUC)

	ehStation := ehLogistics.NewStationEventHandler(stationUC)
	e.Register(events.EventHandlerFunc(ehStation.HandleCarAdded), domainLogistics.CarAdded{}.Name())
	e.Register(events.EventHandlerFunc(ehStation.HandleCarParked), domainLogistics.CarParked{}.Name())
	e.Register(events.EventHandlerFunc(ehStation.HandleCarUnderMaintenance), domainLogistics.CarUnderMaintenance{}.Name())
	e.Register(events.EventHandlerFunc(ehStation.HandleCarInTransfer), domainLogistics.CarInTransfer{}.Name())
	e.Register(events.EventHandlerFunc(ehStation.HandleSyncCarParked), domainLogistics.SyncCarParked{}.Name())
	e.Register(events.EventHandlerFunc(ehStation.HandleSyncCarReserved), domainLogistics.SyncCarReserved{}.Name())

	cons := consumer.NewOrderConsumer(e)
	chOpenedOrder := b.Subscribe(string(consumer.OrderOpened))
	chConfirmedOrder := b.Subscribe(string(consumer.OrderConfirmed))
	chCanceledOrder := b.Subscribe(string(consumer.OrderCanceled))
	chClosedOrder := b.Subscribe(string(consumer.OrderClosed))
	go broker.Consume(chOpenedOrder, broker.ConsumerFunc(cons.ConsumeOpenedOrder))
	go broker.Consume(chConfirmedOrder, broker.ConsumerFunc(cons.ConsumeConfirmedOrder))
	go broker.Consume(chCanceledOrder, broker.ConsumerFunc(cons.ConsumeCanceledOrder))
	go broker.Consume(chClosedOrder, broker.ConsumerFunc(cons.ConsumeClosedOrder))

	r.HandleFunc("/stations/{id}/capacity", stationController.UpdateStationCapacity).Methods("PUT")
	r.HandleFunc("/stations/{id}", stationController.GetStationById).Methods("GET")
	r.HandleFunc("/stations/{id}", stationController.DeleteStation).Methods("DELETE")
	r.HandleFunc("/stations/", stationController.GetStations).Methods("GET")
	r.HandleFunc("/stations/", stationController.CreateStation).Methods("POST")

	// LOGISTICS CAR

	carRepo := repoLogistics.NewCarRepositorySqlx(context.Background(), db, e)
	carUC := appLogistics.NewCarUseCase(carRepo, stationRepo)
	carController := hLogistics.NewCarController(carUC)

	carIPC := ipcLogistics.NewCarIPC(carUC)

	ehCar := ehLogistics.NewCarEventHandler(carUC)
	e.Register(events.EventHandlerFunc(ehCar.HandleSyncCarInTransit), domainLogistics.SyncCarInTransit{}.Name())
	e.Register(events.EventHandlerFunc(ehCar.HandleSyncCarParked), domainLogistics.SyncCarParked{}.Name())
	e.Register(events.EventHandlerFunc(ehCar.HandleSyncCarReserved), domainLogistics.SyncCarReserved{}.Name())

	r.HandleFunc("/cars/{id}/maintenance/", carController.UpdateCarToMaintenance).Methods("PUT")
	r.HandleFunc("/cars/{id}/park/", carController.UpdateCarToPark).Methods("PUT")
	r.HandleFunc("/cars/{id}/transfer/", carController.UpdateCarToTransfer).Methods("PUT")
	r.HandleFunc("/cars/{id}", carController.GetCarById).Methods("GET")
	r.HandleFunc("/cars/{id}", carController.DeleteCar).Methods("DELETE")
	r.HandleFunc("/cars/", carController.SearchCars).Methods("GET")
	r.HandleFunc("/cars/", carController.CreateCar).Methods("POST")

	return carIPC
}

func setupPricing(db *sqlx.DB, r *mux.Router, e events.Dispatcher, b broker.Publisher) ipc.PricingIPC {
	categoryRepo := repoPricing.NewCategoryRepositorySqlx(context.Background(), db)
	categoryUC := appPricing.NewCategoryUseCase(categoryRepo)
	categoryController := hPricing.NewCategoryController(categoryUC)

	categoryIPC := ipcPricing.NewCategoryIPC(categoryUC)

	r.HandleFunc("/categories/{id}/model/", categoryController.UpdateAddModelInCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}/model/", categoryController.UpdateDelModelInCategory).Methods("DELETE")
	r.HandleFunc("/categories/{id}/policy/", categoryController.UpdateAddPolicyInCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}/policy/{policyId}", categoryController.UpdateDelPolicyInCategory).Methods("DELETE")
	r.HandleFunc("/categories/{id}", categoryController.GetCategoryById).Methods("GET")
	r.HandleFunc("/categories/{id}", categoryController.DeleteCategory).Methods("DELETE")
	r.HandleFunc("/categories/", categoryController.GetCategories).Methods("GET")
	r.HandleFunc("/categories/", categoryController.CreateCategory).Methods("POST")

	return categoryIPC
}

func setupRental(db *sqlx.DB, r *mux.Router, e events.Dispatcher, b broker.Publisher, l ipc.LogisticsIPC, p ipc.PricingIPC) {
	orderSvc := svcRental.NewOrderServiceIPC(l, p)
	orderRepo := repoRental.NewOrderRepositorySqlx(context.Background(), db, e)
	orderUC := appRental.NewOrderUseCase(orderRepo, orderSvc)
	orderController := hRental.NewOrderController(orderUC)

	ehOrder := ehRental.NewOrderEventHandler(b)
	e.Register(events.EventHandlerFunc(ehOrder.HandleOpenedOrder), domainRental.OpenedOrder{}.Name())
	e.Register(events.EventHandlerFunc(ehOrder.HandleConfirmedOrder), domainRental.ConfirmedOrder{}.Name())
	e.Register(events.EventHandlerFunc(ehOrder.HandleClosedOrder), domainRental.ClosedOrder{}.Name())
	e.Register(events.EventHandlerFunc(ehOrder.HandleCanceledOrder), domainRental.CanceledOrder{}.Name())

	r.HandleFunc("/orders/{id}/confirm/", orderController.UpdateToComfirmOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}/close/", orderController.UpdateToCloseOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}/cancel/", orderController.UpdateToCancelOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}", orderController.GetOrderById).Methods("GET")
	r.HandleFunc("/orders/", orderController.CreateOrder).Methods("POST")
}

func runAPI(r *mux.Router, c config.AppConfig) {
	fmt.Printf("API is running on port %d", c.Server.Port)
	addr := fmt.Sprintf("%v:%d", c.Server.Host, c.Server.Port)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	config := config.GetConfig()

	db := database.GetDBConn(config.Database)
	defer db.Close()

	pubsub := broker.NewPubSub()
	defer pubsub.Close()

	dispatcher := events.NewEventDispatcher()

	router := mux.NewRouter()

	logisticsIPC := setupLogistics(db, router, dispatcher, pubsub)
	pricingIPC := setupPricing(db, router, dispatcher, pubsub)
	setupRental(db, router, dispatcher, pubsub, logisticsIPC, pricingIPC)

	// API

	router.HandleFunc("/status/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		fmt.Fprint(w, `{"name":"API", "version":"v1"}`)
	})

	runAPI(router, config)
}
