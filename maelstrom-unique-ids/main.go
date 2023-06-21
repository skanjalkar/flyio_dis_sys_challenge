package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	var count int = 1

	// Register a handler for the "echo" message that responds with an "echo_ok".
	n.Handle("generate", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type.
		body["type"] = "generate_ok"
		body["id"] = strconv.Itoa(count)+n.ID()
		// os.Stderr.WriteString("Counnt of ID" + n.ID())
		count++

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})


	// Execute the node's message loop. This will run until STDIN is closed.
	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
