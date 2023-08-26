package pkg

// don't remove these!
const svcSystemName = "FSWatcherSvc"
const svcDisplayName = "FSWatcherSvc"
const svcDescription = "Svc to watch folders on file systems and alert"

// called by Execute()
func InitCustomCode() {
	FSInit()
}
