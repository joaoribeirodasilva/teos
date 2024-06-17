package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppAppsList(c *gin.Context) {

	//vars := controllers.GetValues(c)

	/* 	model := models.AppAppModel{}
	   	models := []models.AppAppModel{}

	   	coll := models.NewBaseCollection(
	   		model.GetCollectionName(),
	   		vars.Db,
	   		nil,
	   	)

	   	logger.Info("will insert test on the database")

	   	userId := primitive.NewObjectID()
	   	model.SetUserID(&vars.UserID)

	   	count, err := coll.FindAll(vars.Query.filter, &models)
	   	if err != nil {
	   		if err != models.ErrNotFound {
	   			httpErr := logger.Error(logger.LogStatusNotFound, nil, "", err, nil)
	   			c.AbortWithStatusJSON(httpErr.Status, httpErr)
	   			return
	   		}
	   		httpErr := logger.Error(logger.LogStatusNotFound, nil, "", err, nil)
	   		c.AbortWithStatusJSON(httpErr.Status, httpErr)
	   		return
	   	}

	   	m.Count = count
	   	m.Rows = docs */

	//c.JSON(http.StatusOK, &m)
	c.Status(http.StatusOK)
}

func AppAppsGet(c *gin.Context) {

	/*
		 	doc := models.NewAppAppModel(c)

			if appErr := doc.FindByQueryID(doc); appErr != nil {
				c.AbortWithStatusJSON(appErr.HttpCode, appErr)
				return
			}

			c.JSON(http.StatusOK, &doc)
	*/
	c.Status(http.StatusOK)
}

func AppAppsCreate(c *gin.Context) {

	/* 	doc := models.NewAppAppModel(c)

	   	if appErr := doc.Bind(); appErr != nil {
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	}

	   	if appErr := doc.Validate(); appErr != nil {
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	}

	   	doc.FillMeta(true, false)

	   	if appErr := doc.Create(bson.D{{Key: "appKey", Value: doc.AppKey}}, doc, "appKey"); appErr != nil {
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	}

	   	response := responses.ResponseCreated{
	   		ID: doc.ID,
	   	}
	   	c.JSON(http.StatusCreated, response) */
	c.Status(http.StatusOK)
}

func AppAppsUpdate(c *gin.Context) {

	/* 	doc := models.NewAppAppModel(c)

	   	if appErr := doc.Bind(); appErr != nil {
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	}

	   	if appErr := doc.Validate(); appErr != nil {
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	}

	   	id := doc.GetValues().Query.ID
	   	if id == nil {
	   		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::AppApps::Update:", "", "invalid id")
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	}

	   	doc.ID = *id

	   	exists := models.NewAppAppModel(c)
	   	if appErr := exists.First(bson.D{{Key: "appKey", Value: doc.AppKey}}, exists); appErr != nil {
	   		if appErr.HttpCode != http.StatusNotFound {
	   			c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   			return
	   		}
	   	}

	   	if exists.ID != doc.ID {
	   		appErr := service_log.Error(0, http.StatusConflict, "[CONTROLLERS]::AppApp::Update", "appKey", "document already exists with id: %s", exists.ID.Hex())
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	}

	   	exists.Name = doc.Name
	   	exists.Description = doc.Description
	   	exists.AppKey = doc.AppKey
	   	exists.Active = doc.Active

	   	doc.FillMeta(false, false)

	   	if appErr := doc.Update(bson.D{{Key: "_id", Value: exists.ID}}, exists, "appKey"); appErr != nil {
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	} */

	c.Status(http.StatusOK)
}

func AppAppsDelete(c *gin.Context) {

	/* 	doc := models.NewAppAppModel(c)

	   	id := doc.GetValues().Query.ID
	   	if id == nil {
	   		appErr := service_log.Error(0, http.StatusBadRequest, "CONTROLLER::AppApps::Update:", "", "invalid id")
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	}

	   	if appErr := doc.Delete(*id); appErr != nil {
	   		c.AbortWithStatusJSON(appErr.HttpCode, appErr)
	   		return
	   	} */

	c.Status(http.StatusOK)
}
