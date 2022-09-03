package main

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type SmsMessage struct {
	XMLName  xml.Name `json:"Message,omitempty"`
	Smstat   uint     `json:"Smstat,omitempty"`
	Index    uint     `json:"Index,omitempty"`
	Phone    string   `json:"Phone,omitempty"`
	Content  string   `json:"Content,omitempty"`
	Date     string   `json:"Date,omitempty"`
	Sca      any      `json:"Sca,omitempty"`
	SaveType uint     `json:"SaveType,omitempty"`
	Priority uint     `json:"Priority,omitempty,omitempty"`
	SmsType  uint     `json:"SmsType,omitempty"`
}

func main() {
	r := gin.Default()
	db := CreateDao()

	r.GET("/mailhub", func(c *gin.Context) {
		smss, err := db.GetAllSmss()
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
		smss, err := db.GetSmssTo(name)

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
		var sms SmsMessage
		if err := c.ShouldBindJSON(&sms); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		name := c.Param("name")
		err := db.Save(name, sms)
		if err != nil {
			c.JSON(http.StatusTeapot, gin.H{"error": err.Error()})
		} else {
			c.String(http.StatusOK, "Hello %s. You received an SMS from %s, saying %s", name, sms.Phone, sms.Content)
		}
	})

	v, exists := os.LookupEnv("SERVER_ADDR")
	if exists {
		r.Run(v)
	} else {
		r.Run()
	}
}
