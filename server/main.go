package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

func main() {
	

}

func streamHandler(w http.ResponseWriter ,r http.Request){
	w.Header().Set("Content-Type","application/json")

	if r.Body == nil{
		json.NewEncoder(w).Encode("Invalid Request")
		return
	}

	_,err:= io.ReadAll(r.Body)

	if err != nil{
		json.NewEncoder(w).Encode("inv")
	}


}

func createStream(streamRules *twitter.StreamFilterParams)(twitterStream *twitter.Stream){
	apiKey,apiSecret,tokenKey,tokenSecret := returnEnvVars()
	config:= oauth1.NewConfig(apiKey,apiSecret)
	token := oauth1.NewToken(tokenKey,tokenSecret)

	httpClient := config.Client(oauth1.NoContext,token)

	client:= twitter.NewClient(httpClient)

	twitterStream ,err := client.Streams.Filter(streamRules)

	if err != nil{
		fmt.Println("Stream failed")
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