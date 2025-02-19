package remotree

import (
	"archive/zip"
	"bufio"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/htfs"
	"github.com/robocorp/rcc/operations"
	"github.com/robocorp/rcc/set"
)

func isSelfRequest(request *http.Request) bool {
	identity, ok := request.Header[operations.X_RCC_RANDOM_IDENTITY]
	return ok && len(identity) > 0 && identity[0] == common.RandomIdentifier()
}

func makeDeltaHandler(queries Partqueries) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		catalog := filepath.Base(request.URL.Path)
		defer common.Stopwatch("Delta of catalog %q took", catalog).Debug()
		if request.Method != http.MethodPost {
			response.WriteHeader(http.StatusMethodNotAllowed)
			common.Trace("Delta: rejecting request %q for catalog %q.", request.Method, catalog)
			return
		}
		if isSelfRequest(request) {
			response.WriteHeader(http.StatusConflict)
			common.Trace("Delta: rejecting /SELF/ request for catalog %q.", catalog)
			return
		}
		reply := make(chan string)
		queries <- &Partquery{
			Catalog: catalog,
			Reply:   reply,
		}
		known, ok := <-reply
		common.Debug("query handler: %q -> %v", catalog, ok)
		if !ok {
			response.WriteHeader(http.StatusNotFound)
			response.Write([]byte("404 not found, sorry"))
			return
		}

		membership := set.Membership(strings.Split(known, "\n"))

		approved := make([]string, 0, 1000)
		todo := bufio.NewReader(request.Body)
	todoloop:
		for {
			line, err := todo.ReadString('\n')
			stopping := err == io.EOF
			candidate := filepath.Base(strings.TrimSpace(line))
			if len(candidate) > 10 {
				if membership[candidate] {
					approved = append(approved, candidate)
				} else {
					common.Trace("DELTA: ignoring extra %q entry, not part of set!", candidate)
					if !stopping {
						continue todoloop
					}
				}
			}
			if stopping {
				break todoloop
			}
			if err != nil {
				common.Trace("DELTA: error %v with line %q", err, line)
				break todoloop
			}
		}

		headers := response.Header()
		headers.Add("Content-Type", "application/zip")
		response.WriteHeader(http.StatusOK)

		sink := zip.NewWriter(response)
		defer sink.Close()

		for _, member := range approved {
			relative := htfs.RelativeDefaultLocation(member)
			fullpath := htfs.ExactDefaultLocation(member)
			err := operations.ZipAppend(sink, fullpath, relative)
			if err != nil {
				common.Debug("DELTA: error %v with %v -> %v", err, fullpath, relative)
				return
			}
		}

		fullpath := filepath.Join(common.HololibCatalogLocation(), catalog)
		relative, err := filepath.Rel(common.HololibLocation(), fullpath)
		if err != nil {
			common.Debug("DELTA: error %v", err)
			return
		}
		err = operations.ZipAppend(sink, fullpath, relative)
		if err != nil {
			common.Debug("DELTA: error %v", err)
			return
		}
	}
}
