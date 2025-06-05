package test

import (
	"fmt"
	"testing"
	"watchAlert/internal/models"
	"watchAlert/pkg/provider"
	"watchAlert/pkg/tools"
)

func TestHTTPer(t *testing.T) {
	geter()
	//poster()
}

func geter() {
	buildOption := provider.EndpointOption{
		Endpoint: "https://docsify.js.org/adsf",
		Timeout:  10,
		HTTP: provider.Ehttp{
			Method: provider.GetHTTPMethod,
			Header: map[string]string{},
			Body:   "",
		},
	}

	per, err := provider.NewEndpointHTTPer().Pilot(buildOption)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(tools.JsonMarshal(per))
}

func poster() {
	var user models.Member
	user.UserName = "admin"
	user.Password = "123"
	buildOption := provider.EndpointOption{
		Endpoint: "http://8.147.234.89/api/system/login",
		Timeout:  10,
		HTTP: provider.Ehttp{
			Method: provider.PostHTTPMethod,
			Header: map[string]string{},
			Body:   tools.JsonMarshal(user),
		},
	}

	per, err := provider.NewEndpointHTTPer().Pilot(buildOption)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(tools.JsonMarshal(per))
}
