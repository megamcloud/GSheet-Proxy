package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func QRCheck(c *gin.Context) {
	c.HTML(http.StatusOK, "first_checkin.tmpl", gin.H{
		"name":      "Tran Toan Van",
		"company":   "Anphabe",
		"job_title": "Chief Tumlum Tala officer",
	})
}

func Hello(c *gin.Context) {
	c.HTML(http.StatusOK, "first_checkin.tmpl", gin.H{
		"name":      "Tran Toan Van",
		"company":   "Anphabe",
		"job_title": "Chief Tumlum Tala officer",
	})
}
