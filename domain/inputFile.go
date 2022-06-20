package domain

import (
	"reflect"
)

type DataFile struct {
	Tconst         string `tsv:"tconst"`
	TitleType      string `tsv:"titleType"`
	PrimaryTitle   string `tsv:"primaryTitle"`
	OriginalTitle  string `tsv:"originalTitle"`
	IsAdult        string `tsv:"isAdult"`
	StartYear      string `tsv:"startYear"`
	EndYear        string `tsv:"endYear"`
	RuntimeMinutes string `tsv:"runtimeMinutes"`
	Genres         string `tsv:"genres"`
}

func EnvDataFile(tconst, titleType, primaryTitle, originalTitle, isAdult, startYear, endYear, runtimeMinutes, genres string) *DataFile {
	return &DataFile{
		Tconst:         tconst,
		TitleType:      titleType,
		PrimaryTitle:   primaryTitle,
		OriginalTitle:  originalTitle,
		IsAdult:        isAdult,
		StartYear:      startYear,
		EndYear:        endYear,
		RuntimeMinutes: runtimeMinutes,
		Genres:         genres,
	}
}

func (df *DataFile) Filters(conf map[string]interface{}, values []DataFile) []DataFile {
	var returned []DataFile

	for _, val := range values {
		flag := true
		for k := range conf {
			if !reflect.DeepEqual(reflect.ValueOf(val).FieldByName(k).Interface(), conf[k]) && (conf[k] != "") {
				flag = false
				break
			}
		}
		if flag {
			returned = append(returned, val)
		}
	}
	return returned
}
