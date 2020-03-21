package main

import (
	"callr/internal/config"
	"callr/internal/dao"
	"callr/internal/resources"
	"callr/internal/resources/static"
	"crypto/subtle"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

var db dao.Dao
var cfg config.Config

func main() {

	cfg = config.Get()
	db = dao.New(cfg.DataPath)

	fmt.Println("---- API ----")
	fmt.Println("Incident reporting: POST/GET:", cfg.BaseURL+"/incident?token="+cfg.IncidentToken)
	fmt.Println("    Incident is rotten after:", cfg.IncidentRottenDuration)

	fmt.Println("\n---- GUI ----")
	fmt.Println("Page at:", cfg.BaseURL)

	basicAuth := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		validUser := false
		validPass := false
		if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.BasicAuthUser)) == 1 {
			validUser = true
		}
		if subtle.ConstantTimeCompare([]byte(password), []byte(cfg.BasicAuthPass)) == 1 {
			validPass = true
		}
		return validUser && validPass, nil
	})

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.POST("/switchboard/page", resources.Page(db, cfg))
	e.POST("/switchboard/test-call", resources.TestCall())
	e.POST("/incident", resources.Incident(db, cfg))
	e.GET("/incident", resources.Incident(db, cfg))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	e.GET("/", static.Index, basicAuth)
	e.GET("/manage", static.Manage, basicAuth)

	e.GET("/api/incident", resources.GetActiveIncident(db), basicAuth)
	e.DELETE("/api/incident", resources.DeleteActiveIncident(db), basicAuth)

	e.GET("/api/incidents", resources.GetOldIncidents(db), basicAuth)
	e.GET("/api/incidents/:id/log", resources.GetIncidentLogs(db), basicAuth)

	e.GET("/api/people", resources.GetPeople(db), basicAuth)
	e.POST("/api/people", resources.PostPeople(db), basicAuth)
	e.GET("/api/oncall", resources.GetOnCall(db), basicAuth)
	e.POST("/api/oncall", resources.PostOnCall(db), basicAuth)

	e.GET("/api/test-call/:phone", resources.GetTestCall(cfg), basicAuth)

	if cfg.AutoTLS {
		go func() {
			ee := echo.New()
			ee.Pre(middleware.HTTPSRedirect())
			ee.Logger.Fatal(ee.Start(fmt.Sprintf(":%d", cfg.PortHTTP)))
		}()

		if len(cfg.AutoTLSDir) > 0 {
			e.AutoTLSManager.Cache = autocert.DirCache(cfg.AutoTLSDir)
		}
		if len(cfg.AutoTLSDomains) > 0 {
			e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(cfg.AutoTLSDomains...)
		}
		if len(cfg.AutoTLSEmail) > 0 {
			e.AutoTLSManager.Email = cfg.AutoTLSEmail
		}
		e.Logger.Fatal(e.StartAutoTLS(fmt.Sprintf(":%d", cfg.PortHTTPS)))
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.PortHTTP)))

}
