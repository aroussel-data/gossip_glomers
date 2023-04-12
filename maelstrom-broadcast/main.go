package main

import (
	"encoding/json"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	var intStore map[int]bool
	n := maelstrom.NewNode()

	// Register a handler for the "broadcast" message.
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if res, ok := body["message"].(int); ok {
			intStore[res] = true
		}

		newBody := make(map[string]any)

		if src_val, ok := body["src"].(string); ok {
			newBody["dest"] = src_val
		}

		// for k, v := range body {
		// 	if k != "message" {
		// 		newBody[k] = v
		// 	}
		// }

		// Update the message type.
		newBody["type"] = "broadcast_ok"

		// Respond with the updated message type.
		// return n.Reply(msg, newBody)
		return n.Send(newBody["dest"], newBody["type"])
	})

	// Register a handler for the "read" message.
	n.Handle("read", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		var keys []int
		for k := range intStore {
			keys = append(keys, k)
		}

		// Update the message type.
		body["type"] = "read_ok"
		body["messages"] = keys

		// Respond with the updated message type.
		return n.Reply(msg, body)
	})

	// Register a handler for the "topology" message.
	n.Handle("topology", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		newBody := make(map[string]any)
		for k, v := range body {
			if k != "topology" {
				newBody[k] = v
			}
		}

		// Update the message type.
		newBody["type"] = "topology_ok"

		// Respond with the updated message type.
		return n.Reply(msg, newBody)
	})

	// Execute the node's message loop. This will run until STDIN is closed.
	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}

}
