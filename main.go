package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
)

// Config represents the information in the config file
type Config struct {
	SerialPort string
	LivePin    string
	NeutralPin string
	ServoPin   string
	Port       string
}

func main() {
	config := readConfig("./config.toml")
	log.Print("Connecting to ", config.SerialPort)
	plug, err := NewPlug(config.SerialPort, config.LivePin, config.NeutralPin, config.ServoPin)
	defer plug.Disconnect()
	if err != nil {
		log.Print(err.Error())
		return
	}
	plug.Off()
	defer plug.Off()
	power := false

	m := mux.NewRouter()
	m.HandleFunc("/on", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Fprint(w, "Please issue a POST request")
			return
		}
		power = true
		plug.On()
		fmt.Fprint(w, "OK\n")
	})

	m.HandleFunc("/off", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Fprint(w, "Please issue a POST request")
			return
		}
		power = false
		plug.Off()
		fmt.Fprint(w, "OK\n")
	})

	m.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			fmt.Fprintf(w, "Please issue a GET request")
			return
		}
		var response string
		if power {
			response = "ON\n"
		} else {
			response = "OFF\n"
		}
		fmt.Fprint(w, response)
	})

	m.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Fprintf(w, "Please issue a POST request")
			return
		}
		plug.Input()
		fmt.Fprint(w, "OK\n")
	})

	http.Handle("/", m)
	log.Print("Listening on ", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil))
}

func readConfig(configfile string) Config {
	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
