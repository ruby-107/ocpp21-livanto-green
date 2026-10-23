# OCPP 2.1 Charging Station Simulator

## 1. Project Overview

This project is a simple EV Charging Station Simulator developed using **Golang, WebSocket, OCPP 2.1, and PostgreSQL**.

The simulator communicates with the CSMS (Charging Station Management System) through WebSocket and handles three main OCPP messages:

- 🚀 **BootNotification:** Sends charging station information when the station starts.
- 🔋 **StatusNotification:** Sends the connector's current status.
- 💓 **Heartbeat:** Periodically informs the server that the station is still communicating.

PostgreSQL is used to store station details and message history.

### Project Flow

```text
Charging Station Simulator
          |
          | WebSocket + OCPP 2.1
          v
       Go CSMS
          |
          v
      PostgreSQL
```

## 2. Project Structure

```text
ocpp21-livantogreen/
├── server/
│   ├── main.go
│   └── handler.go
├── simulator/
│   └── main.go
├── internal/
│   ├── database/
│   │   ├── postgres.go
│   │   └── repository.go
│   └── ocpp/
│       └── messages.go
├── .env
└── go.mod
```

### File Responsibilities

**`server/main.go`**
- Starts the CSMS HTTP server.
- Connects to PostgreSQL.
- Registers the WebSocket endpoint.

**`server/handler.go`**
- Handles incoming OCPP messages.
- Processes BootNotification, Heartbeat, and StatusNotification.
- Calls database methods to save the data.
- Sends responses to the simulator.

**`simulator/main.go`**
- Connects the simulated charging station to the CSMS.
- Sends BootNotification when starting.
- Sends StatusNotification after startup.
- Sends Heartbeat every 30 seconds.

**`internal/database/postgres.go`**
- Establishes the PostgreSQL connection.

**`internal/database/repository.go`**
- Saves station details.
- Updates the latest Heartbeat timestamp.
- Stores message history.

**`internal/ocpp/messages.go`**
- Builds OCPP CALL and CALLRESULT messages.

**`.env`**
- Contains PostgreSQL connection settings.

## 3. OCPP Messages

### 🚀 BootNotification

BootNotification is sent when the charging station starts and connects to the CSMS.

It contains information such as:
- Station model
- Vendor name
- Boot reason

The server saves the station information in PostgreSQL and responds with an acceptance status and Heartbeat interval.

Example station:
- Station ID: `JP-001`
- Model: `JP-Simulator`
- Vendor: `DemoVendor`

### 🔋 StatusNotification

StatusNotification reports the status of a charging connector.

Example:
- Connector ID: `1`
- Status: `Available`

The server stores the connector status in the message history table.

### 💓 Heartbeat

Heartbeat is sent periodically by the simulator to inform the CSMS that the station is still communicating.

In this project, the Heartbeat interval is 30 seconds.

When a Heartbeat arrives:
- The server updates `last_heartbeat` in `charging_stations`.
- The server adds a Heartbeat record to `station_logs`.
- The server responds with its current time.

## 4. PostgreSQL Tables

The project uses two tables.

### `charging_stations`

Stores the latest information for each charging station.

| Column | Purpose |
|---|---|
| `station_id` | Unique station ID |
| `model` | Charging station model |
| `vendor_name` | Station vendor |
| `last_heartbeat` | Latest Heartbeat time |
| `created_at` | Record creation time |

### `station_logs`

Stores the history of the three OCPP messages.

| Column | Purpose |
|---|---|
| `id` | Unique log ID |
| `station_id` | Charging station ID |
| `message_type` | OCPP message name |
| `status` | Message or connector status |
| `connector_id` | Connector ID, when applicable |
| `created_at` | Log creation time |

## 5. Requirements

Before running the project, make sure you have:

- Go installed
- PostgreSQL installed and running
- A PostgreSQL database named `livantogreen`

## 6. Install Dependencies

Open the terminal in the project root folder and run:

```bash
go get github.com/gorilla/websocket
go get github.com/jackc/pgx/v5/stdlib
go get github.com/joho/godotenv
go mod tidy
```

## 7. Environment Configuration

Create a `.env` file in the project root:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=livantogreen
DB_SSLMODE=disable
```

Replace `your_password` with your PostgreSQL password.

## 8. Run the Project

Open two terminals in the project root directory.

### Terminal 1: Start the CSMS Server

```bash
go run ./server
```

The server listens on port `9000`.

### Terminal 2: Start the Simulator

```bash
go run ./simulator
```

The simulator connects to:

```text
ws://localhost:9000/ocpp/JP-001
```

Both the simulator and server should be configured to use the WebSocket subprotocol `ocpp2.1`.

## 9. Expected Output

The logs should show a flow similar to this:

```text
🔌 Connected to CSMS successfully!

🚀 [BootNotification] Charging Station is starting...
⚡ Station ID: JP-001
📩 [BootNotification] CSMS Response: ...

🔋 [StatusNotification] Connector Status: Available

💓 [Heartbeat] Sending heartbeat to CSMS...
💓 [Heartbeat] CSMS Response: ...
```

Heartbeat messages should repeat approximately every 30 seconds while the simulator is running.

## 10. Verify Data in PostgreSQL

### Check Charging Station Details

```sql
SELECT *
FROM charging_stations;
```

### Check All Message Logs

```sql
SELECT *
FROM station_logs
ORDER BY id DESC;
```

### Check Only the Three OCPP Messages

```sql
SELECT
    station_id,
    message_type,
    status,
    connector_id,
    created_at
FROM station_logs
WHERE message_type IN (
    'BootNotification',
    'Heartbeat',
    'StatusNotification'
)
ORDER BY id DESC;
```

## 11. Troubleshooting

**Database connection failed**
- Verify that PostgreSQL is running.
- Check the values in `.env`.

**WebSocket connection failed**
- Start the CSMS server before the simulator.
- Check the WebSocket URL and `ocpp2.1` subprotocol configuration.

**Heartbeat is not appearing**
- Wait for the 30-second interval.
- Check that the simulator is still running.

**Database table does not exist**
- Verify that the required tables have been created in `livantogreen`.

## 12. Conclusion

This project demonstrates a simple OCPP communication flow using Go, WebSocket, and PostgreSQL.

It covers:
- Charging station registration through BootNotification.
- Connector status reporting through StatusNotification.
- Periodic communication through Heartbeat.
- Saving station details and message history in PostgreSQL.

**Technologies:** Golang, WebSocket, OCPP 2.1, PostgreSQL.
