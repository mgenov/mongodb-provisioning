package tarutil

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Untar reads the gzip-compressed tar file from r and writes it into dir.
func Untar(r io.Reader, dir string) error {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	tr := tar.NewReader(zr)
	for {
		f, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("tar reading error: %v", err)
			return err
		}

		rel := filepath.FromSlash(f.Name)
		abs := filepath.Join(dir, rel)

		fi := f.FileInfo()
		mode := fi.Mode()
		switch {
		case mode.IsRegular():
			// Make the directory. This is redundant because it should
			// already be made by a directory entry in the tar
			// beforehand. Thus, don't check for errors; the next
			// write will fail with the same error.
			os.MkdirAll(filepath.Dir(abs), 0755)
			wf, err := os.OpenFile(abs, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return err
			}
			n, err := io.Copy(wf, tr)
			if closeErr := wf.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
			if err != nil {
				return fmt.Errorf("error writing to %s: %v", abs, err)
			}
			if n != f.Size {
				return fmt.Errorf("only wrote %d bytes to %s; expected %d", n, abs, f.Size)
			}
			log.Printf("wrote %s", abs)
			if !f.ModTime.IsZero() {
				if err := os.Chtimes(abs, f.ModTime, f.ModTime); err != nil {
					// benign error. Gerrit doesn't even set the
					// modtime in these, and we don't end up relying
					// on it anywhere (the gomote push command relies
					// on digests only), so this is a little pointless
					// for now.
					log.Printf("error changing modtime: %v", err)
				}
			}
		case mode.IsDir():
			if err := os.MkdirAll(abs, 0755); err != nil {
				return err
			}
		default:
			return fmt.Errorf("tar file entry %s contained unsupported file type %v", f.Name, mode)
		}
	}
	return nil
}
