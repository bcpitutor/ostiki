package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/models"
	"github.com/gin-gonic/gin"
)

func GetTicket(c *gin.Context, vars middleware.GinHandlerVars) {
	//groupRepository := vars.GroupRepository
	ticketRepository := vars.TicketRepository
	permissionRepository := vars.PermissionRepository

	var ticketRequest map[string]string

	err := c.ShouldBindJSON(&ticketRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "check your request.",
			"status":  "error",
		})
		return
	}

	userEmail := c.Request.Header.Get("email")
	ticketPath := ticketRequest["ticketPath"]

	if !permissionRepository.CanUserAccessToTicket(userEmail, ticketPath) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status": "error",
			"message": fmt.Sprintf("Permission denied. User %s has been denied to use the ticket: [%s]",
				userEmail,
				ticketPath),
		})
		return
	}

	if !permissionRepository.CanUserPerformTicketOperation(userEmail, models.Operation.Info) {
		message := fmt.Sprintf("Permission denied. User: %s cannot perform operation [%s] on ticket [%s]", userEmail, models.Operation.Info, ticketPath)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	newToken, _ := c.Get("newToken")

	ticket, err := ticketRepository.QueryTicketByPath(ticketPath)
	if err != nil {
		message := fmt.Sprintf("Ticket not found for path: %s", ticketPath)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"data":     ticket,
		"newToken": newToken,
	})
}

func ListTickets(c *gin.Context, vars middleware.GinHandlerVars) {
	ticketRepository := vars.TicketRepository
	permissionRepository := vars.PermissionRepository

	// -- Control Section -- //
	userEmail := c.Request.Header.Get("email")
	if !permissionRepository.CanUserPerformTicketOperation(userEmail, models.Operation.List) {
		message := fmt.Sprintf("Permission denied. user: [%v], operation : [%v]", userEmail, models.Operation.List)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}
	// -- end of control section -- //

	tickets, err := ticketRepository.GetAllTickets()
	if err != nil {
		message := fmt.Sprintf("%v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	count := len(tickets)
	newToken, _ := c.Get("newToken")

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"count":    count,
		"tickets":  tickets,
		"newToken": newToken,
	})
}

func CreateTicket(c *gin.Context, vars middleware.GinHandlerVars) {
	ticketRepository := vars.TicketRepository
	permissionRepository := vars.PermissionRepository

	userEmail := c.Request.Header.Get("email")
	if !permissionRepository.CanUserPerformTicketOperation(userEmail, models.Operation.Create) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status": "error",
			"message": fmt.Sprintf("Permisson denied. user: [%v], operation : [%v]",
				userEmail,
				models.Operation.Create,
			),
		})
		return
	}

	var newTicket models.Ticket
	err := c.ShouldBindJSON(&newTicket)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"details": err,
			"message": "check your json post data 1",
		})
		return
	}

	if ticketRepository.DoesTicketExist(newTicket.TicketPath) {
		message := fmt.Sprintf("This ticket [%v] does already exist", newTicket.TicketPath)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if !permissionRepository.IsUserAllowedByDomainScope(userEmail, newTicket.TicketPath, models.DomainScopeOperation.CreateTicket) {
		message := fmt.Sprintf("This ticket [%v] is not allowed to created by [%v] due to group permissons", newTicket.TicketPath, userEmail)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	newTicket.CreatedBy = userEmail
	newTicket.UpdatedBy = userEmail

	err = ticketRepository.CreateTicket(newTicket)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failure during the creation dynmodb ticket record.",
			"details": err,
		})
		return
	}

	newToken, _ := c.Get("newToken")
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "a new ticket created on dynamodb",
		"details":  fmt.Sprintf("Ticket created: %s", newTicket.TicketPath),
		"newToken": newToken,
	})
}

