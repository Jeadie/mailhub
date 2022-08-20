package main

import (
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"net/http"
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

	db := Dao{make(map[string][]SmsMessage)}

	r.GET("/mailhub", func(c *gin.Context) {
		smss, _ := db.GetAllSmss()
		c.JSON(http.StatusOK, gin.H{
			"count": len(smss),
			"smss":  smss,
		})
	})

	r.GET("/sms/:name", func(c *gin.Context) {
		name := c.Param("name")
		smss, _ := db.GetSmssTo(name)

		c.JSON(http.StatusOK, gin.H{
			"count": len(smss),
			"smss":  smss,
		})
	})

	r.POST("/sms/:name", func(c *gin.Context) {
		var sms SmsMessage
		if err := c.ShouldBindJSON(&sms); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		name := c.Param("name")
		db.Save(name, sms)
		c.String(http.StatusOK, "Hello %s. You received an SMS from %s, saying %s", name, sms.Phone, sms.Content)
	})

	r.Run()
}
