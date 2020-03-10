package main

import (
	"callr/internal/config"
	"callr/internal/dao"
	"callr/internal/resources"
	"crypto/subtle"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sfreiberg/gotwilio"
	"io/ioutil"
	"time"
)

var db dao.Dao
var cfg config.Config

func main() {

	cfg = config.Get()
	db.Store = cfg.DataPath

	fmt.Println("---- API ----")
	fmt.Println("Incident reporting: POST/GET:", cfg.BaseURL+"/incident?token="+cfg.IncidentToken)
	fmt.Println("      Close Incident: DELETE:", cfg.BaseURL+"/incident?token="+cfg.IncidentToken)

	fmt.Println("\n---- GUI ----")
	fmt.Println("Page at:", cfg.BaseURL)

	basicAuth := middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		validUser := false
		validPass := false
		if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.BasicAuthUser)) == 1{
			validUser = true
		}
		if subtle.ConstantTimeCompare([]byte(password), []byte(cfg.BasicAuthPass)) == 1 {
			validPass = true
		}
		return validUser && validPass, nil
	})

	e := echo.New()
	e.POST("/switchboard/page", Page)
	e.POST("/incident", Incident)
	e.GET("/incident", Incident)


	e.DELETE("/incident", CloseIncident, basicAuth)
	e.GET("/", func(c echo.Context) error {
		return c.Blob(200, "text/html", resources.Index())
	}, basicAuth)
	e.GET("/manage", func(c echo.Context) error {
		return c.Blob(200, "text/html", resources.Manage())
	}, basicAuth)

	e.GET("/api/incident", func(c echo.Context) error {
		is, err := db.GetIncident()
		if err == nil {
			return c.JSON(200, is)

		}
		if err.Error() == "no incident exists" {
			return c.JSON(200, nil)
		}
		return err
	},basicAuth)
	e.GET("/api/incidents", func(c echo.Context) error {
		is, err := db.GetIncidents()
		if err != nil {
			return err
		}
		return c.JSON(200, is)
	}, basicAuth)

	e.GET("/api/incidents/:id/log", func(c echo.Context) error {
		logs, err := db.GetLogs(c.Param("id"))
		if err != nil {
			return err
		}
		return c.JSON(200, logs)
	},basicAuth)

	e.Logger.Fatal(e.Start(":8080"))
}

func CloseIncident(c echo.Context) error {
	i, err := db.CloseIncident()
	if err != nil {
		return err
	}
	return c.JSON(200, i)
}

func Incident(c echo.Context) error {
	if c.QueryParams().Get("token") != cfg.IncidentToken {
		return c.String(400, "correct token was not provided")
	}

	err := db.NewIncident()
	if err == nil { // first once that creates it...
		go caller()
	}
	if err != nil && err.Error() != "incident already exist" {
		fmt.Println("[Error]", err)
	}

	i, err := db.GetIncident()
	if err != nil {
		fmt.Println("[Error]", err)
		return err
	}

	log := dao.Log{
		CreatedAt:   time.Now(),
		ContentType: c.Request().Header.Get("content-type"),
	}

	params := c.QueryParams()
	params.Del("token")
	log.Params = params

	if c.Request().Method == "POST" {
		b, err := ioutil.ReadAll(c.Request().Body)
		if err == nil {
			log.Body = string(b)
		}
	}
	return db.AddLog(i, log)
}

func caller() {
	for {
		i, err := db.GetIncident()
		if err != nil {
			fmt.Println("[Error]", err)
		}

		switch i.Status {
		case "Init":
			i.Status = "Calling"
			err = db.WriteIncident(i)
			if err != nil {
				fmt.Println("[Error]", err)
			}
			go mkcall(0)
		case "Calling":

		}
		<-time.After(time.Minute)
	}

}

