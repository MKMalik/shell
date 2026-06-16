package handlers

import "sync"

var Builtins = map[Builtin]func(string) string{
	Echo:     HandleEcho,
	Type:     HandleType,
	Exit:     func(s string) string { return "" },
	Pwd:      HandlePwd,
	Cd:       HandleCd,
	Complete: HandleComplete,
	Jobs:     HandleJobs,
}

var BuiltinNames = map[Builtin]struct{}{}

func IsBuiltin(cmd string) bool {
	_, ok := BuiltinNames[Builtin(cmd)]
	return ok
}

func init() {
	for builtin := range Builtins {
		BuiltinNames[builtin] = struct{}{}
	}
}

type Builtin string

const (
	Exit     Builtin = "exit"
	Echo     Builtin = "echo"
	Type     Builtin = "type"
	Pwd      Builtin = "pwd"
	Cd       Builtin = "cd"
	Complete Builtin = "complete"
	Jobs     Builtin = "jobs"
)

type Job struct {
	ID            int
	ProcessID     int
	CommandString string
	Status        string
}

var (
	CurrentJobId int = 1
	JobList          = []Job{}
	JobMutex     sync.Mutex
)

func GetNextJobID() int {
	id := 1

	for {
		found := false

		for _, job := range JobList {
			if job.ID == id {
				found = true
				break
			}
		}

		if !found {
			return id
		}

		id++
	}
}
