package delivery

import (
	// "database/sql"
	"fmt"
	"mnc-techtest/config"
	"mnc-techtest/shared/service"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Server struct {
	jwtService service.JwtService
	engine     *gin.Engine
	host       string
}

func (s *Server) initRoutes() {
	// rg := s.engine.Group(config.ApiGroup)
}

func NewServer() *Server {
	cfg, _ := config.NewConfig()
	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbConfig.Host, cfg.DbConfig.Port, cfg.DbConfig.User, cfg.DbConfig.Password, cfg.DbConfig.Name)

	// db, err := sql.Open(cfg.DbConfig.Driver, dsn)
	// if err != nil {
	// 	panic("connection error")
	// }
	jwtService := service.NewJwtService(cfg.JwtConfig)
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)

	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("failed to create file")
	}
	logrus.SetOutput(file)

	return &Server{
		jwtService: jwtService,
		engine:     engine,
		host:       host,
	}
}
