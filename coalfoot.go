package coalfoot

import (
	"log/slog"
)

func Main() {
	slog.Debug("coalfoot", "test", true)

	cf := NewTxtarTemplate()
	slog.Debug("Main", "path", cf.LocalPathRendered)
	slog.Debug("Main", "url", cf.RemoteURL)
}
