package bins

import (
	"net/http"
	"os"
	"path/filepath"
)

func BinHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	Path := filepath.Clean(r.URL.Path)
	if Path == "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	FilePath := filepath.Join("bins", "compiles", Path)

	File, err := os.Open(FilePath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer File.Close()

	Info, err := File.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, Info.Name(), Info.ModTime(), File)
}
