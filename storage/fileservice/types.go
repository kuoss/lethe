package fileservice

// LogFile
type LogFile struct {
	Fullpath  string
	Subpath   string
	LogType   string
	Target    string
	Name      string
	Extension string
	Size      int64
}

type LogDir struct {
	Fullpath    string
	Subpath     string
	LogType     string
	Target      string
	FileCount   int
	FirstFile   string
	LastFile    string
	Size        int64
	LastForward string
}
