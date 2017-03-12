package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/mediocregopher/radix.v2/pool"

	"github.com/bilinguliar/gokiezen/msg"
	"github.com/bilinguliar/gokiezen/score"
	"github.com/bilinguliar/gokiezen/voting"
)

const (
	candidatesEndpoint = "/candidates"
	statsEndpoint      = "/stats"
	statsWSEndpoint    = "/stats/ws"
	voteEndpoint       = "/track"
	frontend           = "/"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// Flags declared in main to prevent direct usage, they will be passed only as function parameters.
	var (
		port          string
		event         string
		token         string
		redisHost     string
		redisPort     string
		redisConType  string
		redisPoolSize int
	)

	/*
		According to https://12factor.net the best place to store config - environment variables.
		Will use this approach.	Env vars will be declared in Dockerfile, startup script will pass them as execution params.
		Exception is Token. Use secrets to provide this environment variable.
	*/

	flag.StringVar(&port, "port", "8080", "Specifies port that server will use to accept connections")
	flag.StringVar(&event, "event", "WrldDomntn", "Event name. Eurovision for example. 11 symbols max.")
	flag.StringVar(&token, "token", "", "SMS Gateway API token")
	flag.StringVar(&redisHost, "redis_host", "redis", "Redis host")
	flag.StringVar(&redisPort, "redis_port", "6379", "Redis server port")
	flag.StringVar(&redisConType, "redis_conn_type", "tcp", "Redis connetction type")
	flag.IntVar(&redisPoolSize, "redis_pool_size", 10, "Redis pool size")

	flag.Parse()

	scoreKeeper := score.NewKeeper(
		newPool(
			redisHost+":"+redisPort,
			redisConType,
			redisPoolSize,
		),
	)

	msgChan := make(chan msg.Request)

	birdClient := msg.NewMsgBirdClient(
		token,
		msgChan,
	)

	go msg.StartSendingMessages(context.TODO(), msgChan, birdClient)

	votingSvc := voting.New(
		birdClient, // Messenger
		birdClient, // Enquirer
		scoreKeeper,
		event,
	)

	candsSvc := voting.NewCandidates(scoreKeeper)

	ctrl := voting.NewController(votingSvc, candsSvc)

	http.HandleFunc(candidatesEndpoint, ctrl.HandleCandidates) // Add/Delete candidates.
	http.HandleFunc(statsWSEndpoint, ctrl.GetStatsWS)          // Current voting score via WebSocket.
	http.HandleFunc(statsEndpoint, ctrl.GetStats)              // Voting score via REST API.
	http.HandleFunc(voteEndpoint, ctrl.HandleVote)             // Web hook that accepts requests from SMS web service.
	http.HandleFunc(frontend, voting.ServeHTML)                // HTML file handler. Simple page that listens to WebSocket.

	// TODO handle graceful shutdown.
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// newPool inits new Redis pool.
func newPool(url, conType string, poolSize int) *pool.Pool {
	p, err := pool.New(conType, url, poolSize)
	if err != nil {
		log.Fatal("Redis pool init failed, error: ", err)
	}

	return p
}
