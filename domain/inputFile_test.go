package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var (
	lines = []string{"tt0001216	short	The Final Settlement	The Final Settlement	0	1910	\\N	16	Drama,Short",
		"tt0001219	short	Flores y perlas	Flores y perlas	0	1910	\\N	\\N	Drama,Short",
		"tt0001220	short	The Flower of the Ranch	The Flower of the Ranch	0	1910	\\N	\\N	Short,Western",
		"tt0001218	short	A Fallen Spirit	L'épée du spirite	0	1910	\\N	6	Animation,Fantasy,Horror",
	}
	params = ArgsFilter{
		TitleType:      "short",
		PrimaryTitle:   "",
		OriginalTitle:  "",
		Genres:         "Drama,Short",
		StartYear:      "1910",
		EndYear:        "",
		RuntimeMinutes: "",
	}
	paramsFail = ArgsFilter{
		TitleType:      "",
		PrimaryTitle:   "",
		OriginalTitle:  "",
		Genres:         "1",
		StartYear:      "",
		EndYear:        "",
		RuntimeMinutes: "",
	}

	omdbResult, omdbResultFail []DataFile
	returned, returnedFail     []DataFile
)

func init() {

}

func TestDataFile_Filters_Success(t *testing.T) {
	for _, values := range lines {
		value := strings.Split(values, "\t")
		assert.NotNil(t, value)

		inputService := EnvDataFile(value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7], value[8])
		assert.NotNil(t, inputService)

		omdbResult = append(omdbResult, *inputService)
		assert.NotNil(t, omdbResult)

		returned = InputService.Filters(inputService, map[string]interface{}{
			"PrimaryTitle":   params.PrimaryTitle,
			"TitleType":      params.TitleType,
			"OriginalTitle":  params.OriginalTitle,
			"Genres":         params.Genres,
			"StartYear":      params.StartYear,
			"EndYear":        params.EndYear,
			"RuntimeMinutes": params.RuntimeMinutes,
		}, omdbResult)
		assert.NotNil(t, returned)
	}
	fmt.Println(returned)
}

func TestDataFile_Filters_Fail(t *testing.T) {
	for _, values := range lines {
		value := strings.Split(values, "\t")
		assert.NotNil(t, value)

		inputService := EnvDataFile(value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7], value[8])
		assert.NotNil(t, inputService)

		omdbResultFail = append(omdbResultFail, *inputService)
		assert.NotNil(t, omdbResultFail)

		returnedFail = InputService.Filters(inputService, map[string]interface{}{
			"PrimaryTitle":   paramsFail.PrimaryTitle,
			"TitleType":      paramsFail.TitleType,
			"OriginalTitle":  paramsFail.OriginalTitle,
			"Genres":         paramsFail.Genres,
			"StartYear":      paramsFail.StartYear,
			"EndYear":        paramsFail.EndYear,
			"RuntimeMinutes": paramsFail.RuntimeMinutes,
		}, omdbResultFail)
		assert.Nil(t, returnedFail)
	}
}

//Benchmark below

func BenchmarkDataFile_Filters_Success(b *testing.B) {
	for _, values := range lines {
		value := strings.Split(values, "\t")
		assert.NotNil(b, value)

		inputService := EnvDataFile(value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7], value[8])
		assert.NotNil(b, inputService)

		omdbResult = append(omdbResult, *inputService)
		assert.NotNil(b, omdbResult)

		returned = InputService.Filters(inputService, map[string]interface{}{
			"PrimaryTitle":   params.PrimaryTitle,
			"TitleType":      params.TitleType,
			"OriginalTitle":  params.OriginalTitle,
			"Genres":         params.Genres,
			"StartYear":      params.StartYear,
			"EndYear":        params.EndYear,
			"RuntimeMinutes": params.RuntimeMinutes,
		}, omdbResult)
		assert.NotNil(b, returned)
	}
	fmt.Println(returned)
}

func BenchmarkDataFile_Filters_Fail(b *testing.B) {
	for _, values := range lines {
		value := strings.Split(values, "\t")
		assert.NotNil(b, value)

		inputService := EnvDataFile(value[0], value[1], value[2], value[3], value[4], value[5], value[6], value[7], value[8])
		assert.NotNil(b, inputService)

		omdbResultFail = append(omdbResultFail, *inputService)
		assert.NotNil(b, omdbResultFail)

		returnedFail = InputService.Filters(inputService, map[string]interface{}{
			"PrimaryTitle":   paramsFail.PrimaryTitle,
			"TitleType":      paramsFail.TitleType,
			"OriginalTitle":  paramsFail.OriginalTitle,
			"Genres":         paramsFail.Genres,
			"StartYear":      paramsFail.StartYear,
			"EndYear":        paramsFail.EndYear,
			"RuntimeMinutes": paramsFail.RuntimeMinutes,
		}, omdbResultFail)
		assert.Nil(b, returnedFail)
	}
}
