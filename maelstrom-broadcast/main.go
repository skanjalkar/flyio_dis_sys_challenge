package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	var store []string
	var slice []int // Declaration of an integer slice

	// Register a handler for the "echo" message that responds with an "echo_ok".
	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		flag := false
		// Update the message type.
		value := int(body["message"].(float64))
		for i := 0; i < len(slice); i++ {
			if (slice[i] == value) {
				flag = true
			}
		}
		if (flag) {
			return nil
		}
		for i := 0; i < len(store); i++ {
			n.Send(store[i], body)
		}
		// os.Stderr.WriteString(fmt.Sprintln(value))
		{
			slice = append(slice, value)
			// os.Stderr.WriteString("ok")
			// os.Stderr.WriteString(fmt.Sprintf("%v", slice))
		}
		body["type"] = "broadcast_ok"

		delete(body, "message")
		// os.Stderr.WriteString("Counnt of ID" + n.ID())

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type.
		body["type"] = "read_ok"
		body["messages"] = slice
		// os.Stderr.WriteString("Counnt of ID" + n.ID())

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		// Unmarshal the message body as an loosely-typed map.
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type.
		body["type"] = "topology_ok"
		topologyData := body["topology"].(map[string]interface{})
		os.Stderr.WriteString(fmt.Sprintf("%v", topologyData)+"\n")
		os.Stderr.WriteString("\n" + n.ID() + fmt.Sprintf("%v", topologyData[n.ID()]) + "\n")
		if neighbors, ok := topologyData[n.ID()]; ok {
			// os.Stderr.WriteString("got inside if")
			neighborsSlice := neighbors.([]interface{})
			store = make([]string, len(neighborsSlice))

			// Convert the neighbors to a string slice
			for i, neighbor := range neighborsSlice {
				store[i] = neighbor.(string)
			}
		}
		os.Stderr.WriteString("Store here: " + fmt.Sprintf("%d", len(store)) + "\n")

		delete(body, "topology")

		// Echo the original message back with the updated message type.
		return n.Reply(msg, body)
	})

	// Execute the node's message loop. This will run until STDIN is closed.
	if err := n.Run(); err != nil {
		log.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
