package file

// LogFile
type LogFile struct {
	FullPath  string
	SubPath   string
	Typ       string
	Target    string
	Name      string
	Extention string
	Size      int64
}

type LogDir struct {
	FullPath    string
	SubPath     string
	Typ         string
	Target      string
	FileCount   int
	FirstFile   string
	LastFile    string
	Size        int64
	LastForward string
}
