package types

import (
	"crypto/md5"
	"encoding/hex"
	"time"
	"watchAlert/pkg/tools"
)

var NamespaceMetricsMap = map[string][]string{
	"AWS/EBS":     {"BurstBalance", "VolumeConsumedReadWriteOps", "VolumeIdleTime", "VolumeQueueLength", "VolumeReadBytes", "VolumeReadOps", "VolumeThroughputPercentage", "VolumeTotalReadTime", "VolumeTotalWriteTime", "VolumeWriteBytes", "VolumeWriteOps"},
	"AWS/EC2":     {"CPUCreditBalance", "CPUCreditUsage", "CPUSurplusCreditBalance", "CPUSurplusCreditsCharged", "CPUUtilization", "DiskReadBytes", "DiskReadOps", "DiskWriteBytes", "DiskWriteOps", "EBSByteBalance%", "EBSIOBalance%", "EBSReadBytes", "EBSReadOps", "EBSWriteBytes", "EBSWriteOps", "MetadataNoToken", "NetworkIn", "NetworkOut", "NetworkPacketsIn", "NetworkPacketsOut", "StatusCheckFailed", "StatusCheckFailed_Instance", "StatusCheckFailed_System"},
	"AWS/RDS":     {"ActiveTransactions", "AuroraBinlogReplicaLag", "AuroraGlobalDBDataTransferBytes", "AuroraGlobalDBReplicatedWriteIO", "AuroraGlobalDBReplicationLag", "AuroraReplicaLag", "AuroraReplicaLagMaximum", "AuroraReplicaLagMinimum", "AvailabilityPercentage", "BacktrackChangeRecordsCreationRate", "BacktrackChangeRecordsStored", "BacktrackWindowActual", "BacktrackWindowAlert", "BackupRetentionPeriodStorageUsed", "BinLogDiskUsage", "BlockedTransactions", "BufferCacheHitRatio", "BurstBalance", "CPUCreditBalance", "CPUCreditUsage", "CPUUtilization", "ClientConnections", "ClientConnectionsClosed", "ClientConnectionsNoTLS", "ClientConnectionsReceived", "ClientConnectionsSetupFailedAuth", "ClientConnectionsSetupSucceeded", "ClientConnectionsTLS", "CommitLatency", "CommitThroughput", "DDLLatency", "DDLThroughput", "DMLLatency", "DMLThroughput", "DatabaseConnectionRequests", "DatabaseConnectionRequestsWithTLS", "DatabaseConnections", "DatabaseConnectionsBorrowLatency", "DatabaseConnectionsCurrentlyBorrowed", "DatabaseConnectionsCurrentlyInTransaction", "DatabaseConnectionsCurrentlySessionPinned", "DatabaseConnectionsSetupFailed", "DatabaseConnectionsSetupSucceeded", "DatabaseConnectionsWithTLS", "Deadlocks", "DeleteLatency", "DeleteThroughput", "DiskQueueDepth", "EBSByteBalance%", "EBSIOBalance%", "EngineUptime", "FailedSQLServerAgentJobsCount", "FreeLocalStorage", "FreeStorageSpace", "FreeableMemory", "InsertLatency", "InsertThroughput", "LoginFailures", "MaxDatabaseConnectionsAllowed", "MaximumUsedTransactionIDs", "NetworkReceiveThroughput", "NetworkThroughput", "NetworkTransmitThroughput", "OldestReplicationSlotLag", "Queries", "QueryDatabaseResponseLatency", "QueryRequests", "QueryRequestsNoTLS", "QueryRequestsTLS", "QueryResponseLatency", "RDSToAuroraPostgreSQLReplicaLag", "ReadIOPS", "ReadLatency", "ReadThroughput", "ReplicaLag", "ReplicationSlotDiskUsage", "ResultSetCacheHitRatio", "SelectLatency", "SelectThroughput", "ServerlessDatabaseCapacity", "SnapshotStorageUsed", "SwapUsage", "TotalBackupStorageBilled", "TransactionLogsDiskUsage", "TransactionLogsGeneration", "UpdateLatency", "UpdateThroughput", "VolumeBytesUsed", "VolumeReadIOPs", "VolumeWriteIOPs", "WriteIOPS", "WriteLatency", "WriteThroughput"},
	"AWS/Route53": {"ChildHealthCheckHealthyCount", "ConnectionTime", "DNSQueries", "HealthCheckPercentageHealthy", "HealthCheckStatus", "SSLHandshakeTime", "TimeToFirstByte"},
	"AWS/S3":      {"4xxErrors", "5xxErrors", "AllRequests", "BucketSizeBytes", "BytesDownloaded", "BytesUploaded", "DeleteRequests", "FirstByteLatency", "GetRequests", "HeadRequests", "ListRequests", "NumberOfObjects", "PostRequests", "PutRequests", "SelectRequests", "SelectReturnedBytes", "SelectScannedBytes", "TotalRequestLatency"},
	"AWS/SES":     {"Bounce", "Clicks", "Complaint", "Delivery", "Opens", "Reject", "Rendering Failures", "Reputation.BounceRate", "Reputation.ComplaintRate", "Send"},
	"AWS/SNS":     {"NumberOfMessagesPublished", "NumberOfNotificationsDelivered", "NumberOfNotificationsFailed", "NumberOfNotificationsFilteredOut", "NumberOfNotificationsFilteredOut-InvalidAttributes", "NumberOfNotificationsFilteredOut-NoMessageAttributes", "PublishSize", "SMSMonthToDateSpentUSD", "SMSSuccessRate"},
}

var NamespaceDimensionKeysMap = map[string][]string{
	"AWS/EBS":     {"VolumeId"},
	"AWS/EC2":     {"AutoScalingGroupName", "ImageId", "InstanceId", "InstanceType"},
	"AWS/RDS":     {"DBClusterIdentifier", "DBInstanceIdentifier"},
	"AWS/Route53": {"HealthCheckId", "Region", "HostedZoneId"},
	"AWS/S3":      {"BucketName", "FilterId", "StorageType"},
	"AWS/SES":     {},
	"AWS/SNS":     {"Application", "Country", "Platform", "SMSType", "TopicName"},
}

type CloudWatchQuery struct {
	Endpoint   string    `json:"endpoint"`
	Dimension  string    `json:"dimension"`
	Namespace  string    `json:"namespace"`
	MetricName string    `json:"metricName"`
	Statistic  string    `json:"statistic"`
	Period     int32     `json:"period"`
	Form       time.Time `json:"form"`
	To         time.Time `json:"to"`
}

func (c CloudWatchQuery) GetFingerprint() string {
	newMetric := map[string]interface{}{
		"namespace":  c.Namespace,
		"metricName": c.MetricName,
		"statistic":  c.Statistic,
	}
	h := md5.New()
	streamString := tools.JsonMarshalToString(newMetric)
	h.Write([]byte(streamString))
	fingerprint := hex.EncodeToString(h.Sum(nil))

	return fingerprint
}

func (c CloudWatchQuery) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"instance":   c.Endpoint,
		"namespace":  c.Namespace,
		"metricName": c.MetricName,
		"statistic":  c.Statistic,
	}
}

type MetricNamesQuery struct {
	MetricType string `json:"metricType" form:"metricType"`
}

type RdsInstanceReq struct {
	DatasourceId string `json:"datasourceId" form:"datasourceId"`
}

type RdsClusterReq struct {
	DatasourceId string `json:"datasourceId" form:"datasourceId"`
}

type RdsDimensionReq struct {
	MetricType string `json:"metricType" form:"metricType"`
}
