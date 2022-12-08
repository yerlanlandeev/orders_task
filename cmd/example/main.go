package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var logger = log.Default()

var AvailableRooms = map[string]struct{}{"econom": {}, "standart": {}, "lux": {}}

type Order struct {
	Room      string
	UserEmail string
	From      string
	To        string
}

var ActualOrders = []Order{}

func makeOrder(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("email")
	if userEmail == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	room := r.URL.Query().Get("room")
	if room == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if _, isOK := AvailableRooms[room]; !isOK {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	from := r.URL.Query().Get("from")
	if from == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	toTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newOrder := Order{
		Room:      room,
		UserEmail: userEmail,
		From:      from,
		To:        to,
	}

	for _, order := range ActualOrders {
		currentOrderFromTime, _ := time.Parse("2006-01-02", order.From)
		currentOrderToTime, _ := time.Parse("2006-01-02", order.To)
		if !(currentOrderToTime.Before(fromTime) || currentOrderFromTime.After(toTime)) {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	ActualOrders = append(ActualOrders, newOrder)

	w.WriteHeader(http.StatusCreated)
	LogInfo("Method makeOrder was successfully done")
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("email")
	if userEmail == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	res := []Order{}
	for _, item := range ActualOrders {
		if item.UserEmail == userEmail {
			res = append(res, item)
		}
	}

	b, err := json.Marshal(res)
	if err != nil {
		LogErrorf("error in getOrders method: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)

	LogInfo("Method getOrders was successfully done")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/order", makeOrder)
	mux.HandleFunc("/orders", getOrders)

	err := http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		LogInfo("server closed")
	} else if err != nil {
		LogErrorf("error listening for server: %s", err)
		os.Exit(1)
	}
}

func LogErrorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v)
	logger.Printf("[Error]: %s\n", msg)
}

func LogInfo(format string, v ...any) {
	msg := fmt.Sprintf(format, v)
	logger.Printf("[Info]: %s\n", msg)
}
