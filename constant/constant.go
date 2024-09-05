package constant

const AppDir string = "push-me-client"
const AppID string = "PushMeClient"
const AppName string = "PushMe Client"
const AppVersion string = "v2.0.0"
const UrlVersion string = "https://push.i-i.me/api/version"

const (
	LnkName      string = AppID + ".lnk"
	LinkSuffix   string = "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup/" + LnkName
	StartMenuLnk string = "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/" + LnkName
)
