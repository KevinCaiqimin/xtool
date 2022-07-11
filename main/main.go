package main

import (
	"flag"
	"fmt"
	"xtool/cfgexport"

	// "os"
	"helloserver/log"
	"xtool/apps/reversefile"
	"xtool/filedecrypt"
	"xtool/fileencrypt"
)

var app = flag.String("app", "", "app to be run")
var js = flag.Bool("js", false, "export to json")
var srcdir = flag.String("srcdir", "./", "where config excel files will be found")
var jsdir = flag.String("jsdir", "./", "exported js file save directory")
var lua = flag.Bool("lua", false, "export to lua")
var luadir = flag.String("luadir", "./", "exported lua file save directory")
var py = flag.Bool("py", false, "export to python")
var pydir = flag.String("pydir", "./", "exported python file save directory")
var cs = flag.Bool("cs", false, "export to csharp")
var csdir = flag.String("csdir", "./", "exported csharp file save directory")
var ignore = flag.String("i", "", "ignore files")
var usesheet = flag.Bool("usesheet", false, "when export excel, sheet is used only")
var f = flag.String("f", "", "file specify")
var d = flag.String("d", "", "directory specify")
var inf = flag.String("inf", "", "input file")
var outf = flag.String("outf", "", "output file")
var ind = flag.String("ind", "", "input directory")
var outd = flag.String("outd", "", "output directory")
var aesKey = flag.String("aesKey", "", "encrypt key")
var aesIV = flag.String("aesIV", "", "encrypt iv")
var encFileExt = flag.Bool("encFileExt", false, "encrypt file ext")
var encFileName = flag.Bool("encFileName", false, "encrypt file name")

func main() {
	flag.Parse()

	log.InitLog("console", "none", log.LV_DEBUG)

	switch *app {
	case "enc": //文件加密
		fileencrypt.AES_KEY = *aesKey
		fileencrypt.AES_IV = *aesIV
		fileencrypt.EncFileExt = *encFileExt
		fileencrypt.EncFileName = *encFileName
		fileencrypt.Start(*inf, *outf, *ind, *outd)
		break
	case "dec": //文件解密
		filedecrypt.AES_KEY = *aesKey
		filedecrypt.AES_IV = *aesIV
		filedecrypt.EncFileExt = *encFileExt
		filedecrypt.EncFileName = *encFileName
		filedecrypt.Start(*inf, *outf, *ind, *outd)
		break
	case "reversefile":
		reversefile.Start(*f, *d)
		break
	case "configexport":
		cfgexport.ExportCfg(&cfgexport.ExportConfig{
			Js:       *js,
			Jsdir:    *jsdir,
			Lua:      *lua,
			Luadir:   *luadir,
			Py:       *py,
			Pydir:    *pydir,
			Cs:       *cs,
			Csdir:    *csdir,
			SrcDir:   *srcdir,
			Ignore:   *ignore,
			UseSheet: *usesheet,
		})
		break
	default:
		fmt.Printf("error app: %s\n", *app)
	}
}
