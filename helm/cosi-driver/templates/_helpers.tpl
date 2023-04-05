{{/*
Expand the name of the chart.
*/}}
{{- define "cosi-driver.name" }}
  {{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cosi-driver.fullname" }}
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
{{- define "cosi-driver.chart" }}
  {{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
# COSI driver log level
# Possible values: "trace" "debug" "info" "warn" "error" "fatal" "panic"
# Default value: "debug"
*/}}
{{- define "cosi-driver.logLevel" }}
  {{- $logLevelValues := list "trace" "debug" "info" "warn" "error" "fatal" "panic" }}
  {{- if (has .Values.logLevel $logLevelValues) }}
    {{- .Values.provisioner.logLevel }}
  {{- else }}
    {{- "debug" }}
  {{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cosi-driver.labels" }}
helm.sh/chart: {{ include "cosi-driver.chart" . }}
{{- include "cosi-driver.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cosi-driver.selectorLabels" }}
app.kubernetes.io/name: {{ include "cosi-driver.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the role to use
*/}}
{{- define "cosi-driver.roleName" }}
  {{- if and .Values.rbac.create .Values.rbac.role.create }}
      {{- default (printf "%s-role" (include "cosi-driver.fullname" .)) .Values.rbac.role.name }}
  {{- else }}
    {{- .Values.rbac.role.name }}
  {{- end }}
{{- end }}

{{/*
Create the name of the role binding to use
*/}}
{{- define "cosi-driver.roleBindingName" }}
  {{- if and .Values.rbac.create .Values.rbac.roleBinding.create }}
    {{- default (printf "%s-rolebinding" (include "cosi-driver.fullname" .)) .Values.rbac.roleBinding.name }}
  {{- else }}
    {{- .Values.rbac.roleBinding.name }}
  {{- end }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "cosi-driver.serviceAccountName" }}
  {{- if and .Values.rbac.create .Values.rbac.serviceAccount.create }}
    {{- default (printf "%s-sa" (include "cosi-driver.fullname" .)) .Values.rbac.serviceAccount.name }}
  {{- else }}
    {{- .Values.rbac.serviceAccount.name }}
  {{- end }}
{{- end }}

{{/*
Create the name of provisioner container
*/}}
{{- define "cosi-driver.provisionerContainerName" }}
  {{- default "cosi-provisioner" .Values.provisioner.name }}
{{- end }}

{{/*
Create the name of provisioner sidecar container
*/}}
{{- define "cosi-driver.provisionerSidecarContainerName" }}
  {{- default "cosi-provisioner-sidecar" .Values.sidecar.name }}
{{- end }}

{{/*
Create the full name of provisioner image from repository and tag
*/}}
{{- define "cosi-driver.provisionerImageName" }}
  {{- .Values.provisioner.image.repository }}:{{ .Values.provisioner.image.tag | default .Chart.AppVersion }}
{{- end }}

{{/*
Create the full name of provisioner sidecar image from repository and tag
*/}}
{{- define "cosi-driver.provisionerSidecarImageName" }}
  {{- .Values.sidecar.image.repository }}:{{ .Values.sidecar.image.tag }}
{{- end }}

{{/*
Create the secret name
*/}}
{{- define "cosi-driver.secretName" }}
  {{- if .Values.configuration.create }}
    {{- default (printf "%s-secret" (include "cosi-driver.name" . )) .Values.configuration.name }}
  {{- else }}
    {{- .Values.configuration.name }}
  {{- end }}
{{- end }}

{{/*
Create the name for secret volume
*/}}
{{- define "cosi-driver.secretVolumeName" }}
  {{- printf "%s-secret-volume" (include "cosi-driver.name" . ) }}
{{- end }}
