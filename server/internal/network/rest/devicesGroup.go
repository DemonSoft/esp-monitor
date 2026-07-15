package rest

import (
	"encoding/json"
	"remoteesp/internal/domain"
	"remoteesp/internal/model"

	"github.com/gin-gonic/gin"
)

type Database interface {
	DevicesUploaded(devices []model.Device) error
	Devices(filter string) ([]model.Device, error)
	DeleteDevice(ssdp string) error
}

func (r *Rest) devicesGroup(root *gin.RouterGroup, groupName string) *gin.RouterGroup {
	group := root.Group(groupName)
	{
		group.GET("", r.getDeviceList)
		group.POST("upload", r.uploadDevicesCSV)
		group.POST("cfg", r.postDeviceCfg)
		group.DELETE("", r.deleteDevice)
	}

	return group
}

func (r *Rest) getDeviceList(c *gin.Context) {
	filter := c.Query("filter")
	devices, err := r.db.Devices(filter)
	if err != nil {
		domain.Log.Log(err)
		SomeError[Empty](c, 500)
		return
	}
	Success(c, devices)
}

func (r *Rest) postDeviceCfg(c *gin.Context) {

	currentDevice := c.GetHeader("SSDP")

	var config model.ConfigPayload

	err := json.NewDecoder(c.Request.Body).Decode(&config)
	if err != nil {
		RestError[Empty](c, "JSON decoding error: "+err.Error(), "Bad Request")
		return
	}

	r.mqtt.SendDeviceConfig(currentDevice, config)
	Success(c, "Ok")
}

func (r *Rest) deleteDevice(c *gin.Context) {
	type Request struct {
		SSDP string `json:"ssdp" binding:"required"`
	}
	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		// Если JSON некорректен или нет обязательного поля ssdp
		RestError[Empty](c, "The ssdp parameter must be set: "+err.Error(), "DELETING_NOT_SET_SSDP")
		return
	}

	err := r.db.DeleteDevice(req.SSDP)
	if err != nil {
		RestError[Empty](c, "Device wasn't deleted", "DELETING_FAIL")
		return
	}
	Success(c, "Ok")
}
