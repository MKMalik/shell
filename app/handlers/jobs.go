package handlers

import (
	"fmt"
	"strings"
)

func HandleJobs(_ string) string {
	var b strings.Builder

	for i, job := range JobList {
		marker := " "

		if i == len(JobList)-1 {
			marker = "+"
		} else if i == len(JobList)-2 {
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

	// Reap jobs that were already shown as Done
	filtered := JobList[:0]

	for _, job := range JobList {
		if job.Status == "Done" {
			continue
		}
		filtered = append(filtered, job)
	}

	JobList = filtered

	return b.String()
}
