package api

import (
	"NULL/knowledgebase/pkg/util"
	"github.com/gin-gonic/gin"
	"html/template"
)

func KlibIndex(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}

func JkxmIndex(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static_jkxm/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}

func KfqZmq(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static_kfqzmq/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}

func Kfqky(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static_kfqky/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}

func Device(c *gin.Context) {
	if c.Request.Method == "GET" {
		t, _ := template.ParseFiles("runtime/static_device/index.html")
		util.ShowError("template parseFiles err", t.Execute(c.Writer, nil))
	}
}
