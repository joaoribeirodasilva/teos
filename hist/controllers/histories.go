package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/teos/common/controllers"
	"github.com/joaoribeirodasilva/teos/common/responses"
	"github.com/joaoribeirodasilva/teos/hist/requests"
)

func HistoriesList(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func HistoriesGet(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func HistoriesCreate(c *gin.Context) {

	_, err := controllers.MustGetAll(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	request := requests.HistoriesCreate{}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	request.AuthenticationKey = strings.TrimSpace(request.AuthenticationKey)
	if request.AuthenticationKey == "" {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "authenticationKey",
				Message: "authenticationKey filed is required",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	//TODO: check authentication key

	request.Table = strings.TrimSpace(request.Table)
	if request.Table == "" {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "table",
				Message: "table filed is required",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if request.OriginalID == 0 {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "originalId",
				Message: "originalId filed is required",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	request.Data = strings.TrimSpace(request.Data)
	if request.Data == "" {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "data",
				Message: "data filed is required",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	request.Data = strings.TrimSpace(request.Data)
	if request.Data == "" {
		response := responses.ResponseErrorField{
			Error: responses.ErrField{
				Field:   "data",
				Message: "data filed is required",
			},
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	c.Status(http.StatusOK)
}
