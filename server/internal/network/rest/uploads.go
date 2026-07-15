package rest

import (
	"remoteesp/internal/domain"
	"remoteesp/internal/model"

	"github.com/gin-gonic/gin"
)

// This methods relate to different REST groups.
// ---------------------------------------------

func (r *Rest) uploadDevicesCSV(c *gin.Context) {
	maxFileSize := int64(5 << 20) // 5 MB

	// Get file
	file, header, err := c.Request.FormFile("body")
	if err != nil {
		domain.Log.Log("ERROR: ", err)
		domain.Log.Log("header: ", header)
		RestError[Empty](c, "", "File was not uploaded")
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > maxFileSize {
		RestError[Empty](c, "File too large, max 5MB", "File is too large")
		return
	}

	// Check file size
	if header.Size == 0 {
		RestError[Empty](c, "File can't be empty", "File is empty")
		return
	}

	devices, err := model.ParseDevices(file)
	if err != nil {
		domain.Log.Log("Devices list not parsed:", err)
		SomeError[Empty](c, 500)
		return
	}

	err = r.db.DevicesUploaded(devices)
	if err != nil {
		domain.Log.Log("Devices list not saved:", err)
		SomeError[Empty](c, 500)
		return
	}

	Success(c, "OK")
}
