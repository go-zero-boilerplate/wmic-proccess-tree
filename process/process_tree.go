package process

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-zero-boilerplate/wmic-xml-command"
)

func LoadProcessTree(pid int) (*Process, error) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}
	wrapped := &Process{Process: proc}
	if err = wrapped.LoadChildren(); err != nil {
		return nil, err
	}
	return wrapped, nil
}

type Process struct {
	Process  *os.Process
	Children []*Process
}

func (p *Process) String() string {
	jsonBytes, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		//Use plain formatting
		return fmt.Sprintf("Pid: %d, Children:\n  %+v", p.Process.Pid, p.Children)
	}
	return string(jsonBytes)
}

func (p *Process) LoadChildren() error {
	columnNames := []string{
		"Caption",
		"ProcessId",
	}

	wmicArgs := []string{
		"process",
		"where",
		fmt.Sprintf("(ParentProcessId=%d)", p.Process.Pid),
		"get",
		strings.Join(columnNames, ","),
	}

	responseXml, err := wmic_command.Run(wmicArgs)
	if err != nil {
		parseErr, ok := err.(*wmic_command.XmlParseError)
		if !ok {
			return fmt.Errorf("Cannot find child processes. Error: %s", err.Error())
		}

		xmlTrimmed := strings.TrimSpace(string(parseErr.XmlContent))
		if xmlTrimmed == "" {
			return fmt.Errorf("Command returned empty response, cannot parse xml")
		}
		if !strings.HasPrefix(xmlTrimmed, "<") {
			if strings.HasPrefix(xmlTrimmed, "No Instance(s) Available") {
				//In this case the process does not have children processes
				return nil
			}
			return fmt.Errorf("Invalid xml returned from wmic command, xml was: %s", xmlTrimmed)
		}
	}

	for _, res := range responseXml.Results {
		var childProcId int64
		for _, prop := range res.Properties {
			if prop.Name == "ProcessId" {
				childProcId, err = strconv.ParseInt(prop.Value, 10, 32)
				if err != nil {
					return fmt.Errorf("Got process ID from wmic of '%s' but is invalid int, error: %s", prop.Value, err.Error())
				}
			}
		}

		if childProcId > 0 {
			childProc, err := os.FindProcess(int(childProcId))
			if err != nil {
				return fmt.Errorf("Cannot get os Process for process id %d (parent was %d), error: %s", childProcId, p.Process.Pid, err.Error())
			}
			wrappedChildProc := &Process{Process: childProc}

			if err = wrappedChildProc.LoadChildren(); err != nil {
				return fmt.Errorf("Cannot get child processes of pid %d, error: %s", p.Process.Pid, err.Error())
			}
			p.Children = append(p.Children, wrappedChildProc)
		}
	}

	return nil
}
