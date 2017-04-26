package SubSystem

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
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

type Manager struct {
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
	Manager      Manager            `json:"manager"`
}

//SubSystem main structure.
type SubSystem struct {
	Data SubSystemData
	Port int
}

func (s *SubSystem) loadConfig() {
	file, _ := ioutil.ReadFile("./subsystem.json")

	d := SubSystemData{}
	json.Unmarshal(file, &d)
	s.Data = d
}

func (s *SubSystem) Register(url string) (int, error) {
	res, err := http.Post(url+"/subsystems", "application/json; charset=utf-8", Shared.ToJSON(s.Data))
	if err != nil {
		log.Fatalf("Failed to register with Manager - %v\n", err)
	}

	log.Printf("Registered with Manager: %v\n", res.Body)
	s.Port = 9000

	return 0, err
}

func (s *SubSystem) HealthCheck() error {
	url := fmt.Sprintf("http://localhost:%v/health", s.Port)
	_, err := http.Head(url)
	if err != nil {
		log.Printf("Healthcheck failed. %v is no more.\n", s.Data.Name)
	}

	return err
}

func (s *SubSystem) Serve() {
	r := gin.Default()
	r.HEAD("/health", func(c *gin.Context) {
		c.String(200, "Alive")
	})
	url := fmt.Sprintf(":%v", s.Port)
	r.Run(url)
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
