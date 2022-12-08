package main

import (
	"applicationDesignTest/repository/reporder"
	"applicationDesignTest/repository/reporoom"
	v1 "applicationDesignTest/server/rest/router/v1"
	"applicationDesignTest/service/orderserve"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	orderRepo := reporder.NewRoomOrder()
	roomRepo := reporoom.NewRoom()

	orderService := orderserve.NewOrder(orderRepo, roomRepo)

	var router = v1.NewRoutingServer(
		logger.With(zap.String("version", "v1"), zap.String("type", "router")),
		orderService,
	)

	address := ":8800"
	logger.Info("started serving http requests", zap.String("address", address))

	go func() {
		err := router.ListenAndServe(address)
		fmt.Printf("listening the port error, err: %v\n", err)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
