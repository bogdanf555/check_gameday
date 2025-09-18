package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gregdel/pushover"
	"github.com/joho/godotenv"
)

// curl -H "X-Auth-Token: "Your token here" "https://api.football-data.org/v4/teams/64/matches?season=2025&status=SCHEDULED"

// NOTE : This program now when ran will display the next game liverpool has.
//       Inhencements to be done:
//       1.Insert proper checks on the unmarshaled map for the keys to exist
//       4.Make it a project with multiple files: Football-data fetcher; Notification Sender ; main.go etc.
//       5.Write some tests
//       6. make it so that the cron job can be set fast with sh script
//       7. Do documentation

var (
	MY_FOOTBALL_DATA_ORG_AUTH_TOKEN string
	MY_PUSHOVER_RECEPIENT           string
	MY_PUSHOVER_API_TOKEN           string
)

var LiverpoolScheduledGamesUrl = "https://api.football-data.org/v4/teams/64/matches?season=2025&status=SCHEDULED&limit=1"
var LiverpoolId = 64

type Match struct {
	opponentTeamName    string
	matchStartTimeLocal time.Time
}

func GetNextLiverpoolGame() (map[string]interface{}, error) {

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", LiverpoolScheduledGamesUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", MY_FOOTBALL_DATA_ORG_AUTH_TOKEN)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := make(map[string]interface{})
	json.Unmarshal(body, &r)

	return r, nil
}

func FetchNextLiverpoolGame() (Match, error) {

	nextMatch, err := GetNextLiverpoolGame()
	if err != nil {
		fmt.Println(err)
		return Match{}, err
	}

	matchesKey, ok := nextMatch["matches"]
	if !ok {
		panic("the key doesn't exist")
	}

	matchesSlice, ok := matchesKey.([]interface{})
	if !ok {
		panic("couldn't assert the key")
	}

	for _, match := range matchesSlice {
		// TODO : instead of looping only to look at the first element extract it from here
		matchMap := match.(map[string]interface{})

		//TODO : some checks should be done in fact to see if the keys exists in the map

		homeTeam := matchMap["homeTeam"].(map[string]interface{})
		awayTeam := matchMap["awayTeam"].(map[string]interface{})

		var opponentTeamName string
		if int(homeTeam["id"].(float64)) != LiverpoolId {
			opponentTeamName = homeTeam["name"].(string)
		} else {
			opponentTeamName = awayTeam["name"].(string)
		}

		dateTimeUTCString := matchMap["utcDate"].(string)

		matchStartTimeUTC, err := time.Parse(time.RFC3339, dateTimeUTCString)

		if err != nil {
			panic(err)
		}

		return Match{opponentTeamName, matchStartTimeUTC.Local()}, nil
	}

	return Match{}, nil
}

func SendMatchNotificationToPushOver(matchMessage string) error {
	message := pushover.NewMessageWithTitle(matchMessage, "Liverpool Game Today")
	app := pushover.New(MY_PUSHOVER_API_TOKEN)
	recepient := pushover.NewRecipient(MY_PUSHOVER_RECEPIENT)

	response, err := app.SendMessage(message, recepient)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}

func LoadEnv() {
	// load .env file
	godotenv.Load(".env")

	MY_FOOTBALL_DATA_ORG_AUTH_TOKEN = os.Getenv("MY_FOOTBALL_DATA_ORG_AUTH_TOKEN")
	MY_PUSHOVER_API_TOKEN = os.Getenv("MY_PUSHOVER_API_TOKEN")
	MY_PUSHOVER_RECEPIENT = os.Getenv("MY_PUSHOVER_RECIPIENT")

	fmt.Println("Loaded env vars")
}

func main() {

	LoadEnv()

	match, err := FetchNextLiverpoolGame()
	if err != nil {
		panic(err)
	}

	t := match.matchStartTimeLocal
	now := time.Now()
	var sameDay bool = (now.Year() == t.Year()) && (now.Month() == t.Month()) && (now.Day() == t.Day())

	if sameDay {

		hour, minute, _ := t.Clock()
		matchMessage := fmt.Sprintf("Liverpool will face %v today at: %02d:%02d!", match.opponentTeamName, hour, minute)

		err := SendMatchNotificationToPushOver(matchMessage)

		if err != nil {
			panic(err)
		}
		fmt.Println("Message successfully sent! There's a game today :>")
	} else {
		fmt.Println("No game today...")
	}
}
