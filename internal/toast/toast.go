package toast

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

func Toast(content string) {
	app := application.Get()

	app.Event.Emit("toast", content)
}
