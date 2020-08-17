package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"bigo.digital/alertmanager-discord/pkg/alertmanager"
	"bigo.digital/alertmanager-discord/pkg/discord"
	"github.com/gorilla/mux"
	"github.com/plaid/go-envvar/envvar"
)

type envVars struct {
	Port       int    `envvar:"PORT" default:"9094"`
	Address    string `envvar:"ADDRESS" default:"0.0.0.0"`
	WebhookURL string `envvar:"DISCORD_WEBHOOK"`
}

var (
	vars   envVars
	logger *log.Logger
)

func sendMessage(payload alertmanager.Payload) {
	message := discord.MakeMessage(payload)
	encoded, err := json.Marshal(message)

	if err != nil {
		logger.Panicf("%+v\n", err)
		return
	}

	_, err = http.Post(vars.WebhookURL, "application/json", bytes.NewBuffer(encoded))

	if err != nil {
		logger.Panicf("%+v\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("%+v\n", *r)
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Panicf("%+v\n", err)
		return
	}

	payload := alertmanager.Payload{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		logger.Panicf("%+v\n", err)
		return
	}

	logger.Printf("Received alert: %+v\n", payload)
	sendMessage(payload)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	if err := envvar.Parse(&vars); err != nil {
		logger.Fatal(err)
	}

	mux := mux.NewRouter()
	mux.HandleFunc("/", handler).Methods("POST")
	mux.HandleFunc("/", healthCheck).Methods("GET")

	address := fmt.Sprintf("%s:%d", vars.Address, vars.Port)
	logger.Printf("Listening on: %s", address)
	logger.Fatalln(http.ListenAndServe(address, mux))
}
