package rest

import (
	"net/http"
	"remoteesp/internal/domain"
	"remoteesp/internal/domain/env"

	"github.com/gin-gonic/gin"
)

func (r *Rest) rootGroup(root *gin.RouterGroup) *gin.RouterGroup {
	root.GET("", r.getRoot)
	root.GET("/favicon.ico", r.getFavoriteIco)
	root.GET("devices", r.getDevicesPage)
	root.GET("devices/info", r.getDeviceInfPage)
	root.GET("devices/cfg", r.getDeviceCfgPage)
	root.GET("devices/upload", r.getUploadPage)
	root.GET("login", r.getLoginPage)
	return root
}

func (r *Rest) getRoot(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (r *Rest) getDevicesPage(c *gin.Context) {
	filter := c.Query("filter")
	xToken := c.Query("token")
	devices, err := r.db.Devices(filter)
	if err != nil {
		domain.Log.Log(err)
		ErrorHtml[Empty](c)
		return
	}

	switch xToken {
	case env.XToken:
		c.HTML(http.StatusOK, "xdevices.html", gin.H{
			"devices": devices,
		})
	default:
		c.HTML(http.StatusOK, "devices.html", gin.H{
			"devices": devices,
		})
	}
}

func (r *Rest) getDeviceInfPage(c *gin.Context) {
	filter := c.Query("ssdp")
	if filter == "" {
		ErrorHtml[Empty](c)
		return
	}

	xToken := c.Query("token")
	devices, err := r.db.Devices(filter)
	if err != nil {
		domain.Log.Log(err)
		ErrorHtml[Empty](c)
		return
	}

	switch xToken {
	case env.XToken:
		c.HTML(http.StatusOK, "deviceInfo.html", gin.H{
			"device": devices[0],
		})
	default:
		ErrorHtml[Empty](c)
	}
}

func (r *Rest) getDeviceCfgPage(c *gin.Context) {
	filter := c.Query("ssdp")
	if filter == "" {
		ErrorHtml[Empty](c)
		return
	}

	xToken := c.Query("token")
	devices, err := r.db.Devices(filter)
	if err != nil {
		domain.Log.Log(err)
		ErrorHtml[Empty](c)
		return
	}

	switch xToken {
	case env.XToken:
		c.HTML(http.StatusOK, "deviceCfg.html", gin.H{
			"device": devices[0],
		})
	default:
		ErrorHtml[Empty](c)
	}
}

func (r *Rest) getUploadPage(c *gin.Context) {

	xToken := c.Query("token")

	switch xToken {
	case env.XToken:
		c.HTML(http.StatusOK, "upload.html", xToken)
	default:
		ErrorHtml[Empty](c)
	}
}

func (r *Rest) getLoginPage(c *gin.Context) {

	c.HTML(http.StatusOK, "login.html", nil)
}

func (r *Rest) getFavoriteIco(c *gin.Context) {
	c.File("./templates/static/favicon.ico")
}
