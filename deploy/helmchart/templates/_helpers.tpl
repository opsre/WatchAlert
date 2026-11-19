{{/* Common helpers for watchalert chart */}}
{{- define "watchalert.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "watchalert.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "watchalert.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "watchalert.labels" -}}
{{ include "watchalert.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "watchalert.selectorLabels" -}}
app.kubernetes.io/name: {{ include "watchalert.name" . }}
{{- end }}

{{- define "watchalert.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "watchalert.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/* Build a component-scoped name like <fullname>-<component> */}}
{{- define "watchalert.componentName" -}}
{{- printf "%s-%s" .root.Release.Name .component | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* Resolve the MySQL endpoint used by both the application and the init job. */}}
{{- define "watchalert.mysqlHost" -}}
{{- if .Values.mysql.enabled -}}
{{- include "watchalert.componentName" (dict "root" . "component" "mysql") -}}
{{- else -}}
{{- required "mysql.host must be set when mysql.enabled is false" .Values.mysql.host -}}
{{- end -}}
{{- end -}}

{{- define "watchalert.mysqlPort" -}}
{{- .Values.mysql.port | default 3306 -}}
{{- end -}}

{{/* Resolve the Redis endpoint used by the application and its wait container. */}}
{{- define "watchalert.redisHost" -}}
{{- if .Values.redis.enabled -}}
{{- include "watchalert.componentName" (dict "root" . "component" "redis") -}}
{{- else -}}
{{- required "redis.host must be set when redis.enabled is false" .Values.redis.host -}}
{{- end -}}
{{- end -}}

{{- define "watchalert.redisPort" -}}
{{- .Values.redis.port | default 6379 -}}
{{- end -}}
