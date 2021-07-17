{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "warm-images.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "warm-images.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "warm-images.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "warm-images.labels" -}}
helm.sh/chart: {{ include "warm-images.chart" . }}
{{ include "warm-images.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "warm-images.selectorLabels" -}}
app.kubernetes.io/name: {{ include "warm-images.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "warm-images.serviceAccountName" -}}
{{ printf "%s-%s" .Release.Name "svc-acc" }}
{{- end -}}

{{/*
Create the name of the cluster role binding to use
*/}}
{{- define "warm-images.clusterRoleBindingName" -}}
{{ printf "%s-%s" .Release.Name "clusterrole-binding" }}
{{- end -}}

{{/*
Create the name of the Controller
*/}}
{{- define "warm-images.controllerName" -}}
{{ printf "%s-%s" .Release.Name "controller" }}
{{- end -}}

{{/*
Create the name of the Agent
*/}}
{{- define "warm-images.agentName" -}}
{{ printf "%s-%s" .Release.Name "agent" }}
{{- end -}}

{{/*
Create the name of the ConfigMap
*/}}
{{- define "warm-images.configmap" -}}
{{ printf "%s-%s" .Release.Name "cm" }}
{{- end -}}

{{/*
Custom registry name
*/}}
{{- define "warm-images.customRegistry" -}}
{{- if .Values.image.customRegistry -}}
{{ printf "%s/" .Values.image.customRegistry -}}
{{- else -}}
{{ print "" }}
{{- end -}}
{{- end -}}
