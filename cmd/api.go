package main

import (
	"log"
	"nekeroo/myowntrip/internal/moneyrates"
	"nekeroo/myowntrip/internal/poi/clients"
	"nekeroo/myowntrip/internal/poi/httppoi"
	"nekeroo/myowntrip/internal/poi/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server started at %s", app.config.addr)

	return srv.ListenAndServe()
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	moneyRatesService := moneyrates.NewService()
	moneyRatesHandler := moneyrates.NewMoneyRatesHandler(*moneyRatesService)

	userAgent := "my-go-api/1.0 (contact: dev@example.com)"

	nominatimClient := clients.NewNominatimClient(userAgent)
	overpassClient := clients.NewOverpassClient(userAgent)

	poiService := service.New(nominatimClient, overpassClient)
	poiHandler := httppoi.NewHandler(poiService)

	r.Get("/rates", moneyRatesHandler.RetrieveLatestMoneyRates)

	r.Post("/pois", poiHandler.SearchPOIs)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API OK"))
	})

	return r
}

type application struct {
	config config
}

type config struct {
	addr string
}
