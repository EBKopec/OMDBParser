package domain

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.neoway.com.br/IMDBParser/internal/infrastructure/omdb"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

type Service struct {
	omdbServiceClient OMDBServiceClient
	argsBase          ArgsBase
	argsFilter        ArgsFilter
	returned          []DataFile
	omdbResult        []DataFile
	counter           int
	Err               error
}

func NewService(omdbServiceClient *omdb.OMDBServiceClient, base *ArgsBase) *Service {
	return &Service{
		omdbServiceClient: omdbServiceClient,
		argsBase:          *base,
	}
}

func (s *Service) Working(params ArgsFilter, shutdown chan os.Signal, timeOut <-chan time.Time) {

	var wg sync.WaitGroup

	result := make(chan *[]DataFile)
	workerChannel := make(chan bool, s.argsBase.Workers)

	file, err := os.Open(s.argsBase.FilePath)
	if err != nil {
		fmt.Println(err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	s.TimeoutFunc(timeOut, func() {

		for scanner.Scan() {
			workerChannel <- true
			wg.Add(1)

			go func(lines string) {
				defer func() {
					wg.Done()
					<-workerChannel
				}()

				value := strings.Split(lines, "\t")
				inputService := EnvDataFile(value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7], value[8])

				s.omdbResult = append(s.omdbResult, *inputService)
				s.returned = InputService.Filters(inputService, map[string]interface{}{
					"PrimaryTitle":   params.PrimaryTitle,
					"TitleType":      params.TitleType,
					"OriginalTitle":  params.OriginalTitle,
					"Genres":         params.Genres,
					"StartYear":      params.StartYear,
					"EndYear":        params.EndYear,
					"RuntimeMinutes": params.RuntimeMinutes,
				}, s.omdbResult)
			}(scanner.Text())
			result <- &s.returned
		}
		close(result)
		wg.Wait()

	}, result, shutdown)

}

func (s *Service) TimeoutFunc(timeout <-chan time.Time, myFunc func(), result chan *[]DataFile, gracefulShutdown chan os.Signal) {
	finished := make(chan bool)

	go func() {
		myFunc()
		finished <- true
	}()

loop:
	for {

		// select statement will block this thread until one of the three conditions below is met
		select {
		case <-result:
		case <-gracefulShutdown:
			fmt.Println("Release the Kraken :D")
			_, cancel := context.WithCancel(context.Background())
			data := <-result
			s.handleTermination(cancel, data)
			break loop
		case <-timeout:
			fmt.Println("Release the Kraken :D")
			data := <-result
			err := s.RequestData(data)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Timeout reached")
			os.Exit(10)
		case <-finished:
			break loop
		}
	}
}

func (s *Service) handleTermination(cancel context.CancelFunc, data *[]DataFile) {
	cancel()
	err := s.RequestData(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("program exiting")
}

func (s *Service) RequestData(data *[]DataFile) error {
	var responseOMDB *OMDBResponse
	var fail FailRequest

	if len(*data) == 0 {
		s.Err = fmt.Errorf("there are no records with the filters passed")
		return s.Err
	}
	dataResult := s.Unique(data)
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape|tabwriter.Debug)
	for _, v := range dataResult {
		if s.argsBase.MaxApiRequests > 0 {
			s.argsBase.MaxApiRequests--
			filters := fmt.Sprintf("&t=%s&i=%s&plot=%s&type=%s&y=%s", s.argsFilter.PrimaryTitle, v.Tconst,
				s.argsBase.PlotFilter, s.argsFilter.TitleType, s.argsFilter.StartYear)

			response, err, url := s.omdbServiceClient.GetData(filters)
			if err != nil {
				jErr := json.Unmarshal(response, &fail)
				if jErr != nil {
					s.Err = fmt.Errorf("error on trying unmarshal")
					return s.Err
				}
				s.Err = fmt.Errorf("URL: [%s]\n%s - Unauthorized", url, fail.Error)
				return s.Err
			}

			err = json.Unmarshal(response, &responseOMDB)
			if err != nil {
				s.Err = err
			}
			s.FormatOutput(responseOMDB, w)
		}
	}
	return s.Err
}

func (s *Service) FormatOutput(data *OMDBResponse, w *tabwriter.Writer) {
	fmt.Fprintln(w, "IMDB_ID  \t  Title    \t  Plot  ")
	fmt.Fprintln(w, fmt.Sprintf("%s  \t  %s    \t  %s  ", data.ImdbID, data.Title, data.Plot))
	w.Flush()
}

func (s *Service) Unique(sample *[]DataFile) []DataFile {
	var unique []DataFile
sampleLoop:
	for _, v := range *sample {
		for i, u := range unique {
			if v.Tconst == u.Tconst && v.Tconst != " " && u.Tconst != " " {
				unique[i] = v
				continue sampleLoop
			}
		}
		unique = append(unique, v)
	}
	return unique
}
