package main

import (
	"fmt"

	"github.com/djgilcrease/survey"
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
)

type user struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Username string `json:"username"`
	Email string `json:"email"`
	Address *address `json:"address"`
}

type address struct {
	Street string `json:"street"`
	Suite string `json:"suite"`
	City string `json:"city"`
	Zip string `json:"zipcode"`
}
type users = []*user

// the questions to ask
var userPrompt = survey.NewSingleSelect().SetMessage("Select User:")
var usersPrompt = survey.NewMultiSelect().SetMessage("Select Users:")
var simpleQs = []*survey.Question{
	{
		Name: "user",
		Prompt: userPrompt,
		Validate: survey.Required,
	},
	{
		Name: "users",
		Prompt: usersPrompt,
		Validate: survey.Required,
	},
}

func init() {
	var (
		userData []byte
		request *http.Request
		response *http.Response
		err error
	)
	httpClient := &http.Client{Timeout: 5*time.Second}
	if request, err = http.NewRequest("GET", "https://jsonplaceholder.typicode.com/users", nil); err != nil {
		fmt.Println(err.Error())
		return
	}
	if response, err = httpClient.Do(request); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer response.Body.Close()
	userData, err = ioutil.ReadAll(response.Body)
	var us users
	if err = json.Unmarshal(userData, &us); err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, _user := range us {
		userPrompt.AddOption(_user.Username, _user, false)
		usersPrompt.AddOption(_user.Username, _user, false)
	}

}

func main() {
	answers := struct {
		User  *user
		Users users
	}{}


	// ask the question
	err := survey.Ask(simpleQs, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("%s has the username %s and thier address is %+v\r\n", answers.User.Name, answers.User.Username, answers.User.Address)
	fmt.Printf("Selected %d Users\r\n", len(answers.Users))
	fmt.Printf("%s has the username %s and thier address is %+v\r\n", answers.Users[0].Name, answers.Users[0].Username, answers.Users[0].Address)
}
