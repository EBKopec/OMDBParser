package domain

import (
	"github.com/stretchr/testify/assert"
	"gitlab.neoway.com.br/IMDBParser/internal/infrastructure/omdb"
	httpclient "gitlab.neoway.com.br/IMDBParser/pkg"
	"gitlab.neoway.com.br/IMDBParser/utils/misc"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"text/tabwriter"
	"time"
)

var (
	line = []string{"tt0001216	short	The Final Settlement	The Final Settlement	0	1910	\\N	16	Drama,Short",
		"tt0001219	short	Flores y perlas	Flores y perlas	0	1910	\\N	\\N	Drama,Short",
		"tt0001220	short	The Flower of the Ranch	The Flower of the Ranch	0	1910	\\N	\\N	Short,Western",
		"tt0001218	short	A Fallen Spirit	L'épée du spirite	0	1910	\\N	6	Animation,Fantasy,Horror",
	}
	line0 = []string{"								"}

	param = ArgsFilter{
		TitleType:      "short",
		PrimaryTitle:   "",
		OriginalTitle:  "",
		Genres:         "Drama,Short",
		StartYear:      "1910",
		EndYear:        "",
		RuntimeMinutes: "",
	}
	omdbResults, returns, returns0 []DataFile
	argsBase                       = &ArgsBase{
		FilePath:       "../test/data2.tsv",
		Workers:        2,
		MaxRunTime:     "1s",
		MaxRequests:    1,
		MaxApiRequests: 1,
		PlotFilter:     "Short",
		APIKey:         "000000",
	}
	argsBase404 = &ArgsBase{
		FilePath:       "x",
		Workers:        2,
		MaxRunTime:     "1s",
		MaxRequests:    1,
		MaxApiRequests: 1,
		PlotFilter:     "Short",
		APIKey:         "000000",
	}

	OMDBResp = OMDBResponse{
		Title:      "When Honor Calls",
		Year:       "1913",
		Rated:      "N/A",
		Released:   "01 Oct 1914",
		Runtime:    "N/A",
		Genre:      "Drama",
		Director:   "Curt A. Stark",
		Writer:     "Richard Voß (play)",
		Actors:     "Henny Porten, Harry Liedtke, Hans Marr, Frida Richard",
		Plot:       "N/A",
		Language:   "German",
		Country:    "Germany",
		Awards:     "N/A",
		Poster:     "N/A",
		Ratings:    []Ratings{{Source: "Internet Movie Database", Value: "5.0/10"}},
		Metascore:  "N/A",
		ImdbRating: "5.0",
		ImdbVotes:  "6",
		ImdbID:     "tt0002163",
		Type:       "movie",
		DVD:        "N/A",
		BoxOffice:  "N/A",
		Production: "N/A",
		Website:    "N/A",
		Response:   "True",
	}
)

func init() {

}

func TestService_Unique(t *testing.T) {
	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	omdbService := NewService(omdbClient, argsBase)

	for _, values := range line {
		value := strings.Split(values, "\t")
		assert.NotNil(t, value)
		inputService := EnvDataFile(value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7], value[8])
		assert.NotNil(t, inputService)
		omdbResults = append(omdbResults, *inputService)
		assert.NotNil(t, omdbResults)
		returns = InputService.Filters(inputService, map[string]interface{}{
			"PrimaryTitle":   param.PrimaryTitle,
			"TitleType":      param.TitleType,
			"OriginalTitle":  param.OriginalTitle,
			"Genres":         param.Genres,
			"StartYear":      param.StartYear,
			"EndYear":        param.EndYear,
			"RuntimeMinutes": param.RuntimeMinutes,
		}, omdbResults)
		assert.NotNil(t, returns)
	}
	result := omdbService.Unique(&returns)
	assert.NotNil(t, result)
}

func TestService_Working_Success(t *testing.T) {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	timeD, err := misc.ParseTimeDuration(argsBase.MaxRunTime)
	assert.NoError(t, err)
	assert.NotNil(t, timeD)
	assert.Equal(t, 1.0, timeD)

	timeOut := time.After(time.Duration(timeD) * time.Second)
	assert.NotNil(t, timeOut)

	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	assert.NotNil(t, omdbClient)

	omdbService := NewService(omdbClient, argsBase)
	assert.NotNil(t, omdbService)

	omdbService.Working(param, gracefulShutdown, timeOut)
}

func TestService_Working_File_404(t *testing.T) {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	timeD, err := misc.ParseTimeDuration(argsBase404.MaxRunTime)
	assert.NoError(t, err)
	assert.NotNil(t, timeD)
	assert.Equal(t, 1.0, timeD)

	timeOut := time.After(time.Duration(timeD) * time.Second)
	assert.NotNil(t, timeOut)

	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase404.APIKey)
	assert.NotNil(t, omdbClient)

	omdbService := NewService(omdbClient, argsBase404)
	assert.NotNil(t, omdbService)

	omdbService.Working(param, gracefulShutdown, timeOut)
}

func TestService_RequestData_API_Unauthorized(t *testing.T) {

	timeD, err := misc.ParseTimeDuration(argsBase.MaxRunTime)
	assert.NoError(t, err)
	assert.NotNil(t, timeD)
	assert.Equal(t, 1.0, timeD)

	timeOut := time.After(time.Duration(timeD) * time.Second)
	assert.NotNil(t, timeOut)

	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	assert.NotNil(t, omdbClient)

	omdbService := NewService(omdbClient, argsBase)
	assert.NotNil(t, omdbService)

	for _, values := range line {
		value := strings.Split(values, "\t")
		assert.NotNil(t, value)
		inputService := EnvDataFile(value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7], value[8])
		assert.NotNil(t, inputService)
		omdbResults = append(omdbResults, *inputService)
		assert.NotNil(t, omdbResults)
		returns = InputService.Filters(inputService, map[string]interface{}{
			"PrimaryTitle":   param.PrimaryTitle,
			"TitleType":      param.TitleType,
			"OriginalTitle":  param.OriginalTitle,
			"Genres":         param.Genres,
			"StartYear":      param.StartYear,
			"EndYear":        param.EndYear,
			"RuntimeMinutes": param.RuntimeMinutes,
		}, omdbResults)
		assert.NotNil(t, returns)
	}

	err = omdbService.RequestData(&returns)
	assert.Error(t, err)

}

