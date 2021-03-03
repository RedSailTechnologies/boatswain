{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "boatswain.mongodbConn" -}}
- name: MONGODB_ROOT_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ .Values.global.mongodb.auth.existingSecret }}
      key: mongodb-root-password
- name: MONGO_CONNECTION_STRING
  value: {{ printf "mongodb://root:$(MONGODB_ROOT_PASSWORD)@%s-mongodb:%s" .Release.Name (toString .Values.global.mongodb.service.port) }}
{{- end }}
