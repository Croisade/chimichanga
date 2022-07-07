package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/croisade/chimichanga/pkg/models"
	"github.com/croisade/chimichanga/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RunController struct {
	RunService services.RunService
	JWTService services.JWTAuthService
}

func NewRunController(runService services.RunService, jwtService services.JWTAuthService) RunController {
	return RunController{
		RunService: runService,
		JWTService: jwtService,
	}
}

func (rc *RunController) getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	}
	return "Unknown error"
}

func (rc *RunController) handleValidationError(ctx *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {

		out := make([]ErrorValidationMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorValidationMsg{fe.Field(), rc.getErrorMsg(fe)}
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": out})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
	return
}

func (rc *RunController) CreateRun(ctx *gin.Context) {
	var run models.Run
	if err := ctx.ShouldBindJSON(&run); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	result, err := rc.RunService.CreateRun(&run)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
	return
}

func (rc *RunController) GetRun(ctx *gin.Context) {
	var run *services.RunRequest

	if err := ctx.ShouldBindJSON(&run); err != nil {
		rc.handleValidationError(ctx, err)
		return
	}
	returnedRun, err := rc.RunService.GetRun(run)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, returnedRun)
	return
}

func (rc *RunController) GetAll(ctx *gin.Context) {
	// var runAccountId *services.RunFetchRequest
	runAccountId := ctx.Param("accountId")
	fmt.Println(runAccountId)
	fmt.Println(ctx.Request.Body)
	if err := ctx.ShouldBind(&runAccountId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	runs, err := rc.RunService.GetAll(runAccountId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, runs)
	return
}

func (rc *RunController) UpdateRun(ctx *gin.Context) {
	var run services.RunUpdateRequest
	if err := ctx.ShouldBindJSON(&run); err != nil {
		rc.handleValidationError(ctx, err)
		return
	}

	updatedRun, err := rc.RunService.UpdateRun(&run)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedRun)
	return
}

func (rc *RunController) DeleteRun(ctx *gin.Context) {
	var run *services.RunRequest

	if err := ctx.ShouldBindJSON(&run); err != nil {
		rc.handleValidationError(ctx, err)
		return
	}

	err := rc.RunService.DeleteRun(run)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	return
}

func (rc *RunController) RegisterRunRoutes(rg *gin.RouterGroup) {
	runRoute := rg.Group("/run")
	// runRoute.Use(gindump.Dump())
	// runRoute.Use(gindump.DumpWithOptions(true, true, true, true, true, func(dumpStr string) {
	// 	fmt.Println(dumpStr)
	// }))
	runRoute.POST("/create", rc.CreateRun)
	runRoute.GET("", rc.GetRun)
	runRoute.GET("/fetch/:accountId", rc.GetAll)
	runRoute.DELETE("/delete", rc.DeleteRun)
	runRoute.PUT("/update", rc.UpdateRun)
}