func TestService_RequestData_Len0(t *testing.T) {
	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	assert.NotNil(t, omdbClient)
	omdbService := NewService(omdbClient, argsBase)
	assert.NotNil(t, omdbService)

	err := omdbService.RequestData(&returns0)
	assert.Error(t, err)
	assert.Equal(t, "there are no records with the filters passed", err.Error())

}

func TestService_FormatOutput(t *testing.T) {
	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	assert.NotNil(t, omdbClient)
	omdbService := NewService(omdbClient, argsBase)
	assert.NotNil(t, omdbService)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape|tabwriter.Debug)
	assert.NotNil(t, w)

	omdbService.FormatOutput(&OMDBResp, w)
}

//Benchmark Below
func BenchmarkService_Unique(b *testing.B) {
	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	omdbService := NewService(omdbClient, argsBase)

	for _, values := range line {
		value := strings.Split(values, "\t")
		assert.NotNil(b, value)
		inputService := EnvDataFile(value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7], value[8])
		assert.NotNil(b, inputService)
		omdbResults = append(omdbResults, *inputService)
		assert.NotNil(b, omdbResults)
		returns = InputService.Filters(inputService, map[string]interface{}{
			"PrimaryTitle":   param.PrimaryTitle,
			"TitleType":      param.TitleType,
			"OriginalTitle":  param.OriginalTitle,
			"Genres":         param.Genres,
			"StartYear":      param.StartYear,
			"EndYear":        param.EndYear,
			"RuntimeMinutes": param.RuntimeMinutes,
		}, omdbResults)
		assert.NotNil(b, returns)
	}
	result := omdbService.Unique(&returns)
	assert.NotNil(b, result)
}

func BenchmarkService_Working_Success(b *testing.B) {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	timeD, err := misc.ParseTimeDuration(argsBase.MaxRunTime)
	assert.NoError(b, err)
	assert.NotNil(b, timeD)
	assert.Equal(b, 1.0, timeD)

	timeOut := time.After(time.Duration(timeD) * time.Second)
	assert.NotNil(b, timeOut)

	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	assert.NotNil(b, omdbClient)

	omdbService := NewService(omdbClient, argsBase)
	assert.NotNil(b, omdbService)

	omdbService.Working(param, gracefulShutdown, timeOut)
}

func BenchmarkService_Working_File_404(b *testing.B) {
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	timeD, err := misc.ParseTimeDuration(argsBase404.MaxRunTime)
	assert.NoError(b, err)
	assert.NotNil(b, timeD)
	assert.Equal(b, 1.0, timeD)

	timeOut := time.After(time.Duration(timeD) * time.Second)
	assert.NotNil(b, timeOut)

	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase404.APIKey)
	assert.NotNil(b, omdbClient)

	omdbService := NewService(omdbClient, argsBase404)
	assert.NotNil(b, omdbService)

	omdbService.Working(param, gracefulShutdown, timeOut)
}

func BenchmarkService_RequestData_API_Unauthorized(b *testing.B) {

	timeD, err := misc.ParseTimeDuration(argsBase.MaxRunTime)
	assert.NoError(b, err)
	assert.NotNil(b, timeD)
	assert.Equal(b, 1.0, timeD)

	timeOut := time.After(time.Duration(timeD) * time.Second)
	assert.NotNil(b, timeOut)

	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	assert.NotNil(b, omdbClient)

	omdbService := NewService(omdbClient, argsBase)
	assert.NotNil(b, omdbService)

	for _, values := range line {
		value := strings.Split(values, "\t")
		assert.NotNil(b, value)
		inputService := EnvDataFile(value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7], value[8])
		assert.NotNil(b, inputService)
		omdbResults = append(omdbResults, *inputService)
		assert.NotNil(b, omdbResults)
		returns = InputService.Filters(inputService, map[string]interface{}{
			"PrimaryTitle":   param.PrimaryTitle,
			"TitleType":      param.TitleType,
			"OriginalTitle":  param.OriginalTitle,
			"Genres":         param.Genres,
			"StartYear":      param.StartYear,
			"EndYear":        param.EndYear,
			"RuntimeMinutes": param.RuntimeMinutes,
		}, omdbResults)
		assert.NotNil(b, returns)
	}

	err = omdbService.RequestData(&returns)
	assert.Error(b, err)

}

func BenchmarkService_RequestData_Len0(b *testing.B) {
	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	assert.NotNil(b, omdbClient)
	omdbService := NewService(omdbClient, argsBase)
	assert.NotNil(b, omdbService)

	err := omdbService.RequestData(&returns0)
	assert.Error(b, err)
	assert.Equal(b, "there are no records with the filters passed", err.Error())

}

func BenchmarkNewService_FormatOutput(b *testing.B) {
	omdbClient := omdb.NewOmdbClient(httpclient.NewHTTPClient(60*time.Second), argsBase.APIKey)
	assert.NotNil(b, omdbClient)
	omdbService := NewService(omdbClient, argsBase)
	assert.NotNil(b, omdbService)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape|tabwriter.Debug)
	assert.NotNil(b, w)

	omdbService.FormatOutput(&OMDBResp, w)
}
