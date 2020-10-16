package sqliteLib

var (
	DBPath = "test.db"
	TokenLength = 64
)

func SetDBPath(newDBPath string) {
	DBPath = newDBPath
}