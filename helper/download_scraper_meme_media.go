package helper

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

// DownloadMemeScrapingResultMedia is a function to download meme media from scraping result. For now, works for both meme
// sourced from 9gag.com and 1cak.com.
// Returning the filename of the downloaded media.
func DownloadMemeScrapingResultMedia(ctx context.Context, url, outputPath, referer string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:10.0) Gecko/20100101 Firefox/10.0")
	req.Header.Set("Referer", referer)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer WrapCloser(resp.Body.Close)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("9gag media returning non 200")
	}

	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	filename := outputPath + "/" + GenerateID()

	outfile, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(outfile, resp.Body)
	if err != nil {
		return "", err
	}

	mime, err := mimetype.DetectFile(filename)
	if err != nil {
		return "", err
	}

	newFilename := filename + mime.Extension()
	return newFilename, os.Rename(filename, newFilename)
}
