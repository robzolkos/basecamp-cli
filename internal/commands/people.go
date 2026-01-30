package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// Person represents a Basecamp user
type Person struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	EmailAddress string  `json:"email_address"`
	Title        string  `json:"title"`
	Bio          string  `json:"bio"`
	Location     string  `json:"location"`
	Admin        bool    `json:"admin"`
	Owner        bool    `json:"owner"`
	Client       bool    `json:"client"`
	Employee     bool    `json:"employee"`
	TimeZone     string  `json:"time_zone"`
	AvatarURL    string  `json:"avatar_url"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	Company      Company `json:"company"`
}

type Company struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PersonOutput is the brief output for person listings
type PersonOutput struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Title    string `json:"title,omitempty"`
	Admin    bool   `json:"admin"`
	Owner    bool   `json:"owner"`
	Company  string `json:"company,omitempty"`
	Location string `json:"location,omitempty"`
}

// PeopleCmd lists all people
type PeopleCmd struct{}

type PeopleOutput struct {
	Count  int            `json:"count"`
	People []PersonOutput `json:"people"`
}

func (c *PeopleCmd) Run(args []string) error {
	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/people.json")
	if err != nil {
		return err
	}

	var people []Person
	if err := json.Unmarshal(data, &people); err != nil {
		return err
	}

	output := PeopleOutput{
		Count:  len(people),
		People: make([]PersonOutput, len(people)),
	}

	for i, p := range people {
		output.People[i] = personToOutput(p)
	}

	return PrintJSON(output)
}

// PersonCmd views a single person
type PersonCmd struct{}

type PersonDetailOutput struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Title     string `json:"title,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Location  string `json:"location,omitempty"`
	Admin     bool   `json:"admin"`
	Owner     bool   `json:"owner"`
	Client    bool   `json:"client"`
	Employee  bool   `json:"employee"`
	TimeZone  string `json:"time_zone"`
	Company   string `json:"company,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	CreatedAt string `json:"created_at"`
}

func (c *PersonCmd) Run(args []string) error {
	if len(args) < 1 {
		return errors.New("usage: basecamp person <person_id>")
	}
	personID := args[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/people/" + personID + ".json")
	if err != nil {
		return err
	}

	var person Person
	if err := json.Unmarshal(data, &person); err != nil {
		return err
	}

	return PrintJSON(PersonDetailOutput{
		ID:        person.ID,
		Name:      person.Name,
		Email:     person.EmailAddress,
		Title:     person.Title,
		Bio:       person.Bio,
		Location:  person.Location,
		Admin:     person.Admin,
		Owner:     person.Owner,
		Client:    person.Client,
		Employee:  person.Employee,
		TimeZone:  person.TimeZone,
		Company:   person.Company.Name,
		AvatarURL: person.AvatarURL,
		CreatedAt: person.CreatedAt,
	})
}

// PeoplePingableCmd lists pingable people
type PeoplePingableCmd struct{}

func (c *PeoplePingableCmd) Run(args []string) error {
	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/circles/people.json")
	if err != nil {
		return err
	}

	var people []Person
	if err := json.Unmarshal(data, &people); err != nil {
		return err
	}

	output := PeopleOutput{
		Count:  len(people),
		People: make([]PersonOutput, len(people)),
	}

	for i, p := range people {
		output.People[i] = personToOutput(p)
	}

	return PrintJSON(output)
}

// PeopleProjectCmd lists people on a project
type PeopleProjectCmd struct{}

type PeopleProjectOutput struct {
	ProjectID int            `json:"project_id"`
	Count     int            `json:"count"`
	People    []PersonOutput `json:"people"`
}

func (c *PeopleProjectCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/projects/" + projectID + "/people.json")
	if err != nil {
		return err
	}

	var people []Person
	if err := json.Unmarshal(data, &people); err != nil {
		return err
	}

	output := PeopleProjectOutput{
		ProjectID: 0,
		Count:     len(people),
		People:    make([]PersonOutput, len(people)),
	}

	// Parse project ID for output
	fmt.Sscanf(projectID, "%d", &output.ProjectID)

	for i, p := range people {
		output.People[i] = personToOutput(p)
	}

	return PrintJSON(output)
}

// MyProfileCmd shows the current user's profile
type MyProfileCmd struct{}

func (c *MyProfileCmd) Run(args []string) error {
	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/my/profile.json")
	if err != nil {
		return err
	}

	var person Person
	if err := json.Unmarshal(data, &person); err != nil {
		return err
	}

	return PrintJSON(PersonDetailOutput{
		ID:        person.ID,
		Name:      person.Name,
		Email:     person.EmailAddress,
		Title:     person.Title,
		Bio:       person.Bio,
		Location:  person.Location,
		Admin:     person.Admin,
		Owner:     person.Owner,
		Client:    person.Client,
		Employee:  person.Employee,
		TimeZone:  person.TimeZone,
		Company:   person.Company.Name,
		AvatarURL: person.AvatarURL,
		CreatedAt: person.CreatedAt,
	})
}

// ProjectAccessCmd updates who can access a project
type ProjectAccessCmd struct{}

type ProjectAccessOutput struct {
	Status  string   `json:"status"`
	Granted []string `json:"granted,omitempty"`
	Revoked []string `json:"revoked,omitempty"`
	Message string   `json:"message"`
}

func (c *ProjectAccessCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse --grant and --revoke flags
	var grant, revoke string

	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--grant":
			if i+1 < len(remaining) {
				grant = remaining[i+1]
				i++
			}
		case "--revoke":
			if i+1 < len(remaining) {
				revoke = remaining[i+1]
				i++
			}
		}
	}

	if grant == "" && revoke == "" {
		return errors.New("at least one of --grant or --revoke required (comma-separated person IDs)")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	payload := map[string]any{}
	if grant != "" {
		payload["grant"] = parseIDList(grant)
	}
	if revoke != "" {
		payload["revoke"] = parseIDList(revoke)
	}

	data, err := cl.Put("/projects/"+projectID+"/people/users.json", payload)
	if err != nil {
		return err
	}

	var result struct {
		Granted []Person `json:"granted"`
		Revoked []Person `json:"revoked"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	output := ProjectAccessOutput{
		Status:  "ok",
		Message: "Project access updated",
	}

	for _, p := range result.Granted {
		output.Granted = append(output.Granted, p.Name)
	}
	for _, p := range result.Revoked {
		output.Revoked = append(output.Revoked, p.Name)
	}

	return PrintJSON(output)
}

// Helper to convert Person to PersonOutput
func personToOutput(p Person) PersonOutput {
	return PersonOutput{
		ID:       p.ID,
		Name:     p.Name,
		Email:    p.EmailAddress,
		Title:    p.Title,
		Admin:    p.Admin,
		Owner:    p.Owner,
		Company:  p.Company.Name,
		Location: p.Location,
	}
}

// Helper to parse comma-separated ID list to int slice
func parseIDList(s string) []int {
	var ids []int
	for _, part := range splitComma(s) {
		var id int
		if _, err := fmt.Sscanf(part, "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}

func splitComma(s string) []string {
	var result []string
	current := ""
	for _, c := range s {
		if c == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else if c != ' ' {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
