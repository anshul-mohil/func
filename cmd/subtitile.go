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
	"encoding/json"
	"fmt"
	"func/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
		if os.Getenv("API_KEY") == ""{
			log.Fatal("API_KEY environment variable needs to be set to access subtitle API")
		}
		//todo: need to add check to know when API_KEY is wrong....
		listOfSubtitles :=util.GetImdbIdList(args[0],"")
		if listOfSubtitles ==nil{
			log.Fatal("Unable to get any subtitles for requested movie: ",args[0])
		}
		fmt.Println(listOfSubtitles.Search[0].ImdbID)
		downloadSubtitle(listOfSubtitles.Search[0].ImdbID, args[0])
		//subtitileFilePath := ""
		//util.CopyDirectory(subtitileFilePath, destination)
	},
}

func downloadSubtitle(imdbId string,name string){
	url := "https://rest.opensubtitles.org/search/imdbid-"+imdbId+"/sublanguageid-eng"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "TemporaryUserAgent")
	//req.Header.Add("x-rapidapi-host", "opensubtitles-subtitle-tools.p.rapidapi.com")
	//req.Header.Add("x-rapidapi-key", "ed6250409dmsh5b1a4c0a2ab6f01p1240c7jsnab8f132961a4")

	res, _ := http.DefaultClient.Do(req)
	//res, _ :=http.Get("http://example.com/")

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	//body := json.NewDecoder(res.Body)
//	fmt.Println(res)
//	//fmt.Println(string(body))
//	bolB, _ := json.Marshal(decoder)
//	fmt.Println(string(bolB))
	// Declared an empty interface of type Array
	var results []map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(body), &results)
	fmt.Println("Below are the details of srts downloaded for movie requested")
	for jsonArrayIndex, result := range results {
		printSrtMeta(result, jsonArrayIndex)
	}
	//topStrFiles := getTopStrFiles(3, results,name)
	//for jsonArrayIndex, result := range topStrFiles {
	//	printSrtMeta(result, jsonArrayIndex)
	//}

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
func getTopStrFiles(topN int, resultJson []map[string]interface{},movieName string) ([]map[string]interface{}) {
//	var filteredMetaObject []map[string]interface{}
	topList := make([]float64, int(topN))
	for jsonArrayIndex, result := range resultJson {
		if strings.ContainsAny( strings.ToLower(result["MovieReleaseName"].(string)),strings.ToLower(movieName))   && result["SubFormat"] == "srt" {
			fmt.Println(jsonArrayIndex)
			//if(result["Score"])
			fmt.Println(topList)
			}
		}
//	sort.Slice(resultJson, func(i, j int) bool {return resultJson[i].["Score"].(float64)() < resultJson[j].["Score"].(float64)})

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