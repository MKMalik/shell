package handlers

import "fmt"

func HandleJobs(_ string) string {
	for _, job := range JobList {
		if job.Status == "Running" {
			return fmt.Sprintf(
				"[%d]+  %-24s%s\n",
				job.ID,
				job.Status,
				job.CommandString,
			)
		}
	}

	return ""
}
