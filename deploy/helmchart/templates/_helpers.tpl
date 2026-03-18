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

{{/* MySQL Database Host */}}
{{- define "watchalert.database.host" -}}
{{- if .Values.externalDatabase.enabled -}}
{{- required "externalDatabase.mysql.host is required when externalDatabase.enabled=true" .Values.externalDatabase.mysql.host -}}
{{- else if .Values.mysql.enabled -}}
{{- printf "%s-mysql" .Release.Name -}}
{{- else -}}
{{- fail "Either mysql.enabled=true or externalDatabase.enabled=true must be set" -}}
{{- end -}}
{{- end -}}

{{/* MySQL Database Port */}}
{{- define "watchalert.database.port" -}}
{{- if .Values.externalDatabase.enabled -}}
{{- .Values.externalDatabase.mysql.port | default 3306 -}}
{{- else if .Values.mysql.enabled -}}
3306
{{- else -}}
{{- fail "Either mysql.enabled=true or externalDatabase.enabled=true must be set" -}}
{{- end -}}
{{- end -}}

{{/* MySQL Database User */}}
{{- define "watchalert.database.user" -}}
{{- if .Values.externalDatabase.enabled -}}
{{- .Values.externalDatabase.mysql.user | default "root" -}}
{{- else if .Values.mysql.enabled -}}
root
{{- else -}}
{{- fail "Either mysql.enabled=true or externalDatabase.enabled=true must be set" -}}
{{- end -}}
{{- end -}}

{{/* MySQL Database Password */}}
{{- define "watchalert.database.password" -}}
{{- if .Values.externalDatabase.enabled -}}
{{- required "externalDatabase.mysql.password is required when externalDatabase.enabled=true" .Values.externalDatabase.mysql.password -}}
{{- else if .Values.mysql.enabled -}}
{{- .Values.mysql.auth.rootPassword -}}
{{- else -}}
{{- fail "Either mysql.enabled=true or externalDatabase.enabled=true must be set" -}}
{{- end -}}
{{- end -}}

{{/* MySQL Database Name */}}
{{- define "watchalert.database.name" -}}
{{- if .Values.externalDatabase.enabled -}}
{{- .Values.externalDatabase.mysql.database | default "watchalert" -}}
{{- else if .Values.mysql.enabled -}}
{{- .Values.mysql.auth.database -}}
{{- else -}}
{{- fail "Either mysql.enabled=true or externalDatabase.enabled=true must be set" -}}
{{- end -}}
{{- end -}}

{{/* Redis Host */}}
{{- define "watchalert.redis.host" -}}
{{- if .Values.externalDatabase.enabled -}}
{{- required "externalDatabase.redis.host is required when externalDatabase.enabled=true" .Values.externalDatabase.redis.host -}}
{{- else if .Values.redis.enabled -}}
{{- printf "%s-redis" .Release.Name -}}
{{- else -}}
{{- fail "Either redis.enabled=true or externalDatabase.enabled=true must be set" -}}
{{- end -}}
{{- end -}}

{{/* Redis Port */}}
{{- define "watchalert.redis.port" -}}
{{- if .Values.externalDatabase.enabled -}}
{{- .Values.externalDatabase.redis.port | default 6379 -}}
{{- else if .Values.redis.enabled -}}
6379
{{- else -}}
{{- fail "Either redis.enabled=true or externalDatabase.enabled=true must be set" -}}
{{- end -}}
{{- end -}}

{{/* Redis Password */}}
{{- define "watchalert.redis.password" -}}
{{- if .Values.externalDatabase.enabled -}}
{{- .Values.externalDatabase.redis.password | default "" -}}
{{- else if .Values.redis.enabled -}}
""
{{- else -}}
{{- fail "Either redis.enabled=true or externalDatabase.enabled=true must be set" -}}
{{- end -}}
{{- end -}}

{{/* Redis Database */}}
{{- define "watchalert.redis.database" -}}
{{- if .Values.externalDatabase.enabled -}}
{{- .Values.externalDatabase.redis.database | default 0 -}}
{{- else if .Values.redis.enabled -}}
0
{{- else -}}
{{- fail "Either redis.enabled=true or externalDatabase.enabled=true must be set" -}}
{{- end -}}
{{- end -}}

{{/*
MySQL Image for Job

This helper determines which MySQL client image to use for the initialization job:

- When using INTERNAL MySQL (mysql.enabled=true, externalDatabase.enabled=false):
  Uses mysql.image to ensure version compatibility between the StatefulSet and job.
  Example: If mysql.image.tag=8.4, the job will also use MySQL 8.4 client.

- When using EXTERNAL MySQL (externalDatabase.enabled=true):
  Uses job.image to allow flexibility in client version selection.
  Example: External database might be MySQL 8.0, so job.image.tag=8.0 is used.

This automatic selection ensures the MySQL client version matches the server version,
preventing compatibility issues during database initialization.
*/}}
{{- define "watchalert.job.mysql.image" -}}
{{- if and (not .Values.externalDatabase.enabled) .Values.mysql.enabled -}}
{{- printf "%s:%s" .Values.mysql.image.repository .Values.mysql.image.tag -}}
{{- else -}}
{{- printf "%s:%s" .Values.job.image.repository .Values.job.image.tag -}}
{{- end -}}
{{- end -}}
