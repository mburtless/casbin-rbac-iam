package main

import (
	"fmt"
	"strings"
)

type User struct {
	Id int
	Name string
	ApiKey string
	Role int
}

var usersDb = map[string]User{
	"larry": {Id: 0, Name: "larry", Role: 0, ApiKey: "larry"},
	"bob": {Id: 1, Name: "bob", Role: 1, ApiKey: "bob"},
}

type Role struct {
	Id int
	Name string
}

var rolesDb = []Role{
	{Id: 0, Name: "admin"},
	{Id: 1, Name: "viewer"},
	{Id: 2, Name: "editor"},
	{Id: 3, Name: "irrelevant"},
}

func GetRoleById(id int) (*Role, error) {
	if len(rolesDb) > id {
		return &rolesDb[id], nil
	}
	return nil, fmt.Errorf("role with ID %d not found", id)
}

// GetCurrentUser returns the User object that corresponds to the provided apiKey
func GetCurrentUser(apiKey string) (*User, error)  {
	for _, u := range usersDb {
		if u.ApiKey == apiKey {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("user with apikey %s not found", apiKey)
}

// Zone resource
type Zone struct {
	Id   int
	Name string
	Org int
}

// SuffixMatch returns True if zone name has given suffix
func (z Zone) SuffixMatch(suffix string) bool {
	return strings.HasSuffix(z.Name, suffix)
}

var zonesDb = []Zone{
	{Id: 0, Name: "gmail.com", Org: 0},
	{Id: 1, Name: "react.net", Org: 0},
	{Id: 2, Name: "oso.com", Org: 0},
	{Id: 3, Name: "authz.net", Org: 0},
}

// GetZoneById fetches a zone resource by it's ID
func GetZoneById(id int) (*Zone, error) {
	if len(zonesDb) > id {
		return &zonesDb[id], nil
	}
	return nil, fmt.Errorf("zone with ID %d not found", id)
}

// Condition modifier for policies
type Condition struct {
	Type string `json:"type"`
	Value interface{} `json:"value"`
}