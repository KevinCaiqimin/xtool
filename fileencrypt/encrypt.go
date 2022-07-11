package fileencrypt

import (
	"helloserver/log"
	"io/ioutil"
	"path"
	"strings"

	"caiqimin.tech/basic/encrypt"
	"caiqimin.tech/basic/utils"
	"github.com/datadog/zstd"
)

var AES_KEY string
var AES_IV string
var EncFileName bool
var EncFileExt bool

func encFile(srcPath, dstPath string) error {
	buf, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}
	dst, err := zstd.CompressLevel(nil, buf, 22)
	if err != nil {
		return err
	}
	dst, err = encrypt.ToAesCbcBytes(dst, AES_KEY, AES_IV, encrypt.PKCS7)
	if err != nil {
		return err
	}
	utils.EnsurePath(dstPath)
	err = ioutil.WriteFile(dstPath, dst, 0666)
	if err != nil {
		return err
	}
	log.Info("file: %v encrypt to %v successfully", srcPath, dstPath)
	return nil
}

func encDir(srcDir, dstDir string) error {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Error("read dir: %v error: %v", srcDir, err)
		return err
	}
	for _, file := range files {
		srcPath := srcDir + "/" + file.Name()
		dstPath := dstDir + "/" + file.Name()
		if file.IsDir() {
			encDir(srcPath, dstPath)
		} else {
			if EncFileName {
				baseName, err := encrypt.ToAesCbcBase64String(file.Name(),
					AES_KEY, AES_IV, encrypt.PKCS7)
				if err != nil {
					log.Error("encrypt file name %v error: %v", file.Name(), err)
					return err
				}
				dstPath = dstDir + "/" + baseName
			} else if EncFileExt {
				ext := path.Ext(file.Name())
				shortName := strings.TrimRight(file.Name(), ext)
				if len(ext) > 0 {
					ext = ext[1:]
					encExt, err := encrypt.ToAesCbcBase64String(ext,
						AES_KEY, AES_IV, encrypt.PKCS7)
					if err != nil {
						log.Error("encrypt file name %v error: %v", file.Name(), err)
						return err
					}
					dstPath = dstDir + "/" + shortName + "." + encExt
				} else {
					dstPath = dstDir + "/" + shortName
				}
			}
			err := encFile(srcPath, dstPath)
			if err != nil {
				log.Error("enc file: %v error: %v", srcPath, err)
				return err
			}
		}
	}
	return nil
}
