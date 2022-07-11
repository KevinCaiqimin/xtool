package filedecrypt

import "helloserver/log"

func Start(srcFile, dstFile, srcDir, dstDir string) {
	if srcFile != "" {
		err := decFile(srcFile, dstFile)
		if err != nil {
			log.Error("decrypt file: %v error: %v", srcFile, err)
		} else {
			log.Info("decrypt file: %v finished", srcFile)
		}
	}
	if srcDir != "" {
		err := decDir(srcDir, dstDir)
		if err != nil {
			log.Error("decrypt dir: %v error: %v", srcDir, err)
		} else {
			log.Info("decrypt dir: %v finished", srcDir)
		}
	}
}
