package main

import (
	"api_gateway/api/handler"
	"api_gateway/config"

	pb "api_gateway/protos"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	conn, err := grpc.Dial(cfg.CalculatorGRPCServiceHost+cfg.CalculatorGRPCServicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	calculatorService := pb.NewCalculatorServicesClient(conn)

	h := handler.New(cfg, calculatorService)

	r := gin.New()

	r.POST("/calculate", h.Addition)

	if err = r.Run(cfg.HTTPPort); err != nil {
		panic(err)
	}

	defer conn.Close()
}
