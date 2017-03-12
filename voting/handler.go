package voting

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

const (
	statsHTML    = "/opt/gokiezen/stats.html"
	updatePeriod = 1 // WebSocket send period in seconds.
)

// Votes is capable of processing vote SMS messages.
type Votes interface {
	GetStats() (Stats, error)
	RegisterVote(msisdn, text string) error
}

// Candidates can add and delete candidates.
type Candidates interface {
	Add(name string) error
	Del(name string) error
}

type Message struct {
	Originator string
	Body       string
}

// Controller is responsible for requests parsing and responses serialization.
type Controller struct {
	voteSvc   Votes
	candsSvc  Candidates
	statsChan chan Stats
}

// NewController is a constructor for Controller instance.
func NewController(v Votes, c Candidates) *Controller {
	ctrl := &Controller{
		voteSvc:   v,
		candsSvc:  c,
		statsChan: make(chan Stats),
	}

	go ctrl.sendUpdates()

	return ctrl
}

func (c *Controller) HandleCandidates(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		p := req.FormValue("name")
		if p == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "name parameter must be provided")
			return
		}
		err := c.candsSvc.Add(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case "DELETE":
		p := req.FormValue("name")
		if p == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "name parameter must be provided")
			return
		}
		err := c.candsSvc.Del(p)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetStats returns statistics with current voting data.
func (c *Controller) GetStats(w http.ResponseWriter, req *http.Request) {
	// Fail early.
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	stats, err := c.voteSvc.GetStats()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		log.Println("Failed to serialize stats response, error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetStatsWS returns statistics with current voting data via WebSocket.
func (c *Controller) GetStatsWS(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("Failed upgrading to WebSocket, error:", err)
		return
	}
	defer conn.Close()

	var update Stats

	for {
		update = <-c.statsChan

		err = conn.WriteJSON(update)
		if err != nil {
			log.Println("Write to WebSocket failed. Dropping connection. Error:", err)
			break
		}
	}
}

// HandleVote accepts requests with SMS data and passes this data to service responcible for processing.
func (c *Controller) HandleVote(w http.ResponseWriter, req *http.Request) {
	// Expecting only POST from messaging service.
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	err := json.NewDecoder(req.Body).Decode(&msg)
	if err != nil {
		log.Println("Request body is not valid, error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	err = c.voteSvc.RegisterVote(msg.Originator, msg.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// ServeHTML serves single HTML file that displays WebSocket data.
func ServeHTML(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, req, statsHTML)
}

func (c *Controller) sendUpdates() {
	t := time.Tick(updatePeriod * time.Second)

	for now := range t {
		stats, err := c.voteSvc.GetStats()
		if err != nil {
			log.Println("Update was not propagated, time: %v, error: %q", now, err)
			continue
		}

		c.statsChan <- stats
	}
}
