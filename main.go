package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hetfdex/tiny-bank/internal/domain"
	"github.com/hetfdex/tiny-bank/internal/handler"
	"github.com/hetfdex/tiny-bank/internal/repository/accountrepo"
	"github.com/hetfdex/tiny-bank/internal/repository/historyrepo"
	"github.com/hetfdex/tiny-bank/internal/repository/userrepo"
	"github.com/hetfdex/tiny-bank/internal/service"
)

func main() {
	userRepo, accountRepo, historyRepo := configRepo()

	svc := configSvc(userRepo, accountRepo, historyRepo)

	router := getRouter()

	configHandlers(router, svc)

	startServer(router)
}

func configRepo() (userrepo.Repo, accountrepo.Repo, historyrepo.Repo) {
	return userrepo.New(make(map[string]domain.User)),
		accountrepo.New(make(map[string]domain.Account)),
		historyrepo.New(make(map[string]domain.History))
}

func configSvc(
	userRepo userrepo.Repo,
	accountRepo accountrepo.Repo,
	historyRepo historyrepo.Repo,
) service.Service {
	return service.New(userRepo, accountRepo, historyRepo)
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
