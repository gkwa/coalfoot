package coalfoot

import (
	"log/slog"
	"path/filepath"
)

func Main() int {
	slog.Debug("coalfoot", "test", true)

	cf := NewTxtarTemplate()
	slog.Debug("Main", "path", cf.LocalPathRendered)
	slog.Debug("Main", "url", cf.RemoteURL)

	dir, _ := filepath.Abs("tmp")
	err := cf.Extract(dir)
	if err != nil {
		slog.Error("extracting template produced error", "error", err.Error())
	}

	return 0
}
