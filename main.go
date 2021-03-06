package main

import (
	"os"
    "flag"
	"fmt"
	"log"
	"time"
    "github.com/joho/godotenv"
	"github.com/whitef0x0/TrendingGitlab/flags"
	"github.com/whitef0x0/TrendingGitlab/storage"
	"github.com/whitef0x0/TrendingGitlab/twitter"
	"github.com/whitef0x0/TrendingGitlab/expvar"
	"github.com/whitef0x0/TrendingGitlab/tweets"
)

const (
	// Version of @TrendingGitlab
	Version = "0.4.0"
)

func main() {
    
    errLogFile, errOpeningErrLogFile := os.OpenFile("./logs/go_err.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if errOpeningErrLogFile != nil {
        log.Fatal("Error opening go_err.log file: %v", errOpeningErrLogFile)
    }   
    defer errLogFile.Close()
    
    log.SetOutput(errLogFile)
   
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		// Twitter
		twitterConsumerKey       = flags.String("twitter-consumer-key", "TrendingGitlab_TWITTER_CONSUMER_KEY", "", "Twitter-API: Consumer key. Env var: TrendingGitlab_TWITTER_CONSUMER_KEY")
		twitterConsumerSecret    = flags.String("twitter-consumer-secret", "TrendingGitlab_TWITTER_CONSUMER_SECRET", "", "Twitter-API: Consumer secret. Env var: TrendingGitlab_TWITTER_CONSUMER_SECRET")
		twitterAccessToken       = flags.String("twitter-access-token", "TrendingGitlab_TWITTER_ACCESS_TOKEN", "", "Twitter-API: Access token. Env var: TrendingGitlab_TWITTER_ACCESS_TOKEN")
		twitterAccessTokenSecret = flags.String("twitter-access-token-secret", "TrendingGitlab_TWITTER_ACCESS_TOKEN_SECRET", "", "Twitter-API: Access token secret. Env var: TrendingGitlab_TWITTER_ACCESS_TOKEN_SECRET")
		twitterFollowNewPerson   = flags.Bool("twitter-follow-new-person", "TrendingGitlab_TWITTER_FOLLOW_NEW_PERSON", false, "Twitter: Follows a friend of one of our followers. Env var: TrendingGitlab_TWITTER_FOLLOW_NEW_PERSON")

		// Timings
		tweetTime                = flags.Duration("twitter-tweet-time", "TRENDINGGITHUB_TWITTER_TWEET_TIME", 120*time.Minute, "Twitter: Time interval to search a new project and tweet it. Env var: TRENDINGGITHUB_TWITTER_TWEET_TIME")
		configurationRefreshTime = flags.Duration("twitter-conf-refresh-time", "TRENDINGGITHUB_TWITTER_CONF_REFRESH_TIME", 24*time.Hour, "Twitter: Time interval to refresh the configuration of twitter (e.g. char length for short url). Env var: TRENDINGGITHUB_TWITTER_CONF_REFRESH_TIME")
		followNewPersonTime      = flags.Duration("twitter-follow-new-person-time", "TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON_TIME", 45*time.Minute, "Growth hack: Time interval to search for a new person to follow. Env var: TRENDINGGITHUB_TWITTER_FOLLOW_NEW_PERSON_TIME")

		// Redis storage
		storageURL  = flags.String("storage-url", "TrendingGitlab_STORAGE_URL", ":6379", "Storage URL (e.g. 1.2.3.4:6379 or :6379). Env var: TrendingGitlab_STORAGE_URL")
		storageAuth = flags.String("storage-auth", "TrendingGitlab_STORAGE_AUTH", "", "Storage Auth (e.g. myPassword or <empty>). Env var: TrendingGitlab_STORAGE_AUTH")

		expVarPort  = flags.Int("expvar-port", "TRENDINGGITHUB_EXPVAR_PORT", 8311, "Port which will be used for the expvar TCP server. Env var: TRENDINGGITHUB_EXPVAR_PORT")
		showVersion = flags.Bool("version", "TRENDINGGITHUB_VERSION", false, "Outputs the version number and exit. Env var: TRENDINGGITHUB_VERSION")
		debugMode   = flags.Bool("debug", "TRENDINGGITHUB_DEBUG", false, "Outputs the tweet instead of tweet it (useful for development). Env var: TRENDINGGITHUB_DEBUG")
	)
	flag.Parse()

	// Output the version and exit
	if *showVersion {
		fmt.Printf("@TrendingGitlab v%s\n", Version)
		return
	}

	log.Println("Hey, nice to meet you. My name is @GitlabTrending. Lets get ready to tweet some trending content!")
	defer log.Println("Nice sesssion. A lot of knowledge was tweeted. Good work and see you next time!")

	// Prepare the twitter client
	twitterClient := twitter.NewClient(*twitterConsumerKey, *twitterConsumerSecret, *twitterAccessToken, *twitterAccessTokenSecret, *debugMode)

	// When we are running in a debug mode, we are running with a debug configuration.
	// So we don`t need to load the configuration from twitter here.
	if *debugMode == false {
		err := twitterClient.LoadConfiguration()
		if err != nil {
			log.Fatalf("Twitter Configuration initialisation failed: %s", err)
		}
		log.Printf("Twitter Configuration initialisation success: ShortUrlLength %d\n", twitterClient.Configuration.ShortUrlLength)
		twitterClient.SetupConfigurationRefresh(*configurationRefreshTime)
	}

	// Activate our growth hack feature
	// Checkout the README for details or read the code (suggested).
	if *twitterFollowNewPerson {
		log.Println("Growth hack \"Follow a friend of a friend\" activated")
		twitterClient.SetupFollowNewPeopleScheduling(*followNewPersonTime)
	}

	// Request a storage backend
	storageBackend := storage.NewBackend(*storageURL, *storageAuth, *debugMode)
	defer storageBackend.Close()
	log.Println("Storage backend initialisation success")

	// Start the exvar server
	err := expvar_server.StartExpvarServer(*expVarPort)
	if err != nil {
		log.Fatalf("Expvar initialisation failed: %s", err)
	}
	log.Println("Expvar initialisation started ...")

	// Let the party begin
	tweets.StartTweeting(twitterClient, storageBackend, *tweetTime)
}
