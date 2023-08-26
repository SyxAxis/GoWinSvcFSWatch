package pkg

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	// #TODO: look this up
	// https://github.com/gohugoio/hugo/tree/master/watcher
)

var wg sync.WaitGroup

type configAPPJSONFile struct {
	WatcherName string `json:"watchername,omitempty"`
	SourcePath  string `json:"sourcepath,omitempty"`
	FileMask    string `json:"filemask,omitempty"`
	ExecuteCmd  string `json:"executeCommand,omitempty"`
}

var activeAPPConfig = []configAPPJSONFile{}

// type outputJSON struct {
// 	WatcherName string `json:"watchername"`
// 	SourcePath  string `json:"sourcepath"`
// 	Operation   string `json:"operation"`
// }

// =====================================================================================
// MAIN
// =====================================================================================

func FSInit() {

	tmpAppCfg := configAPPJSONFile{
		WatcherName: "TestWatcher",
		SourcePath:  "c:\\temp\\crap",
		FileMask:    "*.txt",
		ExecuteCmd:  "",
	}

	activeAPPConfig = append(activeAPPConfig, tmpAppCfg)

	for i, _ := range activeAPPConfig {
		wg.Add(1)
		go runWatcher(activeAPPConfig[i])
	}

	wg.Wait()
}

// =====================================================================================

func runWatcher(activeConfig configAPPJSONFile) {

	defer wg.Done()

	var (
		// Wait 100ms for new events; each new event resets the timer.
		waitFor = 100 * time.Millisecond

		// Keep track of the timers, as path → timer.
		mu     sync.Mutex
		timers = make(map[string]*time.Timer)

		// Callback we run.
		printEvent = func(e fsnotify.Event) {
			// TODO
			windowsEventlog.Info(1, fmt.Sprintf(" - FSWatcher Event : [ %s ][ %s ]\n", e.Op.String(), e.Name))

			// Don't need to remove the timer if you don't have a lot of files.
			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	windowsEventlog.Info(1, fmt.Sprintf(" - Init FSWatcher : [ %s ][ %s ]\n", activeConfig.WatcherName, activeConfig.SourcePath))

	// fmt.Printf("   - Source Path  : [ %s ]\n", activeConfig.SourcePath)
	// fmt.Printf("   - File Mask    : [ %s ]\n", activeConfig.FileMask)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		windowsEventlog.Error(1, err.Error())
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// We just want to watch for file creation, so ignore everything
				// outside of Create and Write.
				// if !event.Has(fsnotify.Create) && !event.Has(fsnotify.Write) {
				if !event.Has(fsnotify.Create) {
					continue
				}

				// Get timer.
				mu.Lock()
				t, ok := timers[event.Name]
				mu.Unlock()

				// No timer yet, so create one.
				if !ok {
					t = time.AfterFunc(math.MaxInt64, func() { printEvent(event) })
					t.Stop()

					mu.Lock()
					timers[event.Name] = t
					mu.Unlock()
				}

				// Reset the timer for this path, so it will start from 100ms again.
				t.Reset(waitFor)

				// log.Printf("EVENT [%v] from [%v]\n", event.Op, event.Name)
				// if event.Op&fsnotify.Write == fsnotify.Write {
				// 	log.Println("modified file:", event.Name)
				// }
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				windowsEventlog.Error(1, err.Error())
			}
		}
	}()

	// NOTE - this is just the path NOT the filemask as well!

	err = watcher.Add(activeConfig.SourcePath)
	if err != nil {
		windowsEventlog.Error(1, err.Error())
	}
	<-make(chan struct{})

	// w := watcher.New()
	// w.SetWatcherName(activeConfig.WatcherName)

	// // Only files that match the regular expression during file listing
	// // will be watched.
	// r := regexp.MustCompile(activeConfig.FileMask)
	// w.AddFilterHook(watcher.RegexFilterHook(r, false))

	// // Uncomment to use SetMaxEvents set to 1 to allow at most 1 event to be received
	// // on the Event channel per watching cycle.
	// //
	// // If SetMaxEvents is not set, the default is to send all events.
	// // w.SetMaxEvents(1)

	// // Uncomment to only notify rename and move events.
	// // watcher.Create, watcher.Write, watcher.Remove, watcher.Rename, watcher.Chmod, watcher.Move
	// w.FilterOps(watcher.Create)

	// go func() {
	// 	jsonReq := []byte{}

	// 	for {

	// 		select {
	// 		case event := <-w.Event:

	// 			// fmt.Printf("==================================================\n")
	// 			// fmt.Printf("Watcher      : %+v\n", w.GetWatcherName())
	// 			// fmt.Printf("Operation    : %+v\n", event.Op.String())
	// 			// fmt.Printf("Fullpath     : %+v\n", event.Path)

	// 			// fmt.Printf("Old Path(mv) : %+v\n", event.OldPath)
	// 			// fmt.Printf("IsDir        : %+v\n", event.IsDir())
	// 			// fmt.Printf("ModTime      : %+v\n", event.ModTime())
	// 			// fmt.Printf("Mode         : %+v\n", event.Mode())
	// 			// fmt.Printf("Name         : %+v\n", event.Name())
	// 			// fmt.Printf("Size (bytes) : %+v\n", event.Size())
	// 			// fmt.Printf("Lcl Sys Attrs: %+v\n", event.Sys())

	// 			if event.Path != "-" {
	// 				if !event.IsDir() {
	// 					opj := outputJSON{w.GetWatcherName(),
	// 						event.Op.String(),
	// 						event.Path}

	// 					jsonReq, _ = json.Marshal(opj)
	// 					fmt.Println(string(jsonReq))

	// 					if len(activeConfig.ExecuteCmd) > 0 {
	// 						go executeSystemCommand(activeConfig, event.Path)
	// 					}

	// 				}
	// 			}

	// 		case err := <-w.Error:
	// 			log.Fatalln(err)
	// 		case <-w.Closed:
	// 			return
	// 		}
	// 	}
	// }()

	// if err := w.Add(activeConfig.SourcePath); err != nil {
	// 	log.Fatalln(err)
	// }

	// // Watch test_folder recursively for changes.
	// // if err := w.AddRecursive("\\\\gjpc\\rimes"); err != nil {
	// // 	log.Fatalln(err)
	// // }

	// // Print a list of all of the files and folders currently
	// // being watched and their paths.
	// // for path, f := range w.WatchedFiles() {
	// // 	fmt.Printf("%s: %s\n", path, f.Name())
	// // }

	// // Trigger 2 events after watcher started.
	// // go func() {
	// // 	w.Wait()
	// // 	w.TriggerEvent(watcher.Create, nil)
	// // 	w.TriggerEvent(watcher.Remove, nil)
	// // }()

	// // Start the watching process - it'll check for changes every 100ms.
	// if err := w.Start(time.Millisecond * 100); err != nil {
	// 	log.Fatalln(err)
	// }

	// wg.Done()

}

// func executeSystemCommand(activeConfig configAPPJSONFile, fileFullPath string) {

// 	//time.Sleep(time.Second * 2)

// 	cmdToExec := strings.Fields(strings.Replace(activeConfig.ExecuteCmd, "%SOURCEFILE%", fileFullPath, 1))
// 	// fmt.Printf("Executing : %v\n", cmdToExec)

// 	cmd := exec.Command(cmdToExec[0], cmdToExec[1:]...)
// 	_, err := cmd.Output()
// 	if err != nil {
// 		fmt.Printf("Cmd Exec Error: %v\n", err.Error())
// 		return
// 	}

// 	// Print the output
// 	// fmt.Println(string(stdout))
// }
