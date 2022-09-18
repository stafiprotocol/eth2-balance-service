package client

import "github.com/stafiprotocol/reth/shared/beacon"

type NimbusClient struct {
	StandardHttpClient
}

// Create a new client instance
func NewNimbusClient(providerAddress string) (*NimbusClient, error) {
	client, err := NewStandardHttpClient(providerAddress)
	if err != nil {
		return nil, err
	}
	return &NimbusClient{
		StandardHttpClient: *client,
	}, nil
}

func (n *NimbusClient) GetClientType() (beacon.BeaconClientType, error) {
	return beacon.SingleProcess, nil
}
