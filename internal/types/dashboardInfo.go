package types

type ResponseDashboardInfo struct {
	CountAlertRules   int64             `json:"countAlertRules"`
	FaultCenterNumber int64             `json:"faultCenterNumber"`
	UserNumber        int64             `json:"userNumber"`
	CurAlertList      []AlertList       `json:"curAlertList"`
	AlarmDistribution AlarmDistribution `json:"alarmDistribution"`
}

type AlarmDistribution struct {
	P0 int64 `json:"P0"`
	P1 int64 `json:"P1"`
	P2 int64 `json:"P2"`
}

type AlertList struct {
	RuleName      string `json:"ruleName"`
	Severity      string `json:"severity"`
	FaultCenterId string `json:"faultCenterId"`
}
