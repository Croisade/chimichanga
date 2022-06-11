package controllers

import (
	"net/http"

	"github.com/croisade/chimichanga/pkg/models"
	"github.com/croisade/chimichanga/pkg/services"
	"github.com/gin-gonic/gin"
)

type RunController struct {
	RunService services.RunService
}

func NewRunController(runService services.RunService) RunController {
	return RunController{
		RunService: runService,
	}
}

func (rc *RunController) CreateRun(ctx *gin.Context) {
	var run models.Run
	if err := ctx.ShouldBindJSON(&run); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	_, err := rc.RunService.CreateRun(&run)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (rc *RunController) GetRun(ctx *gin.Context) {
	var run *services.RunRequest

	if err := ctx.ShouldBindJSON(&run); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	usr, err := rc.RunService.GetRun(run)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, usr)
	return
}

func (rc *RunController) GetAll(ctx *gin.Context) {
	var runAccountId *services.RunFetchRequest

	if err := ctx.ShouldBindJSON(&runAccountId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	runs, err := rc.RunService.GetAll(runAccountId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, runs)
	return
}

func (rc *RunController) UpdateRun(ctx *gin.Context) {
	var run services.RunUpdateRequest
	if err := ctx.ShouldBindJSON(&run); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedRun, err := rc.RunService.UpdateRun(&run)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedRun)
	return
}

func (rc *RunController) DeleteRun(ctx *gin.Context) {
	var run *services.RunRequest

	if err := ctx.ShouldBindJSON(&run); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := rc.RunService.DeleteRun(run)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	return
}

func (rc *RunController) RegisterRunRoutes(rg *gin.RouterGroup) {
	userRoute := rg.Group("/run")
	userRoute.POST("/create", rc.CreateRun)
	userRoute.GET("/", rc.GetRun)
	userRoute.GET("/fetch", rc.GetAll)
	userRoute.DELETE("/delete", rc.DeleteRun)
	userRoute.PUT("/update", rc.UpdateRun)
}