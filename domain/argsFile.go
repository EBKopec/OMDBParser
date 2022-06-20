package domain

type ArgsFilter struct {
	TitleType      string "titleType,omitempty"
	PrimaryTitle   string "primaryTitle,omitempty"
	OriginalTitle  string "originalTitle,omitempty"
	Genres         string "genres,omitempty"
	StartYear      string "startYear,omitempty"
	EndYear        string "endYear,omitempty"
	RuntimeMinutes string "runtimeMinutes,omitempty"
}

func EnvArgsFilter(titleType, primaryTitle, genres, startYear, endYear, runtimeMinutes string) *ArgsFilter {
	return &ArgsFilter{
		TitleType:      titleType,
		PrimaryTitle:   primaryTitle,
		Genres:         genres,
		StartYear:      startYear,
		EndYear:        endYear,
		RuntimeMinutes: runtimeMinutes}
}

type ArgsBase struct {
	FilePath       string
	Workers        int
	MaxRunTime     string
	MaxRequests    int
	MaxApiRequests int
	PlotFilter     string
	APIKey         string
}

func EnvArgsBase(FilePath, MaxRunTime, PlotFilter, APIKey string, Workers, MaxRequests, MaxApiRequests int) *ArgsBase {
	return &ArgsBase{
		FilePath:       FilePath,
		Workers:        Workers,
		MaxRunTime:     MaxRunTime,
		MaxRequests:    MaxRequests,
		MaxApiRequests: MaxApiRequests,
		PlotFilter:     PlotFilter,
		APIKey:         APIKey}
}
