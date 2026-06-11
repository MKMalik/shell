package handlers

import (
	"fmt"
	"strings"
)

func HandleJobs(_ string) string {
	var running []Job

	for _, job := range JobList {
		if job.Status == "Running" {
			running = append(running, job)
		}
	}

	if len(running) == 0 {
		return ""
	}

	var b strings.Builder

	for i, job := range running {
		marker := " "

		switch {
		case i == len(running)-1:
			marker = "+"
		case i == len(running)-2:
			marker = "-"
		}

		fmt.Fprintf(
			&b,
			"[%d]%s  %-24s%s\n",
			job.ID,
			marker,
			job.Status,
			job.CommandString,
		)
	}

	return b.String()
}
