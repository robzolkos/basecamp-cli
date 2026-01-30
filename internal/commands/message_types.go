package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// MessageType represents a message board category
type MessageType struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// MessageTypesCmd lists message types
type MessageTypesCmd struct{}

type MessageTypesOutput struct {
	ProjectID    int                `json:"project_id"`
	MessageTypes []MessageTypeBrief `json:"message_types"`
}

type MessageTypeBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

func (c *MessageTypesCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/buckets/" + projectID + "/categories.json")
	if err != nil {
		return err
	}

	var types []MessageType
	if err := json.Unmarshal(data, &types); err != nil {
		return err
	}

	var pID int
	fmt.Sscanf(projectID, "%d", &pID)

	output := MessageTypesOutput{
		ProjectID:    pID,
		MessageTypes: make([]MessageTypeBrief, len(types)),
	}

	for i, t := range types {
		output.MessageTypes[i] = MessageTypeBrief{
			ID:   t.ID,
			Name: t.Name,
			Icon: t.Icon,
		}
	}

	return PrintJSON(output)
}

// MessageTypeCmd views a single message type
type MessageTypeCmd struct{}

type MessageTypeDetailOutput struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *MessageTypeCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("message_type_id required")
	}
	typeID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/buckets/" + projectID + "/categories/" + typeID + ".json")
	if err != nil {
		return err
	}

	var msgType MessageType
	if err := json.Unmarshal(data, &msgType); err != nil {
		return err
	}

	return PrintJSON(MessageTypeDetailOutput{
		ID:        msgType.ID,
		Name:      msgType.Name,
		Icon:      msgType.Icon,
		CreatedAt: msgType.CreatedAt,
		UpdatedAt: msgType.UpdatedAt,
	})
}

// MessageTypeCreateCmd creates a new message type
type MessageTypeCreateCmd struct{}

type MessageTypeCreateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Message string `json:"message"`
}

func (c *MessageTypeCreateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse flags
	var name, icon string

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--name":
			if i+1 < len(remaining) {
				name = remaining[i+1]
				i++
			}
		case "--icon":
			if i+1 < len(remaining) {
				icon = remaining[i+1]
				i++
			}
		}
	}

	if name == "" {
		return errors.New("--name required")
	}
	if icon == "" {
		return errors.New("--icon required (emoji)")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]any{
		"name": name,
		"icon": icon,
	}

	data, err := cl.Post("/buckets/"+projectID+"/categories.json", payload)
	if err != nil {
		return err
	}

	var msgType MessageType
	if err := json.Unmarshal(data, &msgType); err != nil {
		return err
	}

	return PrintJSON(MessageTypeCreateOutput{
		Status:  "ok",
		ID:      msgType.ID,
		Name:    msgType.Name,
		Icon:    msgType.Icon,
		Message: fmt.Sprintf("Message type '%s' created", msgType.Name),
	})
}

// MessageTypeUpdateCmd updates a message type
type MessageTypeUpdateCmd struct{}

type MessageTypeUpdateOutput struct {
	Status  string `json:"status"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Message string `json:"message"`
}

func (c *MessageTypeUpdateCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("message_type_id required")
	}
	typeID := remaining[0]

	// Parse flags
	var name, icon string

	for i := 1; i < len(remaining); i++ {
		switch remaining[i] {
		case "--name":
			if i+1 < len(remaining) {
				name = remaining[i+1]
				i++
			}
		case "--icon":
			if i+1 < len(remaining) {
				icon = remaining[i+1]
				i++
			}
		}
	}

	if name == "" && icon == "" {
		return errors.New("at least one of --name or --icon required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]any{}
	if name != "" {
		payload["name"] = name
	}
	if icon != "" {
		payload["icon"] = icon
	}

	data, err := cl.Put("/buckets/"+projectID+"/categories/"+typeID+".json", payload)
	if err != nil {
		return err
	}

	var msgType MessageType
	if err := json.Unmarshal(data, &msgType); err != nil {
		return err
	}

	return PrintJSON(MessageTypeUpdateOutput{
		Status:  "ok",
		ID:      msgType.ID,
		Name:    msgType.Name,
		Icon:    msgType.Icon,
		Message: fmt.Sprintf("Message type '%s' updated", msgType.Name),
	})
}

// MessageTypeDeleteCmd deletes a message type
type MessageTypeDeleteCmd struct{}

type MessageTypeDeleteOutput struct {
	Status  string `json:"status"`
	TypeID  string `json:"type_id"`
	Message string `json:"message"`
}

func (c *MessageTypeDeleteCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("message_type_id required")
	}
	typeID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	_, err = cl.Delete("/buckets/" + projectID + "/categories/" + typeID + ".json")
	if err != nil {
		return err
	}

	return PrintJSON(MessageTypeDeleteOutput{
		Status:  "ok",
		TypeID:  typeID,
		Message: "Message type deleted",
	})
}
