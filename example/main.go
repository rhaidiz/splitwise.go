package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rhaidiz/splitwise.go"
	"golang.org/x/oauth2"
)

func main() {
	// auth := splitwise.NewAPIKeyAuth(os.Getenv("API_KEY"))
	// client := splitwise.NewClient(auth)

	conf := &oauth2.Config{
		ClientID:     "LNUtEksV8VScipH5bRgy2eHd1G98WWnP61p6S0mE",
		ClientSecret: "zFyE9zclIiY8Dn4t0aT3Tc4KJIfIKxzWd8kWvh9K",
		RedirectURL:  "http://127.0.0.1:8081/callback",
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://secure.splitwise.com/oauth/token",
			AuthURL:  "https://secure.splitwise.com/oauth/authorize",
		},
	}
	c := splitwise.NewAuth0Client(conf)
	url := c.GetOAuth2AuthorizeURL()
	log.Print(url)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	c.SetOAuth2Code(code)

	userExamples(c)
}

func userExamples(client splitwise.Client) {
	currentUser, err := client.CurrentUser(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(currentUser.FirstName)
	expenses, err := client.Expenses(context.Background(), "2022-05-01T00:00:00Z", "2022-05-31T00:00:00Z")
	if err != nil {
		panic(err)
	}
	fmt.Println(expenses)
	for i := 0; i < len(expenses.Exps); i++ {
		for j := 0; j < len(expenses.Exps[i].Users); j++ {
			if expenses.Exps[i].Users[j].User_id == currentUser.ID {
				date := expenses.Exps[i].Date
				category := expenses.Exps[i].Category.Name
				cost := expenses.Exps[i].Users[j].Owed_share
				description := expenses.Exps[i].Description
				fmt.Printf("%s %s %s %s\n", date, category, description, cost)
			}
		}
	}
	fmt.Println(len(expenses.Exps))
}
