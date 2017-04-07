package main

//to re generate the asset you need to install esc once: go get -u github.com/mjibson/esc
//go:generate esc -o assets.go static/

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	// middlewareChain := alice.New(loggingHandler, combineLogHandler, servedByHandler)
	router := mux.NewRouter().StrictSlash(true)

	// router.
	// 	Methods("GET").
	// 	PathPrefix("/updates/").
	// 	Name("Serve static files").
	// 	Handler(middlewareChain.Then(http.FileServer(FS(false))))

	router.HandleFunc("/ipv4", ipv4Handler)
	router.HandleFunc("/metadata", metadataHandler)

	return router
}

func ipv4Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("192.168.2.1\n192.168.2.2"))
}

func metadataHandler(w http.ResponseWriter, r *http.Request) {
	json := `{"random_seed": "3HqKry3SwZ3lEk75lWiu4/wF8JZ9Bzl3qHUh6wOpPC8po2cvV+FgBRAw1e9VeoqfITSwb3Eoxm5Ekpl4QNO0D1yzszAZX3nbYqSREEDagyaX6bnQtTsenKgDxg1EZVAsmQnFJ5wKFR0l2H4NyKxFG3z+tH2YHt4wbb72724hwmFKdVWwJwNLJ2N52qIx8GP24CM2NFwdzcjM/ab4Yd3YQMygUVN18ja9h1/tH2dZRpA+WicoHoD6lLwLl7ZSDyX9L9wPNvbF0u1Wo51aO0nBm7VwEe/+zXxGGntMfrzmJGRCUXVjNeBihN8HrFOpE5odx4aN5dUu/AE70UJTwjo7/atmJNrd4UN1uTlMil1J0lW42/C9mzCsv1D2vjMm5BtVZxSjAJS+mLubmxsv3yFEwa5xUP7dgO2Xy2WOObkiXIy9EfVawvgae3zT6pgYu5iLwvQg1bjlaXN+8Wb1j1Uwax5MX0kwUjysMmwSbtk8hD4bHVD4q1WVh/GrJfciugAxuRBFejtHxUd30JgQdVvTXk7MgF5mUx2y8AMIp+/oc2o8fR0+4DZxmKo/FyV+ww/FeETPub7MLx6I5M1FcnPt5RDITpWiapgzcyJB7Io+Mo1OPRFCKDmrFRQNBKhcjwY3dFC5+flsz6yb5sq8xsIZrbsfED2JltKO7NTYwjpm3nY=", "uuid": "aa50283d-81d3-40f0-8bbd-42fe1751bff0", "availability_zone": "eu-de-1a", "keys": [{"data": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDeevePzO0cd4sV/iIvKsaVw08s8gC0BuzyJu7/hSt10tt4O46nZpqLLybB/K/telZ2vpgvVSxfpOMWngVCIegv72jVnEMpD1WIL40ackN7TRRfT6JLrSYwvoUYgE3CIk8TBZYf9OalXWWXVgYdHp/12u7NMOEvwBEdor9aWKB39ojnuA3s5guZt4fqBuOoaYE/32W4sL4TL7QqLBBqdGKjOZKvVZITKr4IPn4EDUVoGJ2hKS8f89kNSvmDe4tFgWSu7mohc9V7M8N0TkNz9bXIFe+9tZtpM55ZJIhlejvqHgn0yXX/evtvdwjjZv0aShaqDZkWkGfsmjXbMlmWs2+J a.reuschenbach.puncernau@sap.com", "type": "ssh", "name": "Arturo_std"}], "hostname": "rel7.novalocal", "launch_index": 0, "public_keys": {"Arturo_std": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDeevePzO0cd4sV/iIvKsaVw08s8gC0BuzyJu7/hSt10tt4O46nZpqLLybB/K/telZ2vpgvVSxfpOMWngVCIegv72jVnEMpD1WIL40ackN7TRRfT6JLrSYwvoUYgE3CIk8TBZYf9OalXWWXVgYdHp/12u7NMOEvwBEdor9aWKB39ojnuA3s5guZt4fqBuOoaYE/32W4sL4TL7QqLBBqdGKjOZKvVZITKr4IPn4EDUVoGJ2hKS8f89kNSvmDe4tFgWSu7mohc9V7M8N0TkNz9bXIFe+9tZtpM55ZJIhlejvqHgn0yXX/evtvdwjjZv0aShaqDZkWkGfsmjXbMlmWs2+J a.reuschenbach.puncernau@sap.com"}, "project_id": "1d1ad583e98c4913a0226feac0f010f9", "name": "rel7"}`

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(json))
}
