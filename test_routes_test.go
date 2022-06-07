package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bcpitutor/ostiki/appconfig"
	"github.com/bcpitutor/ostiki/db"
	"github.com/bcpitutor/ostiki/logger"
	"github.com/bcpitutor/ostiki/middleware"
	"github.com/bcpitutor/ostiki/models"
	"github.com/bcpitutor/ostiki/repositories"
	"github.com/bcpitutor/ostiki/routes"
	"github.com/gin-gonic/gin"
)

func TestRouteListTickets(t *testing.T) {
	r := gin.Default()

	r.GET(
		"/tickets",
		func(ctx *gin.Context) {
			config := appconfig.GetAppConfig()
			logger := logger.GetTikiLogger(config)

			dynamo, err := db.NewDynamoDBDriver(config, logger)
			if err != nil {
				t.Fatalf("Failed to create DynamoDB driver: %v", err)
			}
			grt := repositories.ProvideGroupRepository(dynamo, nil)
			trt := repositories.ProvideTicketRepository(dynamo, nil)
			prt := repositories.ProvidePermissionRepository(dynamo, grt.GroupRepository, trt.TicketRepository)

			routes.ListTickets(
				ctx,
				middleware.GinHandlerVars{
					Logger:               logger.Logger,
					AppConfig:            config,
					GroupRepository:      grt.GroupRepository,
					TicketRepository:     trt.TicketRepository,
					PermissionRepository: prt.PermissionRepository,
				},
			)
		},
	)
	req, err := http.NewRequest(http.MethodGet, "/tickets", nil)
	req.Header.Set("email", "ozgur.demir@itutor.com")
	if err != nil {
		t.Fatalf("Couldn't create http request : %v\n", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		type JsonRes struct {
			Count    int    `json:"count"`
			NewToken string `json:"newToken"`
			Success  string `json:"success"`
			Tickets  []any  `json:"tickets"`
		}

		var j JsonRes

		err := json.Unmarshal(w.Body.Bytes(), &j)
		if err != nil {
			t.Errorf("Couldn't unmarshal json : %v\n", err)
		}

		if j.Count != 0 {
			t.Logf("Got %d tickets", j.Count)
		} else {
			t.Errorf("Could't get ticket list")
		}
	} else {
		t.Errorf("Expected status code 200, got %v\n", w.Code)
	}
}

func TestRouteGetTicket(t *testing.T) {
	r := gin.Default()

	r.POST(
		"/ticket/get",
		func(ctx *gin.Context) {
			config := appconfig.GetAppConfig()
			logger := logger.GetTikiLogger(config)

			dynamo, err := db.NewDynamoDBDriver(config, logger)
			if err != nil {
				t.Fatalf("Failed to create DynamoDB driver: %v", err)
			}
			grt := repositories.ProvideGroupRepository(dynamo, nil)
			trt := repositories.ProvideTicketRepository(dynamo, nil)
			prt := repositories.ProvidePermissionRepository(dynamo, grt.GroupRepository, trt.TicketRepository)
			routes.GetTicket(
				ctx,
				middleware.GinHandlerVars{
					Logger:               logger.Logger,
					AppConfig:            config,
					GroupRepository:      grt.GroupRepository,
					TicketRepository:     trt.TicketRepository,
					PermissionRepository: prt.PermissionRepository,
				},
			)
		},
	)

	var a = `{"ticketPath":"tickets/itutor/infra/prod/admin/aws"}`
	reader := strings.NewReader(a)
	req, err := http.NewRequest(http.MethodPost, "/ticket/get", reader)
	req.Header.Set("email", "ozgur.demir@itutor.com")
	if err != nil {
		t.Fatalf("Couldn't create http request : %v\n", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		type JsonRes struct {
			Status   string        `json:"status"`
			Data     models.Ticket `json:"data"`
			NewToken string        `json:"newToken"`
		}

		var j JsonRes

		err := json.Unmarshal(w.Body.Bytes(), &j)
		if err != nil {
			t.Errorf("Couldn't unmarshal json : %v\n", err)
		}

		if j.Data.TicketPath == "tickets/itutor/infra/prod/admin/aws" {
			t.Logf("Got ticket %v", j.Data.TicketPath)
		} else {
			t.Errorf("Could't get ticket %v", j.Data.TicketPath)
		}
	} else {
		t.Errorf("Expected status code 200, got %+v\n", w)
	}
}

// func TestRouteListDomains(t *testing.T) {
// 	r := gin.Default()

// 	r.GET(
// 		"/domain/list",
// 		func(ctx *gin.Context) {
// 			config := appconfig.GetAppConfig()
// 			logger := logger.GetTikiLogger(config)

// 			dynamo, err := db.NewDynamoDBDriver(config, logger)
// 			if err != nil {
// 				t.Fatalf("Failed to create DynamoDB driver: %v", err)
// 			}
// 			grt := repositories.ProvideGroupRepository(dynamo, nil)
// 			trt := repositories.ProvideTicketRepository(dynamo, nil)
// 			prt := repositories.ProvidePermissionRepository(dynamo, grt.GroupRepository, trt.TicketRepository)

// 			routes.ListDomains(
// 				ctx,
// 				middleware.GinHandlerVars{
// 					Logger:               logger.Logger,
// 					AppConfig:            config,
// 					GroupRepository:      grt.GroupRepository,
// 					TicketRepository:     trt.TicketRepository,
// 					PermissionRepository: prt.PermissionRepository,
// 				},
// 			)
// 		},
// 	)
// 	req, err := http.NewRequest(http.MethodGet, "/domain/list", nil)
// 	req.Header.Set("email", "ozgur.demir@itutor.com")
// 	if err != nil {
// 		t.Fatalf("Couldn't create http request : %v\n", err)
// 	}
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	if w.Code == http.StatusOK {
// 		type JsonRes struct {
// 			Count    int    `json:"count"`
// 			NewToken string `json:"newToken"`
// 			Success  string `json:"success"`
// 			Domains  []any  `json:"domains"`
// 		}

// 		var j JsonRes

// 		err := json.Unmarshal(w.Body.Bytes(), &j)
// 		if err != nil {
// 			t.Errorf("Couldn't unmarshal json : %v\n", err)
// 		}

// 		if j.Count != 0 {
// 			t.Logf("Got %d domains", j.Count)
// 		} else {
// 			t.Errorf("Could't get domain list")
// 		}
// 	} else {
// 		t.Errorf("Expected status code 200, got %v\n", w.Code)
// 	}

// }
