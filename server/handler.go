package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"ocpp21-livantogreen/internal/database"
	"ocpp21-livantogreen/internal/ocpp"

	"github.com/gorilla/websocket"
)

type Server struct {
	Repo *database.Repository
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Local development only
	},
	Subprotocols: []string{"ocpp2.1"},
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	stationID := strings.TrimPrefix(r.URL.Path, "/ocpp/")
	if stationID == "" || stationID == r.URL.Path {
		http.Error(w, "Station ID missing", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket error:", err)
		return
	}
	defer conn.Close()

	log.Println("Station connected:", stationID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Station disconnected:", stationID)
			return
		}

		if err := s.handleCall(conn, stationID, message); err != nil {
			log.Println("Message error:", err)
		}
	}
}

func (s *Server) handleCall(
	conn *websocket.Conn,
	stationID string,
	message []byte,
) error {
	var frame []json.RawMessage

	if err := json.Unmarshal(message, &frame); err != nil {
		return err
	}
	if len(frame) != 4 {
		return nil
	}

	var messageType int
	var messageID, action string

	json.Unmarshal(frame[0], &messageType)
	json.Unmarshal(frame[1], &messageID)
	json.Unmarshal(frame[2], &action)

	if messageType != 2 {
		return nil
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(frame[3], &payload); err != nil {
		return err
	}

	ctx := context.Background()
	response := map[string]any{}

	switch action {
	case "BootNotification":
		var station struct {
			ChargingStation struct {
				Model      string `json:"model"`
				VendorName string `json:"vendorName"`
			} `json:"chargingStation"`
		}

		if err := json.Unmarshal(frame[3], &station); err != nil {
			return err
		}

		if station.ChargingStation.Model == "" ||
			station.ChargingStation.VendorName == "" {
			return s.sendResult(conn, messageID, map[string]any{
				"currentTime": time.Now().UTC().Format(time.RFC3339),
				"interval":    30,
				"status":      "Rejected",
			})
		}

		if err := s.Repo.SaveStation(
			ctx,
			stationID,
			station.ChargingStation.Model,
			station.ChargingStation.VendorName,
		); err != nil {
			return err
		}

		log.Println("BootNotification saved:", stationID)

		response = map[string]any{
			"currentTime": time.Now().UTC().Format(time.RFC3339),
			"interval":    30,
			"status":      "Accepted",
		}

	case "Heartbeat":
		if err := s.Repo.SaveHeartbeat(ctx, stationID); err != nil {
			return err
		}

		log.Println("💓 Heartbeat received:", stationID)

		response = map[string]any{
			"currentTime": time.Now().UTC().Format(time.RFC3339),
		}

	case "StatusNotification":
		var status struct {
			ConnectorStatus string `json:"connectorStatus"`
			ConnectorID     int    `json:"connectorId"`
		}

		if err := json.Unmarshal(frame[3], &status); err != nil {
			return err
		}

		if status.ConnectorStatus == "" {
			return nil
		}

		if err := s.Repo.SaveStatus(
			ctx, stationID, status.ConnectorStatus, status.ConnectorID,
		); err != nil {
			return err
		}

		log.Println("StatusNotification saved:", status.ConnectorStatus)
		response = map[string]any{}

	default:
		log.Println("Unsupported OCPP message:", action)
		return nil
	}

	return s.sendResult(conn, messageID, response)
}

func (s *Server) sendResult(
	conn *websocket.Conn,
	messageID string,
	response any,
) error {
	message, err := ocpp.BuildCallResult(messageID, response)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, message)
}
