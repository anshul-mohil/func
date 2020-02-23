/*
Copyright © 2020 Anshul Mohil <anshulmohil.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"func/util"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// move represents the config command
var sub = &cobra.Command{
	Use:   "sub",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("API_KEY") == "" {
			log.Fatal("API_KEY environment variable needs to be set to access subtitle API")
		}
		//todo: need to add check to know when API_KEY is wrong....
		listOfttIds := util.GetImdbIdList(args[0], "")
		if listOfttIds == nil {
			log.Fatal("Unable to get any Imdb ttIds for requested movie: ", args[0])
		}
		fmt.Println(listOfttIds.Search[0].ImdbID)
		downloadSubtitle(listOfttIds.Search[0].ImdbID, "/Users/anshulmohil/Downloads/anukulmovies", listOfttIds.Search[0].Title)
	},
}

func downloadSubtitle(imdbId string, basePath string, movieName string) {
	body, results := getSubtitileMetadata(imdbId)

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(body), &results)
	fmt.Println("Below are the details of srts downloaded for movie requested")
	var fullFilePath string
	var files []string
	for jsonArrayIndex, result := range results {
		printSrtMeta(result, jsonArrayIndex)
		url := result["ZipDownloadLink"].(string)

		fullFilePath = basePath + "/" + movieName + "/"
		createPathIfNotExist(fullFilePath)
		fileName := strings.Replace(result["SubFileName"].(string), ".srt", ".zip", 1)
		if err := DownloadFile(fullFilePath+fileName, url); err != nil {
			panic(err)
		}
	}
	err := filepath.Walk(fullFilePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
		Unzip(file, fullFilePath)
	}
	zipFiles, err := WalkMatch(fullFilePath, "*.zip")
	nfoFiles, err := WalkMatch(fullFilePath, "*.nfo")

	for _, file := range zipFiles {
		fmt.Println("Removing file", file)
		os.Remove(file)
	}
	for _, file := range nfoFiles {
		fmt.Println("Removing file", file)
		os.Remove(file)
	}

	//SliceIndex(len(files), func(i int) bool { return files[i] == "" })
}
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}
func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
func createPathIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
func getSubtitileMetadata(imdbId string) ([]byte, []map[string]interface{}) {
	url := "https://rest.opensubtitles.org/search/imdbid-" + imdbId + "/sublanguageid-eng"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "TemporaryUserAgent")
	//req.Header.Add("x-rapidapi-host", "opensubtitles-subtitle-tools.p.rapidapi.com")
	//req.Header.Add("x-rapidapi-key", "ed6250409dmsh5b1a4c0a2ab6f01p1240c7jsnab8f132961a4")

	res, _ := http.DefaultClient.Do(req)
	//res, _ :=http.Get("http://example.com/")

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var results []map[string]interface{}
	return body, results
}

func printSrtMeta(result map[string]interface{}, jsonArrayIndex int) {
	fmt.Println("========================================")
	fmt.Println("Name:", result["MovieReleaseName"])
	fmt.Println("Year:", result["MovieYear"])
	fmt.Println("MovieKind:", result["MovieKind"])
	fmt.Println("srt file name:", result["SubFileName"])
	fmt.Println("zip download link:", result["ZipDownloadLink"])
	fmt.Println("ZipDownloadLink:", result["ZipDownloadLink"])
	fmt.Println("isFormatSRT:", result["SubFormat"] == "srt")
	fmt.Println("Imdb Rating:", result["MovieImdbRating"])
	fmt.Println("srt Rating:", result["Score"])
	fmt.Println(jsonArrayIndex)
}

//Todo: Need to implement sort correctly...
func getTopStrFiles(topN int, resultJson []map[string]interface{}, movieName string) []map[string]interface{} {
	//	var filteredMetaObject []map[string]interface{}
	topList := make([]float64, int(topN))
	for jsonArrayIndex, result := range resultJson {
		if strings.ContainsAny(strings.ToLower(result["MovieReleaseName"].(string)), strings.ToLower(movieName)) && result["SubFormat"] == "srt" {
			fmt.Println(jsonArrayIndex)
			//if(result["Score"])
			fmt.Println(topList)
		}
	}
	//	sort.Slice(resultJson, func(i, j int) bool {return resultJson[i].["Score"].(float64)() < resultJson[j].["Score"].(float64)})

	return nil
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
func init() {
	rootCmd.AddCommand(sub)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// move.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// move.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//
