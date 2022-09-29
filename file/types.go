package file

// LogFile
type LogFile struct {
	Filepath  string
	Typ       string
	Target    string
	Name      string
	Extention string
	KB        int
}

type LogDir struct {
	Dirpath     string
	Typ         string
	Target      string
	CountFiles  int
	FirstFile   string
	LastFile    string
	KB          int
	LastForward string
}
