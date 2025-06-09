package sql

import (
	"embed"
	"io"
	"io/fs"
	"log/slog"
	"strings"
)

//go:embed *.sql
var embedQueries embed.FS

func GetSQLFromFile(fileName string) (string, error) {
	fsys, err := fs.Sub(embedQueries, ".")
	if err != nil {
		slog.Error("Error reading embedded SQL-file dir", slog.String("error", err.Error()))
		return "Error reading embedded SQL-file", err
	}

	file, err := fsys.Open(fileName)
	if err != nil {
		slog.Error("Error opening SQL-file", slog.String("error", err.Error()))
		return "Error opening SQL-file", err
	}

	buf := new(strings.Builder)
	_, ioErr := io.Copy(buf, file)

	if ioErr != nil {
		slog.Error("Error reading SQL-file", slog.String("error", ioErr.Error()))
		return "Error reading SQL-file", ioErr
	}

	return buf.String(), nil
}
