package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tiki-systems/tikiserver/actions"
	"github.com/tiki-systems/tikiserver/middleware"
	"github.com/tiki-systems/tikiserver/models"
)

func ListDomains(c *gin.Context, vars middleware.GinHandlerVars) {
	domainRepository := vars.DomainRepository

	domains, err := domainRepository.ListDomains()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"count":    len(domains),
		"domains":  domains,
		"newToken": newToken,
	})
}

func GetDomain(c *gin.Context, vars middleware.GinHandlerVars) {
	domainRepository := vars.DomainRepository

	domainPath := c.Param("domainPath")
	domain, err := domainRepository.GetDomain(domainPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"domain":   domain,
		"newToken": newToken,
	})
}

func DeleteDomain(c *gin.Context, vars middleware.GinHandlerVars) {
	domainRepository := vars.DomainRepository
	permissionRepository := vars.PermissionRepository

	var domainPath map[string]string

	err := c.ShouldBindJSON(&domainPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Domain name couldn't found. Check your request again.",
		})
		return
	}

	userEmail := c.Request.Header.Get("email")

	if !permissionRepository.CanUserPerformDomainOperation(userEmail, models.Operation.Delete) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Permission denied on domain operation: user: [%v], operation : [%v]", userEmail, models.Operation.Delete),
		})
		return
	}

	if !permissionRepository.IsUserAllowedByDomainScope(userEmail, domainPath["domainPath"], models.DomainScopeOperation.DeleteDomain) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("This domain [%v] is not allowed to delete by [%v] due to group permissions", domainPath, userEmail),
		})
		return
	}

	if domainPath["domainPath"] == "/" || domainPath["domainPath"] == "/tickets" || domainPath["domainPath"] == "tickets" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Protected domains can't be deleted"),
		})
		return
	}

	dPath := actions.TrimDomainPath(domainPath["domainPath"])

	if !domainRepository.DoesTicketDomainExist(dPath) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("TicketDomain, [%v] does NOT exist", dPath),
		})
		return
	}

	err = domainRepository.DeleteDomain(dPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  fmt.Sprintf("Domain path [%s] is deleted successfully", domainPath),
		"newToken": newToken,
	})
}

func CreateDomain(c *gin.Context, vars middleware.GinHandlerVars) {
	domainRepository := vars.DomainRepository
	permissionRepository := vars.PermissionRepository

	var domain models.TicketDomain
	err := c.ShouldBindJSON(&domain)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Domain name couldn't found. Check your request again.",
		})
		return
	}

	if domainRepository.DoesTicketDomainExist(domain.DomainPath) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("TicketDomain,[%v] does already exist", domain.DomainPath),
		})
		return
	}

	userEmail := c.Request.Header.Get("email")
	if !permissionRepository.CanUserPerformDomainOperation(userEmail, models.Operation.Create) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Permission denied on domain operation: user: [%v], operation : [%v]", userEmail, models.Operation.Create),
		})
		return
	}

	if !permissionRepository.IsUserAllowedByDomainScope(userEmail, domain.DomainPath, models.DomainScopeOperation.CreateDomain) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("This domain [%v] is not allowed to created by [%v] due to group permissions", domain.DomainPath, userEmail),
		})
		return
	}

	domain.Parent = actions.UpdateParent(domain.DomainPath)
	parentObj, err := domainRepository.GetDomain(domain.Parent)
	if err != nil || parentObj.DomainPath == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Parent domain must exists, unless it is not top level domain",
		})
		return
	}
	domain.CreatedBy = c.Request.Header.Get("email")
	domain.UpdatedBy = c.Request.Header.Get("email")

	err = domainRepository.CreateDomain(domain)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  fmt.Sprintf("Domain path [%s] is created successfully", domain.DomainPath),
		"newToken": newToken,
	})
}
