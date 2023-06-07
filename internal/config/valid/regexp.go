package valid

import "regexp"

var (
	Key      = regexp.MustCompile(`^[A-Za-z0-9-_]+$`)
	Name     = regexp.MustCompile(`^[A-Za-z0-9-_.]+$`)
	Path     = regexp.MustCompile(`^.+(?:\\)|.+(?:/)`)
	Password = regexp.MustCompile(`^[A-Za-z0-9!#$%&*+-.:;<=>?@^_{|}~]+$`)
)
