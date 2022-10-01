package server

import (
	"github.com/Jeadie/mailhub/pkg/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ConstructEndpoints on top of a gin Engine and a dao object.
func ConstructEndpoints(r *gin.Engine, dao db.Dao) {

	r.GET("/mailhub", func(c *gin.Context) {
		smss, err := dao.GetAllSmss()
		if err != nil {
			c.JSON(http.StatusTeapot, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"count": len(smss),
				"smss":  smss,
			})
		}
	})

	r.GET("/sms/:name", func(c *gin.Context) {
		name := c.Param("name")
		smss, err := dao.GetSmssTo(name)

		if err != nil {
			c.JSON(http.StatusTeapot, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"count": len(smss),
				"smss":  smss,
			})
		}
	})

	r.POST("/sms/:name", func(c *gin.Context) {
		var sms db.SmsMessage
		if err := c.ShouldBindJSON(&sms); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		name := c.Param("name")
		err := dao.Save(name, sms)
		if err != nil {
			c.JSON(http.StatusTeapot, gin.H{"error": err.Error()})
		} else {
			c.String(http.StatusOK, "Hello %s. You received an SMS from %s, saying %s", name, sms.Phone, sms.Content)
		}
	})
}
