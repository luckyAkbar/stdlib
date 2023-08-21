package multimedia

import (
	"os"
)

func generateConcatFileList(list []string, outputPath string) error {
	var data []byte

	for _, item := range list {
		data = append(data, []byte("file '"+item+"'\n")...)
	}

	return os.WriteFile(outputPath, data, 0777)
}
