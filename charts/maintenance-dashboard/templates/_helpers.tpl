{{/*
Expand the name of the chart.
*/}}
{{- define "maintenance-dashboard.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "maintenance-dashboard.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "maintenance-dashboard.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "maintenance-dashboard.labels" -}}
helm.sh/chart: {{ include "maintenance-dashboard.chart" . }}
{{ include "maintenance-dashboard.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app: maintenance-dashboard
{{- end }}

{{/*
Selector labels
*/}}
{{- define "maintenance-dashboard.selectorLabels" -}}
app.kubernetes.io/name: {{ include "maintenance-dashboard.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "maintenance-dashboard.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "maintenance-dashboard.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Return true if a secret object should be created
*/}}
{{- define "maintenance-dashboard.createSecret" -}}
{{- if not .Values.existingSecret -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Get the password secret.
*/}}
{{- define "maintenance-dashboard.secretName" -}}
{{- if .Values.existingSecret }}
    {{- printf "%s" (tpl .Values.existingSecret $) -}}
{{- else -}}
    {{- printf "%s" (include "maintenance-dashboard.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Get kube-prometheus-stack release name for the ServiceMonitor.
*/}}
{{- define "maintenance-dashboard.kube-prometheus-stack-release" -}}
{{- if .Values.kubePrometheusStackReleaseName }}
    {{- printf "%s" (tpl .Values.kubePrometheusStackReleaseName $) -}}
{{- else -}}
    {{- printf "%s-%s" .Release.Name "kube-prometheus-stack" | trunc 63 | trimSuffix "-" }}
{{- end -}}
{{- end -}}
