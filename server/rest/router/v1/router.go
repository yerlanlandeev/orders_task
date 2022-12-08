package v1

import (
	"applicationDesignTest/service/orderserve"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type RoutingServer struct {
	orderService orderserve.OrderProcessor
	logger       *zap.Logger
}

func NewRoutingServer(logger *zap.Logger, orderService orderserve.OrderProcessor) *RoutingServer {
	return &RoutingServer{
		orderService: orderService,
		logger:       logger,
	}
}

func (router *RoutingServer) ListenAndServe(address string) error {

	mux := http.NewServeMux()

	mux.HandleFunc("/order", router.Order)
	mux.HandleFunc("/orders", router.GetAll)

	return http.ListenAndServe(address, mux)
}

type Order struct {
	Room      string `json:"room"`
	UserEmail string `json:"user_email"`
	From      string `json:"from"`
	To        string `json:"to"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func (router *RoutingServer) Order(w http.ResponseWriter, r *http.Request) {
	var args Order
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		// log
		router.logger.Error("failed to decode the message", zap.Error(err))
		writeMessage(w, http.StatusOK, ErrorMessage{
			Error: fmt.Sprintf("failed to decode the message. err: %v", err),
		})
		return
	}
	defer r.Body.Close()

	// log
	router.logger.Info("creating a new order", zap.String("room", args.Room))

	fromThisTime, err := time.Parse(`2006-01-02 15:04`, args.From)
	if err != nil {
		// log
		router.logger.Error("failed to parse time (from)", zap.Error(err))
		writeMessage(w, http.StatusOK, ErrorMessage{
			Error: fmt.Sprintf("failed to parse time (from). Time format must be 2006-01-02 15:04. err: %v", err),
		})
		return
	}

	toThisTime, err := time.Parse(`2006-01-02 15:04`, args.From)
	if err != nil {
		// log
		router.logger.Error("failed to parse time (to)", zap.Error(err))
		writeMessage(w, http.StatusOK, ErrorMessage{
			Error: fmt.Sprintf("failed to parse time (to). Time format must be 2006-01-02 15:04. err: %v", err),
		})
		return
	}

	responseOrder, err := router.orderService.Order(context.Background(), orderserve.Order{
		Room:      args.Room,
		From:      fromThisTime,
		To:        toThisTime,
		UserEmail: args.UserEmail,
	})
	if err != nil {
		// log
		router.logger.Error("failed to order the room", zap.Error(err))
		writeMessage(w, http.StatusOK, ErrorMessage{
			Error: fmt.Sprintf("failed to order the room. err: %v", err),
		})
		return
	}

	router.logger.Info("successfully created an order")
	router.logger.Debug("successfully created an order", zap.Any("response", responseOrder))

	writeMessage(w, http.StatusCreated, responseOrder)
}

func (router *RoutingServer) GetAll(w http.ResponseWriter, r *http.Request) {
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	var offset, limit int64
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
	}

	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 1000
	}

	router.logger.Info("fetching orders",
		zap.String("offsetStr", offsetStr),
		zap.String("limitStr", limitStr),
		zap.Int64("offset", offset),
		zap.Int64("limit", limit),
	)

	orders, err := router.orderService.GetAll(context.Background(), orderserve.ParamsGetAll{
		Offset: int(offset),
		Limit:  int(limit),
	})
	if err != nil {
		router.logger.Error("failed to fetch all orders", zap.Error(err))
		return
	}

	router.logger.Info("successfully fetched all orders")
	router.logger.Debug("fetched all orders", zap.Int("response-length", len(orders)))

	writeMessage(w, http.StatusOK, orders)
}

func writeMessage(w http.ResponseWriter, statusCode int, responseMessage interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(responseMessage); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
