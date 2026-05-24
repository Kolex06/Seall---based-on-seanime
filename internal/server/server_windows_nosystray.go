//go:build windows && nosystray

package server

import (
	"embed"
	"fmt"
	"net"
	"time"

	"github.com/cli/browser"
)

func StartServer(webFS embed.FS, embeddedLogo []byte) {

	app, flags, selfupdater := startApp(embeddedLogo)

	if !flags.Update && !flags.IsDesktopSidecar {
		go func() {
			addr := fmt.Sprintf("127.0.0.1:%d", app.Config.Server.Port)
			for i := 0; i < 80; i++ {
				conn, err := net.DialTimeout("tcp", addr, 250*time.Millisecond)
				if err == nil {
					_ = conn.Close()
					_ = browser.OpenURL(app.Config.GetServerURI("127.0.0.1"))
					return
				}
				time.Sleep(250 * time.Millisecond)
			}
		}()
	}

	startAppLoop(&webFS, app, flags, selfupdater)
}
