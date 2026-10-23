package main

import (
	"log"
	"strconv"
	"time"

	"ocpp21-livantogreen/internal/ocpp"

	"github.com/gorilla/websocket"
)

func sendCall(
	conn *websocket.Conn,
	id, action string,
	payload any,
) error {
	message, err := ocpp.BuildCall(id, action, payload)
	if err != nil {
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	_, response, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	if action == "Heartbeat" {
		log.Println("💓 [Heartbeat] Sending heartbeat to CSMS...")
		log.Printf("💓 [Heartbeat] CSMS Response: %s", response)
	} else {
		log.Printf("📩 [%s] CSMS Response: %s", action, response)
	}

	return nil
}

func main() {
	url := "ws://localhost:9000/ocpp/JP-001"

	dialer := websocket.Dialer{
		Subprotocols: []string{"ocpp2.1"},
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer conn.Close()

	log.Println("🔌 Connected to CSMS successfully!")
	log.Println("🚀 [BootNotification] Charging Station is starting...")
	log.Println("⚡ Station ID: JP-001")
	log.Println("✅ [BootNotification] Station registered successfully!")

	// 1. BootNotification
	boot := map[string]any{
		"chargingStation": map[string]string{
			"model":      "JP-Simulator",
			"vendorName": "DemoVendor",
		},
		"reason": "PowerUp",
	}

	if err := sendCall(conn, "1", "BootNotification", boot); err != nil {
		log.Fatal(err)
	}

	// 2. StatusNotification
	status := map[string]any{
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
		"connectorStatus": "Available",
		"evseId":          1,
		"connectorId":     1,
	}

	if err := sendCall(conn, "2", "StatusNotification", status); err != nil {
		log.Fatal(err)
	}

	// 3. Heartbeat every 30 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	counter := 3
	for range ticker.C {
		id := strconv.Itoa(counter)

		if err := sendCall(conn, id, "Heartbeat", map[string]any{}); err != nil {
			log.Println("Heartbeat error:", err)
			return
		}

		counter++
	}

	log.Println("🔋 [StatusNotification] Connector Status: Available")
}
