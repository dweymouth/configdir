package configdir

import "os"

var (
	systemConfig []string
	localConfig  string
	localCache   string
	localState   string
)

func findPaths() {
	systemConfig = []string{os.Getenv("PROGRAMDATA")}
	localConfig = os.Getenv("APPDATA")
	localCache = os.Getenv("LOCALAPPDATA")
	localState = os.Getenv("LOCALAPPDATA")
}
