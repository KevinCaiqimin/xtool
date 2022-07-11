package fileencrypt

import "github.com/KevinCaiqimin/log"

func Start(srcFile, dstFile, srcDir, dstDir string) {
	if srcFile != "" {
		err := encFile(srcFile, dstFile)
		if err != nil {
			log.Error("encrypt file: %v error: %v", srcFile, err)
		} else {
			log.Info("encrypt file: %v finished", srcFile)
		}
	}
	if srcDir != "" {
		err := encDir(srcDir, dstDir)
		if err != nil {
			log.Error("encrypt dir: %v error: %v", srcDir, err)
		} else {
			log.Info("encrypt dir: %v finished", srcDir)
		}
	}
}
