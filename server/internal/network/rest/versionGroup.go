package rest

import "github.com/gin-gonic/gin"

func (r *Rest) versionGroup(root *gin.RouterGroup, groupName string) *gin.RouterGroup {
	group := root.Group(groupName)
	{
		r.devicesGroup(group, "devices")
	}

	return group
}
