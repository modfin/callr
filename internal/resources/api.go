package resources

import (
	"callr/internal/config"
	"callr/internal/dao"
	"github.com/labstack/echo"
	"github.com/sfreiberg/gotwilio"
)

func DeleteActiveIncident(db dao.Dao) echo.HandlerFunc {
	return func(c echo.Context) error {
		i, err := db.CloseIncident()
		if err != nil {
			return err
		}
		return c.JSON(200, i)
	}
}

func GetActiveIncident(db dao.Dao) echo.HandlerFunc{
	return func(c echo.Context) error {
		is, err := db.GetIncident()
		if err == nil {
			return c.JSON(200, is)

		}
		if err.Error() == "no incident exists" {
			return c.JSON(200, nil)
		}
		return err
	}
}

func GetOldIncidents(db dao.Dao) echo.HandlerFunc{
	return func(c echo.Context) error {
		is, err := db.GetIncidents()
		if err != nil {
			return err
		}
		return c.JSON(200, is)
	}
}

func GetIncidentLogs(db dao.Dao) echo.HandlerFunc{
	return func(c echo.Context) error {
		logs, err := db.GetLogs(c.Param("id"))
		if err != nil {
			return err
		}
		return c.JSON(200, logs)
	}
}


func GetPeople(db dao.Dao) echo.HandlerFunc{
	return func(c echo.Context) error {
		people, err := db.GetPeople()
		if err != nil {
			return err
		}
		return c.JSON(200, people)
	}
}

func PostPeople(db dao.Dao) echo.HandlerFunc{
	return func(c echo.Context) error {

		var people []dao.Person
		err := c.Bind(&people)
		if err != nil{
			return err
		}


		err = db.WritePeople(people)
		if err != nil {
			return err
		}
		return c.JSON(200, people)
	}
}

func PostOnCall(db dao.Dao) echo.HandlerFunc{
	return func(c echo.Context) error {

		var people []dao.Person
		err := c.Bind(&people)
		if err != nil{
			return err
		}


		err = db.WriteOnCall(people)
		if err != nil {
			return err
		}
		return c.JSON(200, people)
	}
}


func GetOnCall(db dao.Dao) echo.HandlerFunc{
	return func(c echo.Context) error {
		people, err := db.GetOnCall()
		if err != nil {
			return err
		}
		return c.JSON(200, people)
	}
}

func GetTestCall(cfg config.Config) echo.HandlerFunc{
	return func(c echo.Context) error {

		phone := c.Param("phone")

		twilio := gotwilio.NewTwilioClient(cfg.TwilSID, cfg.TwilToken)
		_, _, err := twilio.CallWithUrlCallbacks(cfg.TwilPhone, phone, gotwilio.NewCallbackParameters(cfg.BaseURL+"/switchboard/test-call"))
		return err
	}
}