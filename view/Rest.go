package view

import (
	"bloom/bloom"
	"encoding/json"
	"log"
	"net/http"
)

type RestInputAdd struct {
	Node string `json:"node"`
}

type RestInputExists struct {
	Node string `json:"node"`
}

type RestDaemonInput struct {
	// Addr of the server e.g. ":8080"
	Addr string
}

func GetViewInputRest(daemonInput RestDaemonInput) (input *NodeInput) {
	input = &NodeInput{
		Id: "rest",
		Setup: func(bloomNode *bloom.Node) error {
			return restDaemon(bloomNode, daemonInput)
		},
		Shutdown: func() error {
			return nil
		},
	}
	return
}

func restDaemon(bloomNode *bloom.Node, daemonInput RestDaemonInput) (err error) {
	http.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {

		if request.Method != "POST" {
			log.Println("invalid http method=" + request.Method)
			writer.WriteHeader(404)
			_, writeErr := writer.Write([]byte("failed to parse input params"))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			return
		}

		var input RestInputAdd
		decodeErr := json.NewDecoder(request.Body).Decode(&input)
		if decodeErr != nil {
			log.Println("failed to parse input", decodeErr)
			writer.WriteHeader(400)
			_, writeErr := writer.Write([]byte("failed to parse input params"))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			return
		}

		if len(input.Node) == 0 {
			writer.WriteHeader(400)
			_, writeErr := writer.Write([]byte("node to add must not be empty"))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			return
		}

		addErr := bloomNode.ItemAdd([]byte(input.Node))
		if addErr != nil {
			log.Println("failed to add item", addErr)
			writer.WriteHeader(400)
			_, writeErr := writer.Write([]byte("an error occurred"))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			return
		}

		writer.WriteHeader(200)
		_, writeErr := writer.Write([]byte("successfully added node"))
		if writeErr != nil {
			log.Fatal(writeErr)
		}
	})

	http.HandleFunc("/exists", func(writer http.ResponseWriter, request *http.Request) {

		if request.Method != "POST" {
			log.Println("invalid http method=" + request.Method)
			writer.WriteHeader(404)
			_, writeErr := writer.Write([]byte("failed to parse input params"))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			return
		}

		var input RestInputExists
		decodeErr := json.NewDecoder(request.Body).Decode(&input)
		if decodeErr != nil {
			log.Println("failed to parse input", decodeErr)
			writer.WriteHeader(400)
			_, writeErr := writer.Write([]byte("failed to parse input params"))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			return
		}

		if len(input.Node) == 0 {
			writer.WriteHeader(400)
			_, writeErr := writer.Write([]byte("node to add must not be empty"))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			return
		}

		contains, containsErr := bloomNode.ItemPossiblyContains([]byte(input.Node))
		if containsErr != nil {
			writer.WriteHeader(400)
			_, writeErr := writer.Write([]byte("an error occurred"))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			return
		}

		if !contains {
			writer.WriteHeader(200)
			_, writeErr := writer.Write([]byte("node does not exists"))
			if writeErr != nil {
				log.Fatal(writeErr)
			}
			return
		}

		writer.WriteHeader(200)
		_, writeErr := writer.Write([]byte("node exists"))
		if writeErr != nil {
			log.Fatal(writeErr)
		}
	})

	err = http.ListenAndServe(daemonInput.Addr, nil)
	if err != nil {
		return
	}

	return
}
