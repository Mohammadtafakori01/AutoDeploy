package extract

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractZip(zipPath string, destDir string) error {
	staticBase := "/var/www/"
	fullDestDir := filepath.Join(staticBase, destDir)
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		fpath := filepath.Join(fullDestDir, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(fullDestDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			if _, err = io.Copy(outFile, rc); err != nil {
				return err
			}
		}
	}
	return nil
}