func DeleteTicket(c *gin.Context, vars middleware.GinHandlerVars) {
	var tkt map[string]string
	err := c.ShouldBindJSON(&tkt)
	permissionRepository := vars.PermissionRepository
	ticketRepository := vars.TicketRepository

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "check your json post data 2",
		})
		return
	}

	ticketPath := tkt["ticketPath"]

	if ticketPath == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ticketPath is empty",
		})
		return
	}

	userEmail := c.Request.Header.Get("email")

	if !permissionRepository.CanUserAccessToTicket(userEmail, ticketPath) {
		message := fmt.Sprintf("Permisson denied. user: [%v], ticket: [%v]", userEmail, ticketPath)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if !permissionRepository.CanUserPerformTicketOperation(userEmail, models.Operation.Delete) {
		message := fmt.Sprintf("Permission denied. user: [%v], operation : [%v]", userEmail, models.Operation.Delete)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if !permissionRepository.IsUserAllowedByDomainScope(userEmail, ticketPath, models.DomainScopeOperation.DeleteTicket) {
		message := fmt.Sprintf("This ticket [%v] is not allowed to deleted by [%v] due to group permissons", ticketPath, userEmail)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	ticketData, err := ticketRepository.QueryTicketByPath(ticketPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failure during the query dynamodb ticket record.",
			"details": err,
		})
		return
	}

	var ticket models.Ticket
	ticketBytes, err := json.Marshal(ticketData)

	if err != nil {
		message := fmt.Sprintf("deleting ticket -> find ticket unmarshalling error : [%v]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": message,
			"details": err,
		})
		return
	}

	if err := json.Unmarshal(ticketBytes, &ticket); err != nil {
		message := fmt.Sprintf("deleting ticket -> find ticket unmarshalling error : [%v]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if ticket.TicketType == "" {
		message := fmt.Sprintf("delete ticket, ticket type not found for the ticket, [%v]", ticketPath)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	err = ticketRepository.DeleteTicket(ticketPath, ticket.TicketType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Delete request is failed.",
			"details": err,
		})
		return
	}

	newToken, _ := c.Get("newToken")
	message := fmt.Sprintf("ticket, [%v] is deleted.", tkt["ticketPath"])
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  message,
		"newToken": newToken,
	})
}

