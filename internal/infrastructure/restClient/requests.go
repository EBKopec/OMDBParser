package restClient

import (
	"context"
	"net/http"
)

func Get(ctx context.Context, URL string) (*http.Response, error) {

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	client := &http.Client{}

	return client.Do(req)
}
