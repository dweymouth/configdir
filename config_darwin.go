package configdir

import "os"

var (
	systemConfig []string
	localConfig  string
	localCache   string
	localState   string
)

func findPaths() {
	systemConfig = []string{"/Library/Application Support"}
	localConfig = os.Getenv("HOME") + "/Library/Application Support"
	localCache = os.Getenv("HOME") + "/Library/Caches"
	localState = os.Getenv("HOME") + "/Library/Application Support"
}
