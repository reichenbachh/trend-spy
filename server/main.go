package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("tweet streaming server running on port :8080")
	http.HandleFunc("/",streamHandler)
	log.Fatal(http.ListenAndServe(":8080",nil))
}

func streamHandler(w http.ResponseWriter ,r *http.Request){

	// type StreamFilterParams struct {
	// 	FilterLevel   string   `url:"filter_level,omitempty"`
	// 	Follow        []string `url:"follow,omitempty,comma"`
	// 	Language      []string `url:"language,omitempty,comma"`
	// 	Locations     []string `url:"locations,omitempty,comma"`
	// 	StallWarnings *bool    `url:"stall_warnings,omitempty"`
	// 	Track         []string `url:"track,omitempty,comma"`
	// }
	w.Header().Set("Content-Type","application/json")

	if r.Body == nil{
		json.NewEncoder(w).Encode("Invalid Request")
		return
	}

	defer r.Body.Close()

	bodyBytes,err:= io.ReadAll(r.Body)

	if err != nil{
		json.NewEncoder(w).Encode("Invalid")
	}

	var streamRules *twitter.StreamFilterParams

	jsonErr:= json.Unmarshal(bodyBytes,&streamRules)

	if jsonErr != nil{
		log.Fatal(err)
	}

	twitterStream := createStream(streamRules)

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
	}


	go demux.HandleChan(twitterStream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	twitterStream.Stop()

	// for message := range twitterStream.Messages{
	// 	fmt.Println(message)
	// }
}

func createStream(streamRules *twitter.StreamFilterParams)(twitterStream *twitter.Stream){
	apiKey,apiSecret,tokenKey,tokenSecret := returnEnvVars()
	config:= oauth1.NewConfig(apiKey,apiSecret)
	token := oauth1.NewToken(tokenKey,tokenSecret)
	httpClient := config.Client(oauth1.NoContext,token)

	client:= twitter.NewClient(httpClient)

	twitterStream ,err := client.Streams.Filter(streamRules)

	if err != nil{
		fmt.Println("Stream creation failed")
		log.Fatal(err)
	}
	return twitterStream
}

func returnEnvVars() (apiKey, apiSecret, token, tokenSecret string) {
	err := godotenv.Load()
	if err !=nil{
		log.Fatal("Error loading API keys")
	}

	apiKey = os.Getenv("TWITTER_API_KEY")
	apiSecret = os.Getenv("TWITTER_BEARER_TOKEN")
	token=os.Getenv("TWITTER_ACCESS_TOKEN")
	tokenSecret= os.Getenv("TWITTER_TOKEN_SECRET")
	return apiKey,apiSecret,token,tokenSecret
}