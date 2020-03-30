package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	logfilename := "/log/tweeter-" + hostname + ".log"
	//	logfilename := "tweeter-" + hostname + ".log"
	fmt.Println(logfilename)
	f, err := os.OpenFile(logfilename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	logger := log.New(f, "prefix", log.LstdFlags)

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	keyword1 := flags.String("keyword1", "", "Twitter Search Keyword")
	keyword2 := flags.String("keyword2", "", "Twitter Search Keyword")
	keyword3 := flags.String("keyword3", "", "Twitter Search Keyword")
	keyword4 := flags.String("keyword4", "", "Twitter Search Keyword")
	keyword5 := flags.String("keyword5", "", "Twitter Search Keyword")
	keyword6 := flags.String("keyword6", "", "Twitter Search Keyword")
	keyword7 := flags.String("keyword7", "", "Twitter Search Keyword")
	keyword8 := flags.String("keyword8", "", "Twitter Search Keyword")
	keyword9 := flags.String("keyword9", "", "Twitter Search Keyword")
	keyword10 := flags.String("keyword10", "", "Twitter Search Keyword")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		//		fmt.Println(tweet.Text)
		//tweetJSON, _ := json.Marshal(tweet)
		//fmt.Println(string(tweetJSON))
		logger.Println(tweet.Text)
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		logger.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		logger.Printf("%#v\n", event)
	}

	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{*keyword1, *keyword2, *keyword3, *keyword4, *keyword5, *keyword6, *keyword7, *keyword8, *keyword9, *keyword10},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// USER (quick test: auth'd user likes a tweet -> event)
	// userParams := &twitter.StreamUserParams{
	// 	StallWarnings: twitter.Bool(true),
	// 	With:          "followings",
	// 	Language:      []string{"en"},
	// }
	// stream, err := client.Streams.User(userParams)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// SAMPLE
	// sampleParams := &twitter.StreamSampleParams{
	// 	StallWarnings: twitter.Bool(true),
	// }
	// stream, err := client.Streams.Sample(sampleParams)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
