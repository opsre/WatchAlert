package test

import (
	"fmt"
	"testing"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

func TestNewEndpointSSLer(t *testing.T) {
	buildOption := provider.EndpointOption{
		Endpoint: "www.baidu.com",
		Timeout:  10,
	}
	pilot, err := provider.NewEndpointSSLer().Pilot(buildOption)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(tools.JsonMarshal(pilot))
}
