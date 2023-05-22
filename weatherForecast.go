package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const ()

var (
	tokens []string
)

// MODELS ------------------------------------------------------------------------------------

type LoginRequest struct {
	Login    string
	Password string
}

type LoginResponse struct {
	Token string
}

type Geo struct {
	Lat  int
	Long int
}

type Weather struct {
	Weather     string
	Temperature int
	Humidity    int
}

type Request struct {
	Token string
	Geo   Geo
	Unit  string
}

type Response struct {
	Weather Weather
}

// FUNC --------------------------------------------------------------------------------------

func main() {
	handleRequests()
}

// router
func handleRequests() {
	router := mux.NewRouter()

	router.HandleFunc("/api/test/{data}", test).Methods("GET")            // test
	router.HandleFunc("/api/weather/login", loginAttempt).Methods("POST") // login
	router.HandleFunc("/api/weather", weatherForecast).Methods("POST")    // weather forecast

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowCredentials: true,
	})
	handler := cors.Handler(router)
	log.Fatal(http.ListenAndServe(":10000", handler))
}

// ENDPOINTS ---------------------------------------------------------------------------------

// GET /api/test/{data}
func test(w http.ResponseWriter, r *http.Request) {
	var data string

	vars := mux.Vars(r)
	data = vars["data"]

	if len(data) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if len(data) > 255 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	w.WriteHeader(http.StatusOK)
	errEncode := json.NewEncoder(w).Encode(data)
	if errEncode != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func loginAttempt(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest
	var token string
	var response LoginResponse

	request, errIO := ioutil.ReadAll(r.Body)
	if errIO != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	errUnmarshal := json.Unmarshal(request, &payload)
	if errUnmarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if payload.Login == "user" && payload.Password == "password" {
		token = uuid.New().String()
		tokens = append(tokens, token)
		response = LoginResponse{Token: token}
		w.WriteHeader(http.StatusOK)

		errEncode := json.NewEncoder(w).Encode(response)
		if errEncode != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

}

func weatherForecast(w http.ResponseWriter, r *http.Request) {
	var payload Request
	var validTokenFlag = false
	var unit, weather string
	var temperature int
	var response Response

	request, errIO := ioutil.ReadAll(r.Body)
	if errIO != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	errUnmarshal := json.Unmarshal(request, &payload)
	if errUnmarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set flag if token is in token list
	for _, item := range tokens {
		if item == payload.Token {
			validTokenFlag = true
			break
		}
	}
	// if token is valid -> check temp unit
	if !validTokenFlag {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		switch payload.Unit {
		case "C":
			unit = "C"
			break
		case "F":
			unit = "F"
			break
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	rand.Seed(time.Now().UnixNano())
	// set random weather
	weatherList := [7]string{"sunny", "sunny", "clear sky", "rain", "cloudy", "thunderstorm", "snow"}
	weather = weatherList[rand.Intn((len(weatherList)))]
	// set random temperature depending on unit
	switch unit {
	case "C":
		temperature = rand.Intn(20 - -20) + -20
		break
	case "F":
		temperature = rand.Intn(64 - -4) + -4
		break
	}
	// create&send response
	response = Response{Weather{Weather: weather, Temperature: temperature, Humidity: rand.Intn(80)}}
	w.WriteHeader(http.StatusOK)

	errEncode := json.NewEncoder(w).Encode(response)
	if errEncode != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
