package sqliteLib

var (
	DBPath string = "test.db"
)

func SetDBPath(newDBPath string) {
	DBPath = newDBPath
}