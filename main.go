package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
)

type Message struct {
	Action string
	Data   string
}

//HashData confirms build and source files
type HashData struct {
	Source string `json:"source"`
	Build  string `json:"build"`
}

//RequirementsData outlines required running ssystems.
type RequirementsData struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

//APIData describes inputs and outputs of ssystem.
type APIData struct {
	Read  []string `json:"read"`
	Write []string `json:"write"`
}

type Manager struct {
	URL        string `json:"url"`
	Connection net.Conn
}

type Repository struct {
	URL string `json:"url"`
}

type Body struct {
	Name         string             `json:"name"`
	Version      string             `json:"version"`
	Repository   Repository         `json:"repository"`
	Hash         HashData           `json:"hash"`
	API          APIData            `json:"api"`
	Requirements []RequirementsData `json:"requirements"`
	Limit        int                `json:"limit"`
	Manager      Manager            `json:"manager"`
}

type SubSystem struct {
	Token    string
	Body     Body
	Inbound  chan Message
	Outbound chan Message
}

func runCommand(args []string) {
	cmd := exec.Command("cmd", args...)
	cmd.Run()
}

func (s *SubSystem) Deploy(url string, name string) {
	runCommand([]string{"/k", "go", "get", url})
	go runCommand([]string{"/k", name})
}

func (s *SubSystem) Connect() error {
	conn, err := net.Dial("tcp", s.Body.Manager.URL)

	s.Body.Manager.Connection = conn

	return err
}

func (s *SubSystem) Listen() {
	enc := json.NewEncoder(s.Body.Manager.Connection)
	enc.Encode(Message{Action: "Token", Data: "Test"})

	go func() {
		for {
			select {
			case msg := <-s.Outbound:
				go func() {
					log.Println(msg)
					enc.Encode(Message{Action: msg.Action, Data: msg.Data})
				}()
			case msg := <-s.Inbound:
				log.Println(msg)
			}
		}
	}()
}

func (s *SubSystem) loadConfig(url string) {
	file, _ := ioutil.ReadFile("./subsystem.json")

	d := Body{}
	json.Unmarshal(file, &d)
	d.Manager.URL = url
	s.Body = d
}

func New(url string) *SubSystem {
	s := new(SubSystem)
	s.loadConfig(url)

	if err := s.Connect(); err != nil {
		log.Fatalln("Failed to connect to Manager")
	}

	return s
}

func main() {

	s := New("localhost:8081")

	enc := json.NewEncoder(s.Body.Manager.Connection)
	enc.Encode(Message{Action: "Token", Data: "Test"})
	s.Body.Manager.Connection.Close()

	//enc.Encode(Message{Action: "Close"})
}