func ObtainTicket(c *gin.Context, vars middleware.GinHandlerVars) {
	var rBody map[string]string
	permissionRepository := vars.PermissionRepository
	ticketRepository := vars.TicketRepository
	awsService := vars.AWSService

	fmt.Printf("we are in obtain ticket\n")
	if err := c.ShouldBindJSON(&rBody); err != nil {
		message := fmt.Sprintf("%v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}
	userEmail := c.Request.Header.Get("email")
	tPath := rBody["ticketPath"]

	if !permissionRepository.CanUserAccessToTicket(userEmail, tPath) {
		message := fmt.Sprintf("Permisson denied. user: [%v], ticket: [%v]", userEmail, tPath)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if !permissionRepository.CanUserPerformTicketOperation(userEmail, models.Operation.Show) {
		message := fmt.Sprintf("Permission denied. user: [%v], operation : [%v]", userEmail, models.Operation.Show)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	requestedTicketInfo, err := ticketRepository.QueryTicketByPath(tPath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": err,
		})
		return
	}

	if !permissionRepository.IsUserAllowedByDomainScope(userEmail, tPath, "assumeTicket") {
		message := fmt.Sprintf("Permission denied for the user to assume this ticket: user: [%v], ticket : [%v]", userEmail, tPath)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	authorizationHeader := c.Request.Header.Get("Authorization")
	userEmail = c.Request.Header.Get("email")
	tokenArr := strings.Split(authorizationHeader, " ")
	token := tokenArr[1]

	ticketRoleArn := requestedTicketInfo.AwsAssumeRole.RoleArn
	ticketTtl := requestedTicketInfo.AwsAssumeRole.Ttl
	ticketRegion := requestedTicketInfo.TicketRegion

	result, err := awsService.ObtainAWSRoleWithToken(
		token,
		userEmail,
		ticketRoleArn,
		ticketTtl,
		ticketRegion,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "AWS error",
			"err":     fmt.Sprintf("%v", err),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"AccessKeyId":     result.AccessKeyId,
		"SecretAccessKey": result.SecretAccessKey,
		"SessionToken":    result.SessionToken,
		"Region":          result.Region,
		"newToken":        token,
	})
}

func TicketSetSecret(c *gin.Context, vars middleware.GinHandlerVars) {
	ticketRepository := vars.TicketRepository
	permissionRepository := vars.PermissionRepository
	awsService := vars.AWSService

	type SecretDataType struct {
		TicketPath string
		SecretData string
	}

	var sdo SecretDataType
	if err := c.ShouldBindJSON(&sdo); err != nil {
		message := fmt.Sprintf("marshall error, check you post data: [%v] ", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "check your json post data 3",
			"details": message,
		})
		return
	}

	if !ticketRepository.DoesTicketExist(sdo.TicketPath) {
		message := fmt.Sprintf("This secret ticket [%v] does NOT exist", sdo.TicketPath)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	userEmail := c.Request.Header.Get("email")
	if !permissionRepository.CanUserPerformTicketOperation(userEmail, models.Operation.SetSecret) {
		message := fmt.Sprintf("Permission denied. user: [%v], operation : [%v]", userEmail, models.Operation.Show)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if !permissionRepository.IsUserAllowedByDomainScope(userEmail, sdo.TicketPath, models.DomainScopeOperation.CreateTicket) {
		message := fmt.Sprintf("This ticket [%v] is not allowed to created by [%v] due to group permissons", sdo.TicketPath, userEmail)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	encryptedData, err := awsService.GetEncryptedSecret(sdo.SecretData)
	if err != nil {
		message := fmt.Sprintf("Encryption error: [%+v]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	err = ticketRepository.SetTicketSecret(sdo.TicketPath, encryptedData)
	if err != nil {
		message := fmt.Sprintf("Error during setting secret data on ticket, error: [%+v]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	message := fmt.Sprintf("Secret set on [%s] by [%s]", sdo.TicketPath, userEmail)
	newToken, _ := c.Get("newToken")
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  message,
		"newToken": newToken,
	})
}

func TicketGetSecret(c *gin.Context, vars middleware.GinHandlerVars) {
	ticketRepository := vars.TicketRepository
	permissionRepository := vars.PermissionRepository
	awsService := vars.AWSService

	type SecretDataType struct {
		TicketPath string
	}
	var secretTicketData SecretDataType

	if err := c.ShouldBindJSON(&secretTicketData); err != nil {
		message := fmt.Sprintf("marshall error, check you post data: [%v]", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"details": message,
			"message": "check your json post data 4",
		})
		return
	}

	if !ticketRepository.DoesTicketExist(secretTicketData.TicketPath) {
		message := fmt.Sprintf("This secret ticket [%v] does NOT exist", secretTicketData.TicketPath)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	userEmail := c.Request.Header.Get("email")
	if !permissionRepository.CanUserPerformTicketOperation(userEmail, models.Operation.GetSecret) {
		message := fmt.Sprintf("Permission denied. user: [%v], operation : [%v]", userEmail, models.Operation.Show)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	if !permissionRepository.IsUserAllowedByDomainScope(userEmail, secretTicketData.TicketPath, models.DomainScopeOperation.CreateTicket) {
		message := fmt.Sprintf("This ticket [%v] is not allowed to created by [%v] due to group permissons", secretTicketData.TicketPath, userEmail)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	secretData, err := ticketRepository.GetTicketSecret(secretTicketData.TicketPath)
	if err != nil {
		message := fmt.Sprintf("Error while retrieving ticket: [%v]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"details": message,
			"message": "check your json post data 5",
		})
		return
	}

	fmt.Printf("Secret data: [%v]\n", secretData)
	if secretData == "" {
		message := fmt.Sprintf("Secret data is empty on ticket [%v]", secretTicketData.TicketPath)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": message,
		})
	}

	clearText, err := awsService.GetDecryptedText(secretData)
	if err != nil {
		message := fmt.Sprintf("Decryption error: [%v]", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"details": message,
			"message": "Decryption error, secret data is not valid",
		})
		return
	}

	newToken, _ := c.Get("newToken")
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  clearText,
		"newToken": newToken,
	})
}
