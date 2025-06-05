package test

import (
	"fmt"
	"testing"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

func TestNewEndpointTcper(t *testing.T) {
	buildOption := provider.EndpointOption{
		Endpoint: "8.147.234.89:80",
		Timeout:  10,
	}
	pilot, err := provider.NewEndpointTcper().Pilot(buildOption)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(tools.JsonMarshal(pilot))
}
