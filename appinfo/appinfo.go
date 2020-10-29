package appinfo

import (
	"os/exec"
	"regexp"
)

type AppInfo struct {
	Info string
	Name string
}

func getForegroundInfo() string {
	infoForeground, err := exec.Command(`lsappinfo`, `front`).Output()
	if err != nil {
		panic(err)
	}
	return string(infoForeground)
}

func getForegroundName() string {
	infoForeground, err := exec.Command(`lsappinfo`, `info`, `-only`, `name`, getForegroundInfo()).Output()
	if err != nil {
		panic(err)
	}
	nameForeground := regexp.MustCompile(`"LSDisplayName"="(.*?)"`).FindSubmatch(infoForeground)[1]
	return string(nameForeground)
}

func getForeground() AppInfo {
	return AppInfo{
		Info: getForegroundInfo(),
		Name: getForegroundName(),
	}
}

func GetChangedForegroundInfo(ch chan AppInfo) {
	appInfoDefault := getForeground()
	for {
		appInfoForeground := getForeground()
		if appInfoDefault.Info != appInfoForeground.Info {
			ch <- appInfoForeground
			break
		}
	}
}
