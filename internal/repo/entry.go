package repo

import (
	"watchAlert/internal/global"
	"watchAlert/pkg/client"

	"gorm.io/gorm"
)

type (
	entryRepo struct {
		g  InterGormDBCli
		db *gorm.DB
	}

	InterEntryRepo interface {
		DB() *gorm.DB
		Dashboard() InterDashboardRepo
		Tenant() InterTenantRepo
		AuditLog() InterAuditLogRepo
		Datasource() InterDatasourceRepo
		Duty() InterDutyRepo
		DutyCalendar() InterDutyCalendar
		Event() InterEventRepo
		Notice() InterNoticeRepo
		NoticeTmpl() InterNoticeTmplRepo
		Rule() InterRuleRepo
		RuleGroup() InterRuleGroupRepo
		RuleTmpl() InterRuleTmplRepo
		RuleTmplGroup() InterRuleTmplGroupRepo
		Silence() InterSilenceRepo
		User() InterUserRepo
		UserRole() InterUserRoleRepo
		UserPermissions() InterUserPermissionsRepo
		Setting() InterSettingRepo
		Subscribe() InterSubscribeRepo
		Probing() InterProbingRepo
		FaultCenter() InterFaultCenterRepo
		Ai() InterAiRepo
		Comment() InterCommentRepo
		Topology() InterTopologyRepo
	}
)

func NewRepoEntry() InterEntryRepo {
	dbConfig := global.Config.Database
	db := client.NewDBClient(client.DBConfig{
		Type:    dbConfig.Type,
		Host:    dbConfig.Host,
		Port:    dbConfig.Port,
		User:    dbConfig.User,
		Pass:    dbConfig.Pass,
		DBName:  dbConfig.DBName,
		Timeout: dbConfig.Timeout,
		Path:    dbConfig.Path,
	})
	g := NewInterGormDBCli(db)
	return &entryRepo{
		g:  g,
		db: db,
	}
}

func (e *entryRepo) DB() *gorm.DB                    { return e.db }
func (e *entryRepo) Dashboard() InterDashboardRepo   { return newDashboardInterface(e.db, e.g) }
func (e *entryRepo) Tenant() InterTenantRepo         { return newTenantInterface(e.db, e.g) }
func (e *entryRepo) AuditLog() InterAuditLogRepo     { return newAuditLogInterface(e.db, e.g) }
func (e *entryRepo) Datasource() InterDatasourceRepo { return newDatasourceInterface(e.db, e.g) }
func (e *entryRepo) Duty() InterDutyRepo             { return newDutyInterface(e.db, e.g) }
func (e *entryRepo) DutyCalendar() InterDutyCalendar { return newDutyCalendarInterface(e.db, e.g) }
func (e *entryRepo) Event() InterEventRepo           { return newEventInterface(e.db, e.g) }
func (e *entryRepo) Notice() InterNoticeRepo         { return newNoticeInterface(e.db, e.g) }
func (e *entryRepo) NoticeTmpl() InterNoticeTmplRepo { return newNoticeTmplInterface(e.db, e.g) }
func (e *entryRepo) Rule() InterRuleRepo             { return newRuleInterface(e.db, e.g) }
func (e *entryRepo) RuleGroup() InterRuleGroupRepo   { return newRuleGroupInterface(e.db, e.g) }
func (e *entryRepo) RuleTmpl() InterRuleTmplRepo     { return newRuleTmplInterface(e.db, e.g) }
func (e *entryRepo) RuleTmplGroup() InterRuleTmplGroupRepo {
	return newRuleTmplGroupInterface(e.db, e.g)
}
func (e *entryRepo) Silence() InterSilenceRepo   { return newSilenceInterface(e.db, e.g) }
func (e *entryRepo) User() InterUserRepo         { return newUserInterface(e.db, e.g) }
func (e *entryRepo) UserRole() InterUserRoleRepo { return newUserRoleInterface(e.db, e.g) }
func (e *entryRepo) UserPermissions() InterUserPermissionsRepo {
	return newInterUserPermissionsRepo(e.db, e.g)
}
func (e *entryRepo) Setting() InterSettingRepo         { return newSettingRepoInterface(e.db, e.g) }
func (e *entryRepo) Subscribe() InterSubscribeRepo     { return newInterSubscribeRepo(e.db, e.g) }
func (e *entryRepo) Probing() InterProbingRepo         { return newProbingRepoInterface(e.db, e.g) }
func (e *entryRepo) FaultCenter() InterFaultCenterRepo { return newInterFaultCenterRepo(e.db, e.g) }
func (e *entryRepo) Ai() InterAiRepo                   { return newAiRepoInterface(e.db, e.g) }
func (e *entryRepo) Comment() InterCommentRepo         { return newCommentInterface(e.db, e.g) }
func (e *entryRepo) Topology() InterTopologyRepo       { return newInterTopologyRepo(e.db, e.g) }
