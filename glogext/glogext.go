package glogext

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/golang/glog"

	"github.com/jamessynge/go_io/fileio"
)

// glog doesn't look at the value of log_dir except when it is first logging
// a message of a higher severity than previously logged, or when it is
// rotating log files.  So, call this immediately after parsing the command
// line (i.e. before logging).
// Note that we don't actually set the default value of the flag because that
// won't be examined by glog when the flag is unset.
func SetDefaultLogDir(default_log_dir string) {
	/*if logtostderr, ok := LogToStderrValue(); ok && logtostderr {
		// No need to set it because logging to stderr.
		return
	}*/
	//	fmt.Printf("SetDefaultLogDir(%q)\n", default_log_dir)
	log_dir_flag := flag.Lookup("log_dir")
	if log_dir_flag == nil {
		// Flag doesn't exist!
		fmt.Println("Flag doesn't exist!")
		return
	}
	if len(log_dir_flag.Value.String()) > 0 {
		fmt.Println("Already set.")
		// Already set.
		return
	}
	if !fileio.IsDirectory(default_log_dir) {
		if err := os.MkdirAll(default_log_dir, 0750); err != nil {
			glog.Fatalf(`Unable to create log directory!
 Path: %q
Error: %s`, default_log_dir, err)
		}
	}
	fmt.Printf("Setting --log_dir to %v\n", default_log_dir)

	log_dir_flag.Value.Set(default_log_dir)
	glog.V(1).Infof("Set --log_dir to %q", default_log_dir)
}
