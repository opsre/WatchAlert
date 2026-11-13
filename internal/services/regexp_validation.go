package services

import (
	"fmt"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

func validateRegexpLabel(operator, value string) error {
	if operator != "=~" && operator != "!~" {
		return nil
	}
	if err := tools.ValidateRegexp(value); err != nil {
		return fmt.Errorf("无效的正则表达式 %q: %w", value, err)
	}
	return nil
}

func validateNoticeRoutes(routes []models.NoticeRoute) error {
	for routeIndex, route := range routes {
		for labelIndex, label := range route.NoticeLabels {
			if err := validateRegexpLabel(label.Operator, label.Value); err != nil {
				return fmt.Errorf("通知路由 %d 标签 %d: %w", routeIndex+1, labelIndex+1, err)
			}
		}
	}
	return nil
}

func validateSilenceLabels(labels []models.SilenceLabel) error {
	for labelIndex, label := range labels {
		if err := validateRegexpLabel(label.Operator, label.Value); err != nil {
			return fmt.Errorf("静默规则标签 %d: %w", labelIndex+1, err)
		}
	}
	return nil
}
