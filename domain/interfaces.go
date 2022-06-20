package domain

import (
	"os"
	"text/tabwriter"
	"time"
)

type InputService interface {
	Filters(conf map[string]interface{}, values []DataFile) []DataFile
}

type OMDBServiceClient interface {
	GetData(filters string) ([]byte, error, string)
}

type OMDBService interface {
	RequestData(*[]DataFile) error
	FormatOutput(data *OMDBResponse, w *tabwriter.Writer)
	Working(params ArgsFilter, shutdown chan os.Signal, timeOut <-chan time.Time)
}
