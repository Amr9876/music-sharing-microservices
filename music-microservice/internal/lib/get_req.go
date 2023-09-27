package lib

import (
	"encoding/json"
	"io"
	"net/http"
)

func SendHttpGetRequest(url string, token string) (interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	if len(token) != 0 {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var data interface{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	return data, nil
}
