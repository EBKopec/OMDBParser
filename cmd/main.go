package main

import (
	"flag"
	"fmt"
	"gitlab.neoway.com.br/IMDBParser/domain"
	"gitlab.neoway.com.br/IMDBParser/internal/infrastructure/omdb"
	httpclient "gitlab.neoway.com.br/IMDBParser/pkg"
	"gitlab.neoway.com.br/IMDBParser/utils/misc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	fileName       string
	titleType      string
	primaryTitle   string
	maxRunTime     string
	genres         string
	startYear      string
	endYear        string
	runtimeMinutes string
	workers        int
	plotFilter     string
	maxRequests    int
	maxApiRequests int
	APIKey         string
)

func init() {
	flag.StringVar(&maxRunTime, "timeOut", "1s", "Should not there be timeout here?")
	flag.IntVar(&workers, "workers", 5, "Should not there be workers here?")
	flag.StringVar(&fileName, "filePath", "", "Should not there be file here?")
	flag.StringVar(&primaryTitle, "primaryTitle", "", "Should not there be a primary title here?")
	flag.StringVar(&titleType, "titleType", "", "Should not there be a title type here?")
	flag.StringVar(&genres, "genres", "", "Should not there be genres here?")
	flag.StringVar(&startYear, "startYear", "", "Should not there be a start year here?")
	flag.StringVar(&endYear, "endYear", "", "Should not there be an end year here?")
	flag.StringVar(&runtimeMinutes, "runtimeMinutes", "", "Should not there be run time minutes here?")
	flag.StringVar(&plotFilter, "plot", "short", "Should not there be a plot here?")
	flag.IntVar(&maxRequests, "maxRequests", 1, "Should not there be max requests numbers here?")
	flag.IntVar(&maxApiRequests, "maxApiRequests", 1, "Should not there be max API requests numbers here?")
	flag.StringVar(&APIKey, "apiKey", "", "Should not there be an API key here?")
	flag.Parse()
}

func main() {
	//apiKey := "77b33efe"
	argsFilter := domain.EnvArgsFilter(titleType, primaryTitle, genres, startYear, endYear, runtimeMinutes)
	argsBase := domain.EnvArgsBase(fileName, maxRunTime, plotFilter, APIKey, workers, maxRequests, maxApiRequests)

	//Channels
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	if !isFlagPassed("apiKey") {
		fmt.Printf("You must provide an API Key in the \"apiKey\" Flag %s\n", argsBase.APIKey)
		os.Exit(0)
	}
	//Client
	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)

	timeD, err := misc.ParseTimeDuration(argsBase.MaxRunTime)
	if err != nil {
		fmt.Printf("invalid format %s", err)
	}
	timeOut := time.After(time.Duration(timeD) * time.Second)

	//Service
	omdbService := domain.NewService(omdbClient, argsBase)

	omdbService.Working(*argsFilter, gracefulShutdown, timeOut)

}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
