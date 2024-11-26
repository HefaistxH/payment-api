package delivery

import (
	"database/sql"
	"fmt"
	"mnc-techtest/config"
	"mnc-techtest/delivery/controller"
	"mnc-techtest/repository"
	"mnc-techtest/shared/service"
	"mnc-techtest/usecase"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Server struct {
	jwtService service.JwtService
	authUc     usecase.AuthUsecase
	customerUc usecase.CustomerUsecase
	engine     *gin.Engine
	host       string
}

func (s *Server) initRoute() {
	rg := s.engine.Group(config.ApiGroup)
	controller.NewAuthController(s.authUc, rg).Route()
	controller.NewCustomerController(s.customerUc, rg).Route()

}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	cfg, _ := config.NewConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open(cfg.DbConfig.Driver, dsn)
	if err != nil {
		panic("connection error")
	}
	jwtService := service.NewJwtService(cfg.JwtConfig)

	customerRepo := repository.NewCustomerRepository(db)
	customerUc := usecase.NewCustomerUsecase(customerRepo, jwtService)
	credRepo := repository.NewCredentialRepository(db)
	authUC := usecase.NewAuthUsecase(customerRepo, credRepo, jwtService)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("failed to create file")
	}
	logrus.SetOutput(file)

	return &Server{
		authUc:     authUC,
		customerUc: customerUc,
		jwtService: jwtService,
		engine:     engine,
		host:       host,
	}
}
