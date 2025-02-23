package provider

import (
	"fmt"
	"testing"
	"time"
	"watchAlert/internal/models"
)

func TestElasticsearch_GetIndexName(t *testing.T) {
	type fields struct {
		IndexOption models.EsIndexOption
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "testIndex-YYYYMMdd",
			fields: struct{ IndexOption models.EsIndexOption }{
				IndexOption: models.EsIndexOption{Index: "testIndex", WithDate: true, DatePattern: "YYYYMMdd", Separator: "-"},
			},
			want: fmt.Sprintf("testIndex-%v", time.Now().Format("20060102")),
		},
		{
			name: "testIndexYYYYMMdd",
			fields: struct{ IndexOption models.EsIndexOption }{
				IndexOption: models.EsIndexOption{Index: "testIndex", WithDate: true, DatePattern: "YYYYMMdd", Separator: ""},
			},
			want: fmt.Sprintf("testIndex%v", time.Now().Format("20060102")),
		},
		{
			name: "testIndex-2025.02.08",
			fields: struct{ IndexOption models.EsIndexOption }{
				IndexOption: models.EsIndexOption{Index: "testIndex", WithDate: true, DatePattern: "YYYY.MM.dd", Separator: "-"},
			},
			want: fmt.Sprintf("testIndex-%v", time.Now().Format("2006.01.02")),
		},
		{
			name: "testIndex-2025-02-08",
			fields: struct{ IndexOption models.EsIndexOption }{
				IndexOption: models.EsIndexOption{Index: "testIndex", WithDate: true, DatePattern: "YYYY-MM-dd", Separator: "-"},
			},
			want: fmt.Sprintf("testIndex-%v", time.Now().Format("2006-01-02")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Elasticsearch{
				IndexOption: tt.fields.IndexOption,
			}
			if got := e.GetIndexName(); got != tt.want {
				t.Errorf("GetIndexName() = %v, want %v", got, tt.want)
			}
		})
	}
}
