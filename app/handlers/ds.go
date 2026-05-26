package handlers

var Builtins = map[Builtin]func(string){
	Echo: HandleEcho,
	Type: HandleType,
	Exit: func(s string) {},
	Pwd:  HandlePwd,
	Cd:   HandleCd,
}

var BuiltinNames = map[Builtin]struct{}{}

func IsBuiltin(cmd string) bool {
	_, ok := BuiltinNames[Builtin(cmd)]
	return ok
}

type Builtin string

const (
	Exit Builtin = "exit"
	Echo Builtin = "echo"
	Type Builtin = "type"
	Pwd  Builtin = "pwd"
	Cd   Builtin = "cd"
)
