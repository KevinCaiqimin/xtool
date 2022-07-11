package reversefile

import (
	"io/ioutil"

	"github.com/KevinCaiqimin/log"
)

func reverseFile(filePath string) error {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	for i, b := range buf {
		buf[i] = ^b
	}
	err = ioutil.WriteFile(filePath, buf, 0666)
	if err != nil {
		return err
	}
	log.Info("file: %v reversed successfully", filePath)
	return nil
}

func reverseDir(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error("read dir: %v error: %v", dir, err)
		return
	}
	for _, file := range files {
		filePath := dir + "/" + file.Name()
		if file.IsDir() {
			reverseDir(filePath)
		} else {
			err := reverseFile(filePath)
			if err != nil {
				log.Error("reverse file: %v error: %v", filePath, err)
			}
		}
	}
}

func Start(filePath, dir string) {
	if filePath != "" {
		err := reverseFile(filePath)
		if err != nil {
			log.Error("reverse file: %v error: %v", filePath, err)
		}
	}
	if dir != "" {
		reverseDir(dir)
	}
}
