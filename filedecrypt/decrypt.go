package filedecrypt

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/KevinCaiqimin/log"

	"github.com/DataDog/zstd"
	"github.com/KevinCaiqimin/go-basic/encrypt"
	"github.com/KevinCaiqimin/go-basic/utils"
)

var AES_KEY string
var AES_IV string
var EncFileName bool
var EncFileExt bool

func decFile(srcPath, dstPath string) error {
	buf, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return err
	}
	buf, err = encrypt.FromAesCbc(buf, AES_KEY, AES_IV, encrypt.PKCS7)
	if err != nil {
		return err
	}
	dst, err := zstd.Decompress(nil, buf)
	if err != nil {
		return err
	}
	utils.EnsurePath(dstPath)
	err = ioutil.WriteFile(dstPath, dst, 0666)
	if err != nil {
		return err
	}
	log.Info("file: %v decrypt to %v successfully", srcPath, dstPath)
	return nil
}

func decDir(srcDir, dstDir string) error {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Error("read dir: %v error: %v", srcDir, err)
		return err
	}
	for _, file := range files {
		srcPath := srcDir + "/" + file.Name()
		dstPath := dstDir + "/" + file.Name()
		if file.IsDir() {
			decDir(srcPath, dstPath)
		} else {
			if EncFileName {
				baseName, err := encrypt.FromAesCbcBase64String(file.Name(),
					AES_KEY, AES_IV, encrypt.PKCS7)
				if err != nil {
					log.Error("decrypt file name %v error: %v", file.Name(), err)
					return err
				}
				dstPath = dstDir + "/" + baseName
			} else if EncFileExt {
				ext := path.Ext(file.Name())
				shortName := strings.TrimRight(file.Name(), ext)
				if len(ext) > 0 {
					ext = ext[1:]
					encExt, err := encrypt.FromAesCbcBase64String(ext,
						AES_KEY, AES_IV, encrypt.PKCS7)
					if err != nil {
						log.Error("decrypt file name %v error: %v", file.Name(), err)
						return err
					}
					dstPath = dstDir + "/" + shortName + "." + encExt
				} else {
					dstPath = dstDir + "/" + shortName
				}
			}
			err := decFile(srcPath, dstPath)
			if err != nil {
				log.Error("decrypt file: %v error: %v", srcPath, err)
				return err
			}
		}
	}
	return nil
}
