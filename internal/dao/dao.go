package dao

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Dao struct {
	mu    sync.Mutex
	Store string
}

func (d Dao) GetOnCall() (p []Person, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	filename := d.Store + "/oncall.json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &p)
	return
}

func (d Dao) NewIncident() error {
	d.mu.Lock()
	filename := d.Store + "/incident.json"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		now := time.Now()
		i := Incident{
			Id:        now.Format(time.RFC3339),
			CreatedAt: now,
			Status:    "Init",
		}

		d.mu.Unlock()
		return d.WriteIncident(i)
	}
	d.mu.Unlock()
	return errors.New("incident already exist")
}

func (d Dao) GetIncident() (Incident, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	filename := d.Store + "/incident.json"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return Incident{}, errors.New("no incident exists")
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Incident{}, err
	}
	i := Incident{}
	err = json.Unmarshal(data, &i)
	return i, err
}

func (d Dao) WriteIncident(i Incident) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	filename := d.Store + "/incident.json"

	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (d Dao) CloseIncident(status ... string) (Incident, error) {
	i, err := d.GetIncident()
	if err != nil {
		return Incident{}, err
	}
	i.Status = "Closed"

	if status != nil{
		i.Status = strings.Join(status, ",")
	}

	err = d.WriteIncident(i)
	if err != nil {
		return Incident{}, err
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	filename := d.Store + "/incident.json"
	toFilename := d.Store + "/incidents/" + i.Id + "/incident.json"

	err = os.MkdirAll(d.Store+"/incidents/"+i.Id, 0770)
	if err != nil {
		return Incident{}, err
	}
	return i, os.Rename(filename, toFilename)
}

func (d Dao) AddLog(incident Incident, log Log) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	err := os.MkdirAll(d.Store+"/incidents/"+incident.Id, 0770)
	if err != nil {
		return err
	}

	filename := d.Store + "/incidents/" + incident.Id + "/" + time.Now().Format(time.RFC3339) + ".log.json"

	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0660)
	return err
}

func (d Dao) GetPeople() ([]Person, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	filename := d.Store + "/people.json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var people []Person
	err = json.Unmarshal(data, &people)
	return people, err
}

func (d Dao) WritePeople(people []Person) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	data, err := json.MarshalIndent(people, "", "  ")
	if err != nil {
		return  err
	}

	return ioutil.WriteFile(d.Store + "/people.json", data, 0660)
}

func (d Dao) WriteOnCall(people []Person) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	data, err := json.MarshalIndent(people, "", "  ")
	if err != nil {
		return  err
	}

	return ioutil.WriteFile(d.Store + "/oncall.json", data, 0660)
}

func (d Dao) GetPersonByPhone(phone string) (Person, error) {
	people, err := d.GetPeople()
	if err != nil {
		return Person{}, err
	}
	for _, p := range people {
		if p.Phone == phone {
			return p, nil
		}
	}
	return Person{}, errors.New("could not find person")
}

func (d Dao) GetLogs(id string) ([]Log, error) {

	ff, err := ioutil.ReadDir(d.Store + "/incidents/" + id)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, f := range ff {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".log.json") {
			files = append(files, f.Name())
		}
	}
	sort.Strings(files)

	var logs []Log
	for _, f := range files {
		filename := d.Store + "/incidents/" + id + "/" + f
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			continue
		}
		var log Log
		err = json.Unmarshal(data, &log)
		if err != nil {
			continue
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (d Dao) GetIncidents() ([]Incident, error) {


	filename := d.Store + "/incidents"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil
	}

	files, err := ioutil.ReadDir(filename)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	if len(dirs) == 0 {
		return nil, nil
	}
	sort.Strings(dirs)
	for i, j := 0, len(dirs)-1; i < j; i, j = i+1, j-1 {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	}

	var incidents []Incident

	for _, di := range dirs {

		filename := d.Store + "/incidents/" + di + "/incident.json"

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			continue
		}
		i := Incident{}
		err = json.Unmarshal(data, &i)
		if err != nil {
			continue
		}
		incidents = append(incidents, i)

	}
	return incidents, nil
}
