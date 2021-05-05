package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-template/settings"
	"net/http"
)

func Homepage(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("It works. version: %v", settings.Conf.Version))
}
