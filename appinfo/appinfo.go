package appinfo

import (
	"github.com/go-vgo/robotgo"
	"github.com/mitchellh/go-ps"
)

type AppInfo struct {
	Pid  int32
	Name string
}

func getForeground() (AppInfo, bool) {
	pid := robotgo.GetPID()
	process, err := ps.FindProcess(int(pid))
	if err != nil {
		panic(err)
	}
	appInfo := AppInfo{}
	isProcessNonNil := process != nil
	if isProcessNonNil {
		appInfo = AppInfo{
			Pid:  pid,
			Name: process.Executable(),
		}
	}
	return appInfo, isProcessNonNil
}

func GetChangedForegroundInfo(ch chan AppInfo) {
	appInfoDefault, _ := getForeground()
	for {
		appInfoForeground, ok := getForeground()
		if !ok {
			continue
		}
		if appInfoDefault.Pid != appInfoForeground.Pid {
			ch <- appInfoForeground
			break
		}
	}
}
