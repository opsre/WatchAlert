package cloudwatch

import (
	"context"
	"watchAlert/internal/ctx"
	"watchAlert/pkg/provider"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

type (
	awsRdsService struct {
		ctx *ctx.Context
	}

	InterAwsRdsService interface {
		GetDBInstanceIdentifier(req interface{}) (interface{}, interface{})
		GetDBClusterIdentifier(req interface{}) (interface{}, interface{})
	}
)

func NewInterAWSRdsService(ctx *ctx.Context) InterAwsRdsService {
	return awsRdsService{
		ctx: ctx,
	}
}

func (a awsRdsService) GetDBInstanceIdentifier(req interface{}) (interface{}, interface{}) {
	r := req.(*RdsInstanceReq)
	datasourceObj, err := a.ctx.DB.Datasource().GetInstance(r.DatasourceId)
	if err != nil {
		return nil, err
	}

	cfg, err := provider.NewAWSCredentialCfg(datasourceObj.AWSCloudWatch.Region, datasourceObj.AWSCloudWatch.AccessKey, datasourceObj.AWSCloudWatch.SecretKey, datasourceObj.Labels)
	if err != nil {
		return nil, err
	}

	cli := cfg.RdsCli()
	input := &rds.DescribeDBInstancesInput{}
	result, err := cli.DescribeDBInstances(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var instances []string
	for _, instance := range result.DBInstances {
		instances = append(instances, *instance.DBInstanceIdentifier)
	}

	return instances, nil
}

func (a awsRdsService) GetDBClusterIdentifier(req interface{}) (interface{}, interface{}) {
	r := req.(*RdsClusterReq)
	datasourceObj, err := a.ctx.DB.Datasource().GetInstance(r.DatasourceId)
	if err != nil {
		return nil, err
	}

	cfg, err := provider.NewAWSCredentialCfg(datasourceObj.AWSCloudWatch.Region, datasourceObj.AWSCloudWatch.AccessKey, datasourceObj.AWSCloudWatch.SecretKey, datasourceObj.Labels)
	if err != nil {
		return nil, err
	}

	cli := cfg.RdsCli()
	input := &rds.DescribeDBClustersInput{}
	result, err := cli.DescribeDBClusters(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var clusters []string
	for _, cluster := range result.DBClusters {
		clusters = append(clusters, *cluster.DBClusterIdentifier)
	}

	return clusters, nil
}
