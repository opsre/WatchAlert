package types

type ResponseDashboardInfo struct {
	CountAlertRules   int64       `json:"countAlertRules"`
	FaultCenterNumber int64       `json:"faultCenterNumber"`
	UserNumber        int64       `json:"userNumber"`
	CurAlertList      []AlertList `json:"curAlertList"`
}
type AlertList struct {
	RuleName      string `json:"ruleName"`
	Severity      string `json:"severity"`
	FaultCenterId string `json:"faultCenterId"`
	TiggerTime    int64  `json:"tiggerTime"`
}
