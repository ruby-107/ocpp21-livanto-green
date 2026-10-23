package ocpp

import "encoding/json"

func BuildCall(id, action string, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return json.Marshal([]any{2, id, action, json.RawMessage(body)})
}

func BuildCallResult(id string, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return json.Marshal([]any{3, id, json.RawMessage(body)})
}
