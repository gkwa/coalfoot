package coalfoot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"log/slog"

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
	if !mymazda.FileExists(tpl.LocalPathUnrendered) {
		slog.Debug("fetching", "url", tpl.RemoteURL, "path", tpl.LocalPathUnrendered)
		tpl.FetchRemoteToLocal()
	}
}

func (tpl TxtarTemplate) FetchRemoteToLocal() error {
	dir, err := filepath.Abs(filepath.Dir(tpl.LocalPathUnrendered))
	if err != nil {
		slog.Error("coalfoot", "filepath.abs", tpl.LocalPathUnrendered, "error", err.Error())
		return err
	}

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		slog.Error("coalfoot mkdir", "mkdir", dir, "error", err.Error())
		return err
	}

	fetchTemplateToPath(tpl.RemoteURL, tpl.LocalPathUnrendered)

	return nil
}

func fetchTemplateToPath(url, localPath string) {
	directory := filepath.Dir(localPath)

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		fmt.Printf("Error creating directories: %v\n", err)
		return
	}

	fileName := filepath.Join(directory, filepath.Base(url))

	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		x, _ := filepath.Abs(fileName)
		slog.Debug("file already exists, not refetching", "path", x)
	} else {
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error fetching the file: %v\n", err)
			return
		}
		defer response.Body.Close()

		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("Error creating the file: %v\n", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Printf("Error saving the file: %v\n", err)
			return
		}

		x, _ := filepath.Abs(fileName)
		slog.Debug("file saved", "path", x)
	}
}
