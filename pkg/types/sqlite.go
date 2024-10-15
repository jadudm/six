package types

import "fmt"

type SqliteFile struct {
	Size int64
	Name string
}

func (s *SqliteFile) Weight() int64 {
	return int64(s.Size)
}

func (s *SqliteFile) Value() int64 {
	return 1
}

func (s *SqliteFile) String() string {
	return fmt.Sprintf("%s[%d]", s.Name, s.Size)
}

func (s *SqliteFile) Id() string {
	return s.Name
}
