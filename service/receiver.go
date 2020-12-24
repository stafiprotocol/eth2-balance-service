package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func ReceiveData(apiUrl string) (*BlockRawData, error) {
	//resp, err := http.Get(url)
	//if err != nil {
	//	return nil, fmt.Errorf("ReceiveData: request to %s error: %s", url, err)
	//}
	v := url.Values{}
	resp, err := http.PostForm(apiUrl, v)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ReceiveData: read body error: %s", err)
	}

	ari := new(ApiRawInfo)
	if err := json.Unmarshal(body, ari); err != nil {
		return nil, fmt.Errorf("ReceiveData: Unmarshal ApiRawInfo error: %s", err)
	}

	if ari.Status != "80000" {
		return nil, fmt.Errorf("ReceiveData: body status wrong: %s", ari.Status)
	}

	return ari.Data, nil
}
