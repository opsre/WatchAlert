package test

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
	"watchAlert/pkg/provider"
)

func TestKubernetesClient(t *testing.T) {
	cli, err := provider.NewKubernetesClient(context.Background(), "", nil)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	event, err := cli.GetWarningEvent("", 10)
	if err != nil {
		logrus.Errorf(err.Error())
		return
	}

	fmt.Println(event)
}
