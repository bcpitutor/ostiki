package routes

import (
	"fmt"
	"net/http"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/models"
	"github.com/gin-gonic/gin"
)

func ListGroups(c *gin.Context, vars middleware.GinHandlerVars) {
	groupRepository := vars.GroupRepository
	imoRepository := vars.ImoRepository
	sugar := vars.Logger.Sugar()

	var groups []models.TicketGroup

	groupsFromImo := imoRepository.GetGroups()
	if len(groupsFromImo) != 0 {
		sugar.Infof("Got groups in imo: %+v", groupsFromImo)
		groups = groupsFromImo
	} else {
		sugar.Infof("Groups are not in imo, reading from DB")
		groupsFromDB, err := groupRepository.GetAllGroups()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("%v", err),
			})
			c.Abort()
			return
		}
		groups = groupsFromDB
		imoRepository.SetGroups(groups)
	}

	newToken, _ := c.Get("newToken")
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"count":    len(groups),
		"data":     groups,
		"newToken": newToken,
	})
}

func GetGroup(c *gin.Context, vars middleware.GinHandlerVars) {
	groupRepository := vars.GroupRepository

	group, err := groupRepository.GetGroup(c.Param("groupName"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Group name couldn't found. Check your request again.",
		})
		c.Abort()
		return
	}

	newToken, _ := c.Get("newToken")
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"data":     group,
		"newToken": newToken,
	})
}

func CreateGroup(c *gin.Context, vars middleware.GinHandlerVars) {
	groupRepository := vars.GroupRepository
	permissionRepository := vars.PermissionRepository

	var newGroup models.TicketGroup
	err := c.ShouldBindJSON(&newGroup)
	if err != nil {
		message := fmt.Sprintf("marshall error, check you post data: [%v] ", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message,
		})
	}

	userEmail := c.Request.Header.Get("email")
	if !permissionRepository.CanUserPerformGroupOperation(userEmail, models.Operation.Create) {
		message := fmt.Sprintf("Permission denied on group operation: user: [%v], operation : [%v]", userEmail, models.Operation.Create)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}
	groupName := newGroup.GroupName
	if groupRepository.DoesGroupExist(groupName) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Group [%s] already exists", groupName),
		})
		return
	}

	email := c.GetHeader("email")
	newGroup.CreatedBy, newGroup.UpdatedBy = email, email

	var defaultPerms models.Aperms

	defaultPerms.Group = map[string]bool{
		"create":    false,
		"addMember": false,
		"delMember": false,
		"delete":    false,
	}

	defaultPerms.Domain = map[string]bool{
		"create": true,
		"delete": false,
	}

	defaultPerms.Ticket = map[string]bool{
		"create": true,
		"delete": false,
	}

	defaultPerms.SecretTicket = map[string]bool{
		"read":   false,
		"create": false,
		"delete": false,
	}
	newGroup.AccessPerms = defaultPerms

	err = groupRepository.CreateGroup(newGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	message := fmt.Sprintf("group, [%v] is created.", newGroup.GroupName)
	newToken, _ := c.Get("newToken")
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  message,
		"newToken": newToken,
	})
}

func DeleteGroup(c *gin.Context, vars middleware.GinHandlerVars) {
	groupRepository := vars.GroupRepository
	permissionRepository := vars.PermissionRepository

	groupName, ok := c.Params.Get("groupName")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "groupName is required",
		})
		c.Abort()
		return
	}
	userEmail := c.Request.Header.Get("email")
	if userEmail == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "identity couldn't be found",
		})
		return
	}

	if !permissionRepository.CanUserPerformGroupOperation(userEmail, models.Operation.Delete) {
		message := fmt.Sprintf("Permission denied on group operation: user: [%v], operation : [%v]", userEmail, models.Operation.Delete)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if !groupRepository.DoesGroupExist(groupName) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Group couldn't be found",
		})
		return
	}

	err := groupRepository.DeleteGroup(groupName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	message := fmt.Sprintf("group, [%v] is deleted.", groupName)
	newToken, _ := c.Get("newToken")
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  message,
		"newToken": newToken,
	})
}

func AddMemberToGroup(c *gin.Context, vars middleware.GinHandlerVars) {
	groupRepository := vars.GroupRepository
	permissionRepository := vars.PermissionRepository

	groupName, ok := c.Params.Get("groupName")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "groupName is required",
		})
		c.Abort()
		return
	}

	userEmail := c.Request.Header.Get("email")
	if userEmail == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "identity couldn't be found",
		})
		return
	}

	var newMember map[string]string
	err := c.ShouldBindJSON(&newMember)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("marshall error, check you post data: [%v] ", err),
		})
		return
	}

	if !permissionRepository.CanUserPerformGroupOperation(userEmail, models.Operation.AddMember) {
		message := fmt.Sprintf("Permission denied on group operation: user: [%v], operation : [%v]", userEmail, models.Operation.AddMember)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if !groupRepository.DoesGroupExist(groupName) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Group couldn't be found",
		})
		return
	}

	err = groupRepository.AddMemberToGroup(newMember["newMemberEmail"], groupName, userEmail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "new member couldn't be added to the group",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  fmt.Sprintf("new member [%v] is added", newMember["newMemberEmail"]),
		"newToken": newToken,
	})
}

func DelMemberFromGroup(c *gin.Context, vars middleware.GinHandlerVars) {
	groupRepository := vars.GroupRepository
	permissionRepository := vars.PermissionRepository

	groupName, ok := c.Params.Get("groupName")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "groupName is required",
		})
		c.Abort()
		return
	}

	userEmail := c.Request.Header.Get("email")
	if userEmail == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "identity couldn't be found",
		})
		return
	}

	var memberToRemove map[string]string
	err := c.ShouldBindJSON(&memberToRemove)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("marshall error, check you post data: [%v] ", err),
		})
		return
	}

	if !permissionRepository.CanUserPerformGroupOperation(userEmail, models.Operation.DelMember) {
		message := fmt.Sprintf("Permission denied on group operation: user: [%v], operation : [%v]", userEmail, models.Operation.DelMember)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if !groupRepository.DoesGroupExist(groupName) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Group couldn't be found",
		})
		return
	}

	if !groupRepository.IsUserMemberOfGroup(memberToRemove["deleteMemberEmail"], groupName) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%s is not a member of this group", memberToRemove["deleteMemberEmail"]),
		})
		return
	}

	err = groupRepository.DelMemberFromGroup(memberToRemove["deleteMemberEmail"], groupName, userEmail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Member couldn't be removed from the group",
			"details": fmt.Sprintf("%v", err),
		})
		return
	}

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  fmt.Sprintf("the member [%v] is deleted", memberToRemove["deleteMemberEmail"]),
		"newToken": newToken,
	})
}
