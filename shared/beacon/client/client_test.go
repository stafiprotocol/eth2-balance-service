package client_test

import (
	"testing"

	"github.com/stafiprotocol/reth/shared/beacon"
	"github.com/stafiprotocol/reth/shared/beacon/client"
	"github.com/stafiprotocol/reth/types"
)

func TestStatus(t *testing.T) {
	c := client.NewStandardHttpClient("https://27Y0WDKrX1dYIkBXOugsSLh9hfr:a7c3849eba862fdd67382dab42e2a23c@eth2-beacon-mainnet.infura.io")
	pubkey, err := types.HexToValidatorPubkey("a06aaf5586589ee0953fe3345303d4193636ffcf9f87e1c39090171d6a47c12c0909868ebf4894b9278451ddb796c8c3")
	if err != nil {
		t.Fatal(err)
	}
	slot := uint64(4684154)
	status, err := c.GetValidatorStatus(pubkey, &beacon.ValidatorStatusOptions{
		Slot: &slot,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", status)
}
