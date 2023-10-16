package coalfoot

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	mymazda "github.com/taylormonacelli/forestfish/mymazda"
)

type TxtarTemplate struct {
	RemoteURL           string
	LocalPathUnrendered string
	LocalPathRendered   string
}

var baseDir = "/tmp/coalfoot"

func NewTxtarTemplate() *TxtarTemplate {
	fname := "1.txtar"
	fnameRendered := "1-rendered.txt"
	url := fmt.Sprintf("https://raw.githubusercontent.com/taylormonacelli/navylie/master/templates/%s", fname)

	return &TxtarTemplate{
		RemoteURL:           url,
		LocalPathUnrendered: filepath.Join(baseDir, fname),
		LocalPathRendered:   filepath.Join(baseDir, fnameRendered),
	}
}

func (tpl TxtarTemplate) RemotePath() string {
	return tpl.RemoteURL
}

func (tpl TxtarTemplate) GetRenderedTemplateDir() string {
	return filepath.Dir(tpl.LocalPathRendered)
}

func (tpl TxtarTemplate) GetUnRenderedTemplateDir() string {
	return filepath.Dir(tpl.LocalPathUnrendered)
}

func (tpl TxtarTemplate) FetchIfNotFound() {
	if mymazda.FileExists(tpl.LocalPathUnrendered) && durationSinceFileCreated(tpl.LocalPathUnrendered) < 24*time.Hour {
		return
	}

	slog.Debug("fetching", "url", tpl.RemoteURL, "path", tpl.LocalPathUnrendered)
	tpl.FetchRemoteToLocal()
}

func durationSinceFileCreated(filePath string) time.Duration {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		slog.Error("error getting file info", "path", filePath, "error", err.Error())
		return 0
	}

	durationSince := time.Since(fileInfo.ModTime())

	slog.Debug("duration since file was modified", "path", filePath, "duration", durationSince)

	return durationSince
}

func (tpl TxtarTemplate) FetchRemoteToLocal() error {
	dir, err := filepath.Abs(filepath.Dir(tpl.LocalPathUnrendered))
	if err != nil {
		slog.Error("coalfoot", "filepath.abs", tpl.LocalPathUnrendered, "error", err.Error())
		return err
	}

	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		slog.Error("coalfoot mkdir", "mkdir", dir, "error", err.Error())
		return err
	}

	err = fetchTemplateToPath(tpl.RemoteURL, tpl.LocalPathUnrendered)
	if err != nil {
		slog.Error("fetch template failed", "url", tpl.RemoteURL, "target", tpl.LocalPathUnrendered, "error", err.Error())
	}

	return nil
}

func fetchTemplateToPath(url, localPath string) error {
	directory := filepath.Dir(localPath)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return err
	}

	fileName := filepath.Join(directory, filepath.Base(url))

	absPath, _ := filepath.Abs(fileName)

	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		slog.Debug("file already exists, not refetching", "path", absPath)
		return nil

	} else {
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		file, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			return err
		}

		slog.Debug("file saved", "path", absPath)
	}

	return nil
}
