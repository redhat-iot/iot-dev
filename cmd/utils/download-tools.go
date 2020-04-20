package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	ocUrl      = "https://mirror.openshift.com/pub/openshift-v4/clients/ocp/latest/openshift-client-linux-4.3.5.tar.gz"
	enmasseUrl = "https://github.com/EnMasseProject/enmasse/releases/download/0.30.2/enmasse-0.30.2.tgz"
)

func downloadPkg(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func untar(dst string, r io.Reader) (filename string, err error) {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return filename, err

	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	i := false

	for {

		header, err := tr.Next()
		if !i {
			filename = header.Name
			i = true
		}
		///log.Println(header.Name)
		switch {

		// if no more files are found return
		case err == io.EOF:
			return filename, nil

		// return any other error
		case err != nil:
			return filename, err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return filename, err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return filename, err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return filename, err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
	//return filename, err
}

func DownloadAndUncompress(name string, url string) string {
	err := downloadPkg(name, url)
	if err != nil {
		log.Fatal("Error downloading Package: ", err)
	}

	content, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filename, err := untar(path, content)
	if err != nil {
		log.Fatal("Error uncompressing package: ", err)
	}

	os.Remove(name)

	return filename

}
