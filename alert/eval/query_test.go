package eval

import (
	"reflect"
	"testing"
	"watchAlert/internal/models"
)

func TestSelectPrometheusRulesReturnsDefaultRulesWhenModeIsEmpty(t *testing.T) {
	defaultRules := []models.Rules{
		{Severity: "P0", Expr: "> 95", ForDuration: 300},
	}

	got := selectPrometheusRules(models.PrometheusConfig{
		Rules: defaultRules,
	}, map[string]interface{}{
		"instance": "node-a:9100",
	})

	if !reflect.DeepEqual(got, defaultRules) {
		t.Fatalf("expected default rules, got %#v", got)
	}
}

func TestSelectPrometheusRulesReturnsOverrideRulesWhenLabelsMatch(t *testing.T) {
	defaultRules := []models.Rules{
		{Severity: "P0", Expr: "> 95", ForDuration: 300},
	}
	overrideRules := []models.Rules{
		{Severity: "P0", Expr: "> 98", ForDuration: 600},
	}

	got := selectPrometheusRules(models.PrometheusConfig{
		Rules:         defaultRules,
		ThresholdMode: ThresholdModeNodeOverride,
		ThresholdOverrides: []models.ThresholdOverride{
			{
				MatchLabels: map[string]string{"instance": "node-a:9100"},
				Rules:       overrideRules,
			},
		},
	}, map[string]interface{}{
		"instance": "node-a:9100",
	})

	if !reflect.DeepEqual(got, overrideRules) {
		t.Fatalf("expected override rules, got %#v", got)
	}
}

func TestSelectPrometheusRulesFallsBackToDefaultRulesWhenLabelsDoNotMatch(t *testing.T) {
	defaultRules := []models.Rules{
		{Severity: "P0", Expr: "> 95", ForDuration: 300},
	}
	overrideRules := []models.Rules{
		{Severity: "P0", Expr: "> 98", ForDuration: 600},
	}

	got := selectPrometheusRules(models.PrometheusConfig{
		Rules:         defaultRules,
		ThresholdMode: ThresholdModeNodeOverride,
		ThresholdOverrides: []models.ThresholdOverride{
			{
				MatchLabels: map[string]string{"instance": "node-a:9100"},
				Rules:       overrideRules,
			},
		},
	}, map[string]interface{}{
		"instance": "node-b:9100",
	})

	if !reflect.DeepEqual(got, defaultRules) {
		t.Fatalf("expected default rules, got %#v", got)
	}
}

func TestSelectPrometheusRulesSkipsEmptyOverrideRules(t *testing.T) {
	defaultRules := []models.Rules{
		{Severity: "P0", Expr: "> 95", ForDuration: 300},
	}

	got := selectPrometheusRules(models.PrometheusConfig{
		Rules:         defaultRules,
		ThresholdMode: ThresholdModeNodeOverride,
		ThresholdOverrides: []models.ThresholdOverride{
			{
				MatchLabels: map[string]string{"instance": "node-a:9100"},
			},
		},
	}, map[string]interface{}{
		"instance": "node-a:9100",
	})

	if !reflect.DeepEqual(got, defaultRules) {
		t.Fatalf("expected default rules, got %#v", got)
	}
}

func TestMatchMetricLabelsRequiresAllLabelsToMatch(t *testing.T) {
	metricLabels := map[string]interface{}{
		"instance": "node-a:9100",
		"os_type":  "Linux",
	}

	if !matchMetricLabels(metricLabels, map[string]string{
		"instance": "node-a:9100",
		"os_type":  "Linux",
	}) {
		t.Fatal("expected all labels to match")
	}

	if matchMetricLabels(metricLabels, map[string]string{
		"instance": "node-a:9100",
		"os_type":  "Windows",
	}) {
		t.Fatal("expected mismatched label to fail")
	}
}
