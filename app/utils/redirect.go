package utils

import "strings"

type Redirect struct {
	Cmd    string
	File   string
	FD     int // 1=stdout, 2=stderr
	Append bool
	Valid  bool
}

func ParseRedirect(cmd string) Redirect {
	ops := []struct {
		token  string
		fd     int
		append bool
	}{
		{"2>>", 2, true},
		{"1>>", 1, true},
		{">>", 1, true},
		{"2>", 2, false},
		{"1>", 1, false},
		{">", 1, false},
	}

	for _, op := range ops {
		parts := strings.SplitN(cmd, op.token, 2)

		if len(parts) != 2 {
			continue
		}

		return Redirect{
			Cmd:    strings.TrimSpace(parts[0]),
			File:   strings.TrimSpace(parts[1]),
			FD:     op.fd,
			Append: op.append,
			Valid:  true,
		}
	}

	return Redirect{
		Cmd:   cmd,
		Valid: false,
	}
}
