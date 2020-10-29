package appinfo

import (
	"github.com/go-vgo/robotgo"
	"github.com/mitchellh/go-ps"
)

type AppInfo struct {
	Pid  int32
	Name string
}

func getForeground() AppInfo {
	pid := robotgo.GetPID()
	pidInfo, err := ps.FindProcess(int(pid))
	if err != nil {
		panic(err)
	}
	return AppInfo{
		Pid:  pid,
		Name: pidInfo.Executable(),
	}
}

func GetChangedForegroundInfo(ch chan AppInfo) {
	appInfoDefault := getForeground()
	for {
		appInfoForeground := getForeground()
		if appInfoDefault.Pid != appInfoForeground.Pid {
			ch <- appInfoForeground
			break
		}
	}
}
