package file

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/kuoss/lethe/config"
	"github.com/kuoss/lethe/util"
)

// GET

var hourFile = "2???-??-??_??.log"

func ListFiles() []LogFile {
	var logFiles []LogFile
	result, _, err := util.Execute(fmt.Sprintf("ls %s/*/*/%s", config.GetLogRoot(), hourFile))
	if err != nil {
		log.Printf("GetAllFiles: no files err=%e\n", err)
		return logFiles
	}
	rows := strings.Split(strings.Trim(result, "\n"), "\n")
	separator := regexp.MustCompile("[/.]+")
	for _, row := range rows {
		cols := separator.Split(util.SubstrAfter(row, "log/"), -1)
		logFiles = append(logFiles, LogFile{
			Typ:       cols[0],
			Target:    cols[1],
			Name:      cols[2],
			Extention: cols[3],
			Filepath:  row,
		})
	}
	return logFiles
}

func ListFilesWithSize() []LogFile {
	var logFiles []LogFile
	result, _, err := util.Execute(fmt.Sprintf("du %s/*/*/%s", config.GetLogRoot(), hourFile))
	if err != nil {
		log.Printf("GetAllFilesWithSize: no files err=%e\n", err)
		return logFiles
	}
	rows := strings.Split(strings.Trim(result, "\n"), "\n")
	separator := regexp.MustCompile("[\t/.]+")
	for _, row := range rows {
		kbString := util.SubstrBefore(row, "\t")
		filepath := util.SubstrAfter(row, "\t")
		kb, _ := strconv.Atoi(kbString)
		cols := separator.Split(util.SubstrAfter(row, "log/"), -1)
		logFiles = append(logFiles, LogFile{
			Filepath:  filepath,
			Typ:       cols[0],
			Target:    cols[1],
			Name:      cols[2],
			Extention: cols[3],
			KB:        kb, // KiB
		})
	}
	return logFiles
}

func ListDirs() []string {
	out, _, err := util.Execute(fmt.Sprintf("cd %s; ls -d */%s", config.GetLogRoot(), hourFile))
	if err != nil {
		log.Println(fmt.Errorf("%e", err))
		return []string{}
	}
	return strings.Split(strings.TrimSpace(out), "\n")
}

func ListDirsWithSize() []LogDir {
	var logDirs []LogDir
	result1, _, err := util.Execute(fmt.Sprintf("du -s %s/*/*", config.GetLogRoot()))
	if err != nil {
		log.Println(fmt.Errorf("%e", err))
		return logDirs
	}
	result2, _, err := util.Execute(fmt.Sprintf("du --inode %s/*/*", config.GetLogRoot()))
	if err != nil {
		log.Println(fmt.Errorf("%e", err))
		return logDirs
	}
	rows1 := strings.Split(strings.Trim(result1, "\n"), "\n")
	rows2 := strings.Split(strings.Trim(result2, "\n"), "\n")
	separator := regexp.MustCompile("[\t/.]+")
	for i, row1 := range rows1 {
		dirpath := util.SubstrAfter(row1, "\t")
		kb, _ := strconv.Atoi(util.SubstrBefore(row1, "\t"))
		cols1 := separator.Split(util.SubstrAfter(row1, "log/"), -1)

		row2 := rows2[i]
		countFiles, _ := strconv.Atoi(util.SubstrBefore(row2, "\t"))

		re := regexp.MustCompile(`.*/`)
		firstFile, _, _ := util.Execute(fmt.Sprintf("ls %s/%s | head -1", dirpath, hourFile))
		firstFile = re.ReplaceAllString(strings.TrimSpace(firstFile), "")
		lastFile, _, _ := util.Execute(fmt.Sprintf("ls %s/%s | tail -1", dirpath, hourFile))
		lastFile = re.ReplaceAllString(strings.TrimSpace(lastFile), "")
		logDirs = append(logDirs, LogDir{
			Dirpath:    dirpath,
			Typ:        cols1[0],
			Target:     cols1[1],
			CountFiles: countFiles - 1, // exclude directory itself
			KB:         kb,             // KiB
			FirstFile:  firstFile,
			LastFile:   lastFile,
		})
	}
	return logDirs
}

func ListTargets() []LogDir {
	logDirs := ListDirsWithSize()
	for i, logDir := range logDirs {
		if logDir.LastFile == "" {
			continue
		}
		command := fmt.Sprintf(`tail -1000 %s/%s | grep -Po ^[^Z]\+Z | tail -1`, logDir.Dirpath, logDir.LastFile)
		// fmt.Println("command=", command)
		out, _, err := util.Execute(command)
		if err != nil {
			// log.Println(fmt.Errorf("%e", err))
			continue
		}
		logDirs[i].LastForward = strings.TrimSpace(out)
	}
	return logDirs
}
