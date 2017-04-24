package SubSystem

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/subsystemio/shared"
)

//HashData confirms build and source files
type HashData struct {
	Source string `json:"source"`
	Build  string `json:"build"`
}

//RequirementsData outlines required running subsystems.
type RequirementsData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

//APIData describes inputs and outputs of subsystem.
type APIData struct {
	Read  []string `json:"read"`
	Write []string `json:"write"`
}

type Connection struct {
	URL string `json:"url"`
}

type Repository struct {
	URL string `json:"url"`
}

type SubSystemData struct {
	Name         string             `json:"name"`
	Version      string             `json:"version"`
	Repository   Repository         `json:"repository"`
	Hash         HashData           `json:"hash"`
	API          APIData            `json:"api"`
	Requirements []RequirementsData `json:"requirements"`
	Limit        int                `json:"limit"`
	Connection   Connection         `json:"connection"`
}

//SubSystem main structure.
type SubSystem struct {
	Data   SubSystemData
	Health HealthCheck
}

type HealthCheck struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (s *SubSystem) loadConfig() {
	file, _ := ioutil.ReadFile("./subsystem.json")

	d := SubSystemData{}
	json.Unmarshal(file, &d)
	s.Data = d
}

func (s *SubSystem) Register(url string) {
	res, err := http.Post(url+"/subsystems", "application/json; charset=utf-8", Shared.ToJSON(s.Data))
	if err != nil {
		log.Fatalf("Failed to register with Manager - %v", err)
	}

	body := Shared.ReadBody(res.Body)

	if res.StatusCode == 500 {
		log.Fatalf("Failed to register with Manager [%v] - %s", res.StatusCode, body)

	}

	log.Printf("Registered with Manager")

	s.Data.Connection.URL = url

	json.Unmarshal(body, &s.Health)
}

func (s *SubSystem) SendStatus() {
	res, err := http.Post(s.Data.Connection.URL+"/health", "application/json; charset=utf-8", Shared.ToJSON(s.Health))
	if err != nil {
		log.Fatalf("Failed to check-in with Manager - %v", err)
	}
	body := Shared.ReadBody(res.Body)
	log.Printf("Ba-bump - %v", string(body))
}

func (s *SubSystem) StartHeartbeat() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			s.SendStatus()
		}
	}
}

func New() SubSystem {
	s := new(SubSystem)
	s.loadConfig()

	return *s
}

func Parse(data []byte) SubSystem {
	d := SubSystemData{}
	json.Unmarshal(data, &d)

	s := SubSystem{Data: d}
	return s
}
