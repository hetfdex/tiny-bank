package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/handler"
	"github.com/hetfdex/tiny-bank/internal/repository/accountrepo"
	"github.com/hetfdex/tiny-bank/internal/repository/userrepo"
	"github.com/hetfdex/tiny-bank/internal/service"
)

func main() {
	userRepo, accountRepo := configRepo()

	svc := configSvc(userRepo, accountRepo)

	router := getRouter()

	configHandlers(router, svc)

	startServer(router)
}

func configRepo() (userrepo.Repo, accountrepo.Repo) {
	return userrepo.New(make(map[string]domain.User)),
		accountrepo.New(make(map[string]domain.Account))
}

func configSvc(
	userRepo userrepo.Repo,
	accountRepo accountrepo.Repo,
) service.Service {
	return service.New(userRepo, accountRepo)
}

func getRouter() *gin.Engine {
	return gin.Default()
}

func configHandlers(router *gin.Engine, svc service.Service) {
	hdl := handler.New(svc)

	hdl.ConfigHandlers(router)
}

func startServer(router *gin.Engine) {
	router.Run(":8080")
}
