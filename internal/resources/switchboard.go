package resources

import (
	"callr/internal/config"
	"callr/internal/dao"
	"fmt"
	"github.com/labstack/echo"
	"github.com/sfreiberg/gotwilio"
	"io/ioutil"
	"time"
)

func rotIncident(db dao.Dao, cfg config.Config) error {
	i, err := db.GetIncident()

	if err != nil && err.Error() == "no incident exists" {
		return nil
	}
	if err != nil {
		return err
	}

	if time.Now().After(i.CreatedAt.Add(config.Get().IncidentRottenDuration)) {
		_, err = db.CloseIncident("Rotten")
	}

	return err
}

func Incident(db dao.Dao, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.QueryParams().Get("token") != cfg.IncidentToken {
			return c.String(400, "correct token was not provided")
		}

		err := rotIncident(db, cfg)
		if err != nil {
			return err
		}

		err = db.NewIncident()
		if err == nil { // first once that creates it...
			go caller(db, cfg)
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
}

func caller(db dao.Dao, cfg config.Config) {
	sleep := time.Second * 5
	twilio := gotwilio.NewTwilioClient(cfg.TwilSID, cfg.TwilToken)
	for {
		i, err := db.GetIncident()

		if err != nil && err.Error() == "no incident exists" {
			return
		}

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
			mkcall(0, db, cfg)
		case "Calling":

			if i.OnCallIndex > 10 {
				sleep = time.Minute
				break
			}

			if time.Since(i.LastCall) > time.Minute {
				mkcall(1, db, cfg)
				break
			}

			r, errs, err := twilio.GetCall(i.CallId)
			if err != nil {
				fmt.Println("[Error]", err)
				break
			}
			if errs != nil {
				fmt.Println("[Twil Error]", errs.Error())
				break
			}

			switch r.Status {
			case "busy", "cancelled", "completed", "failed", "no-answer":
				mkcall(1, db, cfg)
			}

		case "Claimed":
			sleep = time.Minute
			err = rotIncident(db, cfg)
			if err != nil {
				fmt.Println("[Error]", err)
			}
		case "Failed":
			sleep = time.Minute
		case "Closed", "Rotten":
			return
		}
		<-time.After(sleep)
	}

}

func mkcall(inc int, db dao.Dao, cfg config.Config) {
	i, err := db.GetIncident()
	oncall, err := db.GetOnCall()
	if err != nil {
		fmt.Println("[Error]", err)
		return
	}
	if len(oncall) == 0 {
		fmt.Println("[Error]", "no one in the oncall list")
		i.Status = "Failed"
		i.Messages = append(i.Messages, "No one on-call")
		err = db.WriteIncident(i)
		if err != nil {
			fmt.Println("[Error]", err)
		}
		return
	}
	if i.OnCallIndex > 10 {
		return
	}

	i.OnCallIndex += inc
	idx := i.OnCallIndex % len(oncall)
	p := oncall[idx]

	twilio := gotwilio.NewTwilioClient(cfg.TwilSID, cfg.TwilToken)
	r, errs, err := twilio.CallWithUrlCallbacks(cfg.TwilPhone, p.Phone, gotwilio.NewCallbackParameters(cfg.BaseURL+"/switchboard/page"))
	if err != nil {
		fmt.Println("[Error]", err)
		return
	}
	if errs != nil {
		fmt.Println("[Twil Error]", errs.Error())
		return
	}

	i.CallId = r.Sid
	i.LastCall = time.Now()
	err = db.WriteIncident(i)
	if err != nil {
		fmt.Println("[Error]", err)
		return
	}

}

func TestCall() echo.HandlerFunc {
	return func(c echo.Context) error {
		resp := `<?xml version="1.0" encoding="UTF-8"?>
				<Response>
						<Say voice="Polly.Joanna" language="en-US">
							This is a Test call from caller
						</Say>
						<Hangup/>
				</Response> 
				`
		return c.Blob(200, "application/xml", []byte(resp))
	}
}

func Page(db dao.Dao, cfg config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {

		params, err := c.FormParams()

		if err != nil {
			fmt.Println(0, err)
			return err
		}

		callId := params.Get("CallSid")
		callStatus := params.Get("CallStatus")
		called := params.Get("Called")
		msg := params.Get("msg")
		digits := params.Get("Digits")

		// Stop calling
		i, err := db.GetIncident()
		if err != nil {
			resp := `<?xml version="1.0" encoding="UTF-8"?>
				<Response>
						<Say voice="Polly.Joanna" language="en-US">
							Something went wrong with callr
						</Say>
						<Hangup/>
				</Response> 
				`
			return c.Blob(200, "application/xml", []byte(resp))
		}
		// Stop calling if for some reason, status is not calling...
		if i.Status != "Calling" || i.CallId != callId {
			resp := `<?xml version="1.0" encoding="UTF-8"?>
				<Response>
						<Hangup/>
				</Response> 
				`
			return c.Blob(200, "application/xml", []byte(resp))
		}

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
					mkcall(1, db, cfg)
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
					mkcall(1, db, cfg)
				}()
			}

			return nil
		}

		return nil
	}
}
