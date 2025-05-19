package main

import (
	"fmt"
	"os"
	"path"
	"text/template"
)

const t = `name: Build and test TYPO3 {{ .Destination }} image

on:
  push:
    paths:
      - '{{ .Destination }}/Dockerfile'
      - '{{ .Destination }}/docker-compose.yml'
  pull_request:
    paths:
      - '{{ .Destination }}/Dockerfile'
      - '{{ .Destination }}/docker-compose.yml'

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v2

      - name: Build and start image
        run: |
          cd {{ .Destination }}
          docker compose build
          docker compose up -d

          sleep 20

          curl -vLf http://localhost
          docker compose down -v
`

func updateWorkflowFile(v *UpdateVersionSpec) (string, error) {
	tmpl, err := template.New("").Parse(t)
	if err != nil {
		return "", err
	}

	outFileName := fmt.Sprintf("verify-%s.yml", v.Destination)
	outFilePath := path.Join(".github", "workflows", outFileName)

	outFile, err := os.Create(outFilePath)
	if err != nil {
		return "", err
	}

	defer outFile.Close()

	if err := tmpl.Execute(outFile, &v); err != nil {
		return "", err
	}

	return outFilePath, nil
}
