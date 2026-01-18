package providers

import (
	"bytes"
	"os/exec"
	"strings"
	"text/template"
)

type OpenCodeCli struct {
	systemPrompt string
	cmd          string
}

func NewOpenCodeCli() *OpenCodeCli {
	return &OpenCodeCli{
		systemPrompt: "Always return output in markdown format. Do not use any tools without explicit request. ",
		cmd:          "bunx opencode-ai run {{.Request}}",
	}
}

func (o *OpenCodeCli) Run(request string) string {
	tmpl, err := template.New("cmd").Parse(o.cmd)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct{ Request string }{Request: o.systemPrompt + request}); err != nil {
		return ""
	}

	cmdParts := strings.Fields(buf.String())
	if len(cmdParts) == 0 {
		return ""
	}

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	stdout, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(stdout))
}
