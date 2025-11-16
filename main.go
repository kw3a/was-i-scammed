package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Response struct {
	WasIScammed bool `json:"was_i_scammed"`
}

type RapidApiHeaders struct {
	ProxySecret   string
	User          string
	Subscription  string
	Version       string
	Host          string
	ForwardedFor  string
	ForwardedHost string
}

func headers(r *http.Request) (RapidApiHeaders, error) {
	secret := r.Header.Get("X-RapidAPI-Proxy-Secret")
	user := r.Header.Get("X-RapidAPI-User")
	subscription := r.Header.Get("X-RapidAPI-Subscription")
	version := r.Header.Get("X-RapidAPI-Version")
	host := r.Header.Get("X-RapidAPI-host")
	forwardedFor := r.Header.Get("X-Forwarded-For")
	forwardedHost := r.Header.Get("X-Forwarded-Host")
	if secret == "" {
		return RapidApiHeaders{}, fmt.Errorf("rapidapi secret is missing")
	}
	if secret != os.Getenv("RAPIDAPI_SECRET") {
		return RapidApiHeaders{}, fmt.Errorf("wrong rapidapi secret")
	}
	return RapidApiHeaders{
		ProxySecret: secret,
		User: user,
		Subscription: subscription,
		Version: version,
		Host: host,
		ForwardedFor: forwardedFor,
		ForwardedHost: forwardedHost,
	}, nil
}

func main() {
	if os.Getenv("RAPIDAPI_SECRET") == "" {
		println("RAPIDAPI SECRET is not set")
		os.Exit(0)
	}
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		rapidapiH, err := headers(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		resp := Response{}
		plan := rapidapiH.Subscription
		if plan == "BASIC" {
			resp.WasIScammed = false
		}
		if plan == "PRO" || plan == "ULTRA" || plan == "MEGA" {
			resp.WasIScammed = true
		}
		json.NewEncoder(w).Encode(resp)
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	fmt.Println("Running in the 8080 port")
	http.ListenAndServe(":8080", nil)
}