func mkcall(inc int) {
	i, err := db.GetIncident()
	oncall, err := db.GetOncall()
	if err != nil {
		fmt.Println("[Error]", err)
		return
	}
	if len(oncall) == 0 {
		fmt.Println("[Error]", "no one in the oncall list")
		return
	}
	if i.OnCallIndex > 10 {
		return
	}

	i.OnCallIndex += inc
	err = db.WriteIncident(i)
	if err != nil {
		fmt.Println("[Error]", err)
		return
	}

	idx := i.OnCallIndex % len(oncall)
	p := oncall[idx]

	twilio := gotwilio.NewTwilioClient(cfg.TwilSID, cfg.TwilToken)
	_, errs, err := twilio.CallWithUrlCallbacks(cfg.TwilPhone, p.Phone, gotwilio.NewCallbackParameters(cfg.BaseURL+"/switchboard/page"))
	if err != nil {
		fmt.Println("[Error]", err)
		return
	}
	if errs != nil {
		fmt.Println("[Twil Error]", errs.Error())
		return
	}
}

func Page(c echo.Context) error {

	params, err := c.FormParams()

	if err != nil {
		fmt.Println(0, err)
		return err
	}

	fmt.Printf("PARAMS:\n%+v\n", params)

	callStatus := params.Get("CallStatus")
	called := params.Get("Called")
	msg := params.Get("msg")
	digits := params.Get("Digits")

	// Init call
	if callStatus == "in-progress" && msg == "" && digits == "" {
		resp := `<?xml version="1.0" encoding="UTF-8"?>
				<Response>
					<Gather input="dtmf" numDigits="1" method="POST" timeout="15" actionOnEmptyResult="true">
						<Say voice="Polly.Joanna" language="en-US">
							This is Caller. There is an ongoing incident and you are on-call.
							Please confirm by pressing 5
						</Say>
					</Gather>
				</Response> 
				`
		return c.Blob(200, "application/xml", []byte(resp))
	}

	if callStatus == "in-progress" {
		var resp string
		switch digits {
		case "9":
			resp = `<?xml version="1.0" encoding="UTF-8"?>
				<Response>
					<Say voice="Polly.Joanna" language="en-US">
						You have been recorded as declining to handle the incident, next person will be called.
					</Say>
					<Hangup/>
				</Response>`

			i, err := db.GetIncident()
			if err != nil {
				fmt.Println("[Error]", err)
			}

			p, err := db.GetPersonByPhone(called)
			if err != nil {
				fmt.Println("[Error]", err)
			}
			i.Declined = append(i.Declined, p)
			err = db.WriteIncident(i)
			if err != nil {
				fmt.Println("[Error]", err)
			}

			go func() {
				time.Sleep(5 * time.Second)
				mkcall(1)
			}()
		case "5":
			resp = `<?xml version="1.0" encoding="UTF-8"?>
				<Response>
					<Say voice="Polly.Joanna" language="en-US">
						Thank you, you have been recorded to handle the incident
					</Say>
					<Hangup/>
				</Response>`

			i, err := db.GetIncident()
			if err != nil {
				fmt.Println("[Error]", err)
			}

			p, err := db.GetPersonByPhone(called)
			if err != nil {
				fmt.Println("[Error]", err)
			}
			i.Responsible = &p
			i.Status = "Claimed"
			err = db.WriteIncident(i)
			if err != nil {
				fmt.Println("[Error]", err)
			}

		default:
			resp = `<?xml version="1.0" encoding="UTF-8"?>
				<Response>
					<Gather input="dtmf" numDigits="1" method="POST" timeout="15"
						actionOnEmptyResult="true"
						>
						<Say voice="Polly.Joanna" language="en-US">
							Please enter 5 to confirm or 9 to have the next person on the call list be called
						</Say>
					</Gather>
				</Response> `
		}

		return c.Blob(200, "application/xml", []byte(resp))

	}

	if callStatus == "completed" {
		i, err := db.GetIncident()
		if err != nil {
			fmt.Println("[Error]", err)
			return err
		}

		if i.Status != "Claimed" {
			fmt.Println("------------- Hangup by " + called + " --------------")

			i, err := db.GetIncident()
			if err != nil {
				fmt.Println("[Error]", err)
			}

			p, err := db.GetPersonByPhone(called)
			if err != nil {
				fmt.Println("[Error]", err)
			}
			i.Declined = append(i.Declined, p)
			err = db.WriteIncident(i)
			if err != nil {
				fmt.Println("[Error]", err)
			}

			go func() {
				time.Sleep(5 * time.Second)
				mkcall(1)
			}()
		}

		return nil
	}

	return nil
}
