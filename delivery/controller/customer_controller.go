package controller

import (
	"mnc-techtest/config"
	"mnc-techtest/entity/dto"
	"mnc-techtest/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CustomerController struct {
	customerUsecase usecase.CustomerUsecase
	rg              *gin.RouterGroup
}

func NewCustomerController(customerUsecase usecase.CustomerUsecase, rg *gin.RouterGroup) *CustomerController {
	return &CustomerController{customerUsecase: customerUsecase, rg: rg}
}

func (c *CustomerController) Payment(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing Authorization header"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Authorization header format"})
		return
	}
	logrus.Infof("Token: %v", token)

	var paymentRequest dto.Payment
	if err := ctx.ShouldBindJSON(&paymentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	paymentResponse, err := c.customerUsecase.Payment(token, paymentRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, paymentResponse)
}

func (c *CustomerController) Route() {
	c.rg.POST(config.PaymentRoute, c.Payment)
}
