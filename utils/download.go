package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/manifoldco/promptui"
)

func FetchAndDownload(url string, headers map[string]string, baseFolder, ep, animeName string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("❌ Error creating request: %v\n", err)
		return
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Error reading response body: %v\n", err)
		return
	}

	htmlData := string(body)
	sources := ScrapeSources(htmlData)
	if len(sources) == 0 {
		fmt.Println("❌ No sources found")
		return
	}

	// Prompt user to select quality
	prompt := promptui.Select{
		Label: "🎥 Select Quality",
		Items: sources,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Quality }}",
			Active:   "🔥 {{ .Quality | cyan }}",
			Inactive: "  {{ .Quality | faint }}",
			Selected: "✅ {{ .Quality | red | cyan }}",
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("❌ Prompt failed %v\n", err)
		return
	}

	selectedSource := sources[i]
	fmt.Printf("📥 Downloading %s quality...\n", selectedSource.Quality)

	animeFolder := filepath.Join(baseFolder, animeName)
	if _, err := os.Stat(animeFolder); os.IsNotExist(err) {
		err = os.Mkdir(animeFolder, os.ModePerm)
		if err != nil {
			fmt.Printf("❌ Error creating folder: %v\n", err)
			return
		}
	}

	filename := fmt.Sprintf("%s - %s.mp4", ep, animeName)
	filePath := filepath.Join(animeFolder, filename)
	err = DownloadFile(filePath, selectedSource.URL, headers)
	if err != nil {
		fmt.Printf("❌ Error downloading file: %v\n", err)
	} else {
		fmt.Println("✅ Download completed!")
	}
}

func DownloadFile(filepath string, url string, headers map[string]string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create a progress bar
	bar := pb.Full.Start64(resp.ContentLength)
	bar.SetTemplateString(`{{ red "🚀 Downloading:" }} {{counters . }} {{bar . "[" "=" ">" "_" "]"}} {{percent . }} {{speed . }}`)
	barReader := bar.NewProxyReader(resp.Body)
	defer bar.Finish()

	// Write the body to file
	_, err = io.Copy(out, barReader)
	if err != nil {
		return err
	}

	return nil
}
