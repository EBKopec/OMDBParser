package omdb

import (
	"context"
	"fmt"
	"gitlab.neoway.com.br/IMDBParser/internal/infrastructure/restClient"
	"io/ioutil"
	"net/http"
)

type OMDBServiceClient struct {
	httpClient *http.Client
	baseURL    string
	APIKey     string
}

func NewOmdbClient(client *http.Client, APIKey string) *OMDBServiceClient {
	return &OMDBServiceClient{
		httpClient: client,
		baseURL:    fmt.Sprintf("http://www.omdbapi.com/?apikey=%s", APIKey),
	}
}

func (os *OMDBServiceClient) GetData(filters string) ([]byte, error, string) {
	url := fmt.Sprintf("%s&%s", os.baseURL, filters)
	ctx := context.Background()
	resp, err := restClient.Get(ctx, url)
	if err != nil {
		return nil, err, url
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			return
		}
	}()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err, url
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s", bodyBytes)
		return bodyBytes, err, url
	}

	return bodyBytes, nil, url
}
