package cloudwatch

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	log "github.com/sirupsen/logrus"
)

func MetricDataQuery(client *cloudwatch.Client, query CloudWatchQuery) ([]time.Time, []float64) {
	input := &cloudwatch.GetMetricDataInput{
		MetricDataQueries: []types.MetricDataQuery{
			{
				Id: aws.String("query"),
				MetricStat: &types.MetricStat{
					Metric: &types.Metric{
						Dimensions: []types.Dimension{
							{
								Name:  aws.String(query.Dimension),
								Value: aws.String(query.Endpoint),
							},
						},
						MetricName: aws.String(query.MetricName),
						Namespace:  aws.String(query.Namespace),
					},
					Stat:   aws.String(query.Statistic),
					Period: aws.Int32(query.Period),
				},
			},
		},
		StartTime: aws.Time(query.Form),
		EndTime:   aws.Time(query.To),
	}
	output, err := client.GetMetricData(context.TODO(), input)
	if err != nil {
		log.Errorf(err.Error())
		return nil, nil
	}

	var times []time.Time
	var values []float64
	for _, result := range output.MetricDataResults {
		for k, value := range result.Values {
			times = append(times, result.Timestamps[k])
			values = append(values, value)
		}
	}

	return times, values
}
