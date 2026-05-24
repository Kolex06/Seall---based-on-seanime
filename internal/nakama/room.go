package nakama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"seall/internal/constants"
	"time"
)

// Room represents a Seall Rooms relay room
type Room struct {
	ID          string    `json:"roomId"`
	HostWsUrl   string    `json:"hostWsUrl"`
	PeerJoinUrl string    `json:"peerJoinUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
	Password    string    `json:"-"` // Not returned by API, kept for reconnection
}

type CreateRoomRequest struct {
	Password string `json:"password"`
	Version  string `json:"version"`
}

type CreateRoomResponse struct {
	RoomID      string    `json:"roomId"`
	HostWsUrl   string    `json:"hostWsUrl"`
	PeerJoinUrl string    `json:"peerJoinUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// createRoom creates a new room on Seall Rooms
func (m *Manager) createRoom(password string) (*Room, error) {
	if constants.SeallRoomsApiUrl == "" {
		return nil, fmt.Errorf("rooms API not configured")
	}

	reqBody := CreateRoomRequest{
		Password: password,
		Version:  constants.Version,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		constants.SeallRoomsApiUrl,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create room: %s, %s", resp.Status, string(bodyBytes))
	}

	var createResp CreateRoomResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	room := &Room{
		ID:          createResp.RoomID,
		HostWsUrl:   createResp.HostWsUrl,
		PeerJoinUrl: createResp.PeerJoinUrl,
		CreatedAt:   createResp.CreatedAt,
		ExpiresAt:   createResp.ExpiresAt,
		Password:    password,
	}

	return room, nil
}
