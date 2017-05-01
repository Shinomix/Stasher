package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Config format of conf.json file
type Config struct {
	BaseURL     string
	AccessToken string
	Username    string
}

// Single-user response format
type user struct {
	ID       string `json:"id"`
	Username string `json:"name"`
}

// HTTP /users.list json response structure
type usersListAPIResponse struct {
	IsOkay    bool   `json:"ok"`
	UsersList []user `json:"members"`
}

func parseResponse(body []byte) (*usersListAPIResponse, error) {
	var s = new(usersListAPIResponse)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("[getUsers] Failed to unmarshal http response body: ", err)
	}
	return s, err
}

func listUsers(config *Config) []user {
	params := url.Values{}
	params.Add("token", config.AccessToken)
	params.Add("presence", "false")

	resp, err := http.Get(config.BaseURL + "/users.list?" + params.Encode())
	if err != nil {
		fmt.Println("[listUsers] Failed to call /users.list: ", err)
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err.Error())
	}
	users, _ := parseResponse(body)
	return users.UsersList
}

// GetSelfID fetch team's users from Slack and returns an ID for user matching Config.Username
func GetSelfID(config *Config) string {
	var toReturn string

	usersList := listUsers(config)
	for i := range usersList {
		if strings.Contains(config.Username, usersList[i].Username) {
			toReturn = usersList[i].ID
			break
		}
	}
	return toReturn
}

// LoadConfig fetch data from conf.json and fill pointer of Config object passed in params
func LoadConfig(conf *Config) bool {
	file, _ := Asset("conf/conf.json")
	decoder := json.NewDecoder(strings.NewReader(string(file)))
	err := decoder.Decode(conf)
	if err != nil {
		fmt.Println("[error] Couldn't load conf file: ", err)
		return false
	}
	return true
}
