package main

import "text/template"

// This file hold the install/bootstrap script templates. used by the /api/v1/pki/token route

var powershellScriptInstaller = template.Must(template.New("name").Parse(`#ps1_sysnative
mkdir C:\\monsoon\\arc
(New-Object System.Net.WebClient).DownloadFile('{{ .UpdateURL }}/arc/windows/amd64/latest','C:\\monsoon\\arc\\arc.exe')
C:\\monsoon\\arc\\arc.exe init --endpoint {{ .EndpointURL }} --update-uri {{ .UpdateURL }} --registration-url {{ .SignURL }} {{if .ApiURL}}--api-uri {{ .ApiURL }} {{end}}
`))

var shellScriptInstaller = template.Must(template.New("name").Parse(`#!/bin/sh
curl -f --create-dirs -o /opt/arc/arc {{ .UpdateURL }}/arc/linux/amd64/latest
chmod +x /opt/arc/arc
/opt/arc/arc init --endpoint {{ .EndpointURL }} --update-uri {{ .UpdateURL }} --registration-url {{ .SignURL }} {{if .ApiURL}}--api-uri {{ .ApiURL }} {{end}}
`))

var cloudConfigInstaller = template.Must(template.New("name").Parse(`#cloud-config
runcmd:
  - - sh
    - -ec
    - |
      curl -f --create-dirs -o /opt/arc/arc {{ .UpdateURL }}/arc/linux/amd64/latest
      chmod +x /opt/arc/arc
      /opt/arc/arc init --endpoint {{ .EndpointURL }} --update-uri {{ .UpdateURL }} --registration-url {{ .SignURL }} {{if .ApiURL}}--api-uri {{ .ApiURL }} {{end}}
`))
