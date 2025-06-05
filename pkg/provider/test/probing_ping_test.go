package test

import (
	"fmt"
	"testing"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

func TestPinger(t *testing.T) {
	buildOption := provider.EndpointOption{
		Endpoint: "8.147.234.89",
		Timeout:  10,
		ICMP: provider.Eicmp{
			Interval: 1,
			Count:    5,
		},
	}

	pinger, err := provider.NewEndpointPinger().Pilot(buildOption)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(tools.JsonMarshal(pinger))
}
