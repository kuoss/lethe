package file

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kuoss/lethe/config"
)

func Cleansing() {
	cleansingLogFiles("host")
	cleansingLogFiles("kube")
}

func cleansingLogFiles(prefix string) {
	files, err := filepath.Glob(fmt.Sprintf("%s/%s.*", config.GetLogRoot(), prefix))
	if err != nil {
		fmt.Printf("error on cleansingLogFiles(%s): %s", prefix, err)
		return
	}
	if len(files) < 1 {
		return
	}
	log.Printf("Warning: need cleansing log files(%s).\n", prefix)
	for _, file := range files {
		log.Printf("deleting file... %s", file)
		e := os.Remove(file)
		if e != nil {
			log.Printf("error on deleting file... %s", file)
		}
	}
}
