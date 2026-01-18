package providers

import (
	"bytes"
	"os/exec"
	"strings"
	"text/template"
)

type CursorCli struct {
	cmd string
}

func NewCursorCli() *CursorCli {
	return &CursorCli{
		cmd: "cursor-agent {{.Request}}",
	}
}

func (o *CursorCli) Run(request string) string {
	tmpl, err := template.New("cmd").Parse(o.cmd)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct{ Request string }{Request: request}); err != nil {
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

	return string(stdout)
}
