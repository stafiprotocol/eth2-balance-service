package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// https://rtoken-api2.stafi.io/
func GetApyFromStafiInfo(stafiInfoEndpoint string) (float64, error) {
	url := stafiInfoEndpoint + "/stafi/webapi/rtoken/allStakeValueList"
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status: %d", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	if len(bodyBytes) == 0 {
		return 0, fmt.Errorf("bodyBytes zero err")
	}

	rsp := RspApy{}
	err = json.Unmarshal(bodyBytes, &rsp)
	if err != nil {
		return 0, err
	}

	if rsp.Status != "80000" {
		return 0, fmt.Errorf("GetApyFromStafiInfo status: %s", rsp.Status)
	}
	for _, apy := range rsp.Data.StakeValueList {
		if apy.RTokenType == -1 {
			return apy.TotalApy, nil
		}
	}

	return 0, fmt.Errorf("reth apy empty")
}

type RspApy struct {
	Data struct {
		TotalStakeValue float64 `json:"totalStakeValue"`
		StakeValueList  []struct {
			StakeAmount string  `json:"stakeAmount"`
			StakeValue  float64 `json:"stakeValue"`
			StakeApy    float64 `json:"stakeApy"`
			MintApy     float64 `json:"mintApy"`
			TotalApy    float64 `json:"totalApy"`
			RTokenType  int64   `json:"rTokenType"`
		} `json:"stakeValueList"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
