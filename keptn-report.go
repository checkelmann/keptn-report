package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/wcharczuk/go-chart"
)

// KeptnReport Keptn Evaluation Result
// Struct geberated with https://mholt.github.io/json-to-go/
type KeptnReport struct {
	Contenttype string `json:"contenttype"`
	Data        struct {
		Deploymentstrategy string `json:"deploymentstrategy"`
		Evaluationdetails  struct {
			IndicatorResults []struct {
				Score   int    `json:"score"`
				Status  string `json:"status"`
				Targets []struct {
					Criteria    string `json:"criteria"`
					TargetValue int    `json:"targetValue"`
					Violated    bool   `json:"violated"`
				} `json:"targets"`
				Value struct {
					Metric  string  `json:"metric"`
					Success bool    `json:"success"`
					Value   float64 `json:"value"`
				} `json:"value"`
			} `json:"indicatorResults"`
			Result         string    `json:"result"`
			Score          int       `json:"score"`
			SloFileContent string    `json:"sloFileContent"`
			TimeEnd        time.Time `json:"timeEnd"`
			TimeStart      time.Time `json:"timeStart"`
		} `json:"evaluationdetails"`
		Labels            interface{} `json:"labels"`
		Project           string      `json:"project"`
		Result            string      `json:"result"`
		Service           string      `json:"service"`
		Stage             string      `json:"stage"`
		Teststrategy      string      `json:"teststrategy"`
		EvaluationHistory []struct {
			Contenttype string `json:"contenttype"`
			Data        struct {
				Deploymentstrategy string `json:"deploymentstrategy"`
				Evaluationdetails  struct {
					IndicatorResults []struct {
						Score   int    `json:"score"`
						Status  string `json:"status"`
						Targets []struct {
							Criteria    string `json:"criteria"`
							TargetValue int    `json:"targetValue"`
							Violated    bool   `json:"violated"`
						} `json:"targets"`
						Value struct {
							Metric  string  `json:"metric"`
							Success bool    `json:"success"`
							Value   float64 `json:"value"`
						} `json:"value"`
					} `json:"indicatorResults"`
					Result         string    `json:"result"`
					Score          int       `json:"score"`
					SloFileContent string    `json:"sloFileContent"`
					TimeEnd        time.Time `json:"timeEnd"`
					TimeStart      time.Time `json:"timeStart"`
				} `json:"evaluationdetails"`
				Labels       interface{} `json:"labels"`
				Project      string      `json:"project"`
				Result       string      `json:"result"`
				Service      string      `json:"service"`
				Stage        string      `json:"stage"`
				Teststrategy string      `json:"teststrategy"`
			} `json:"data"`
			ID             string    `json:"id"`
			Source         string    `json:"source"`
			Specversion    string    `json:"specversion"`
			Time           time.Time `json:"time"`
			Type           string    `json:"type"`
			Shkeptncontext string    `json:"shkeptncontext"`
		} `json:"evaluationHistory"`
	} `json:"data"`
	ID             string    `json:"id"`
	Source         string    `json:"source"`
	Specversion    string    `json:"specversion"`
	Time           time.Time `json:"time"`
	Type           string    `json:"type"`
	Shkeptncontext string    `json:"shkeptncontext"`
	Icon           string    `json:"icon"`
	Label          string    `json:"label"`
}

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func createReport(sourceFile string, reportFile string) {

	// Use the struct
	var jsonReport KeptnReport

	// Load Json to the struct
	json.Unmarshal([]byte(sourceFile), &jsonReport)

	// Array with the Timestamps for the X-Axis of the chart
	var TimeValues []time.Time
	// Map for the data of the evaluation metrics and values for the Y-Axis of the chart
	mYValues := make(map[string][]float64)

	// keptn logo in PNG Format as Base64 encoded String
	logo := "iVBORw0KGgoAAAANSUhEUgAAAGMAAAAfCAYAAADz23MvAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAZESURBVGhD7VlNVhtHEK4eLGBnH4Eb4JwAe5PAyvIN8CoPsQBOEDgBsAC/rHBOYLKSko3xCQwnwD5B7B0Co0591V2ammF+NEgv8YP53psnqWda01NfV31f9zhqiq0Pz+j6ZpnmOxd0+PJrbG0xA0xOxq/9JZqjLf62Ts49I++ZCHdAnYXDlpTZoJ4MkJC43yhhEorg/Rs6XnsXf7WYAuVk9P56zoHeKiUB8HRIx6vbUrpuri7Ju1N6u/omnm3REEn8TLHZf0G9wQdy/lM1Ef6jEAFcX/H1XLqcD+SCHPxPi0ZIyUAAQQI5DixVB9L7C+osduX7xuCEyxiyiMlZC+TdXJ3I/2z0d+V3i4mQknFzdZAhwftvNKI/+fMP/vEltoZ2StZFtDcH2yF7+PyYHCbAua4QNr94IG0tJkJKhqfl+A0Bx8xf4vrfldl+tLYk+iDn3Dod/3Iey9C+kOOTrpDT66+L2KPtlmIba08TIEM3/l5hUuHcHhUCGQgASg1gA7kx6I6DAn0Y0Wsm6DSIO72XdiUnBH1f2hz3/33ts/SnEZe+Bvg+/IeS0RmP59FlVSDjemg14lwCKbOcA46gaO0HESDOj05EsEd+b9yGoKONaIcz6UzIcT5cB3vcohaBjDrBdpQGE+KMLIKevF0LJKmbgr4crR5kCEN5A7kPCZiovf57nqSfZukaAxner8hnwLIEMyzkduToLAYL2xuwyKs4LwTnlLqpC+4T3dQQGRUIU/v7sIBdiK484wwRyLB/itmMmQ5CMMtxqDg72spoiroptHUWwwwRwjjTMoQ9PjG+D1xIM14T5IG9p1v6KegHxFk1wb8UTdB+IIKSF0HEQZg7ybSBMAj70epk+2CbAx+/0Z0+Yhzc0/iLNeznj/FbCujTk2SFx8Cl1Z/zM1zUlkm4N8Ddclmd+8pa903GXoZe/4yfM1aTGI862LHrfUa3X+zYnIgz7GgeIs6sCcgQbHUETQj7UJYc24ZVuyAOEG4KJsC21aGMjMyEYOT3xGScXB6L9M/7A87cPclmC9T8slKjG6HHq3uxJUWGjAJkxi2VorwyeOKYuB2QnxQOvkyc8fB5cda21MJm3ZRiVGMSqjAREZgwJfdwbjvsCuRQVfNxL0e7ItRTwZffA5Ax87PxMzAZOYZrxdm4Kd0krHJTY1Q8eBXqiAAkI8bnL3hsr2V24lNKJgOCK+ueMnDm6oGqkOlXua3DJsf0laMU6bXZsSGrdxMZvAIn8+IctjrCjCtyU9fD/VLCLBw1t7d3iIiZaCEZaDIC48faB8AnFqApwviKgGzWA1XB9qsqM9Al2xdHKcy1GBsWzGP45+ymkuCGAF05l211qJvK702hLU+YBYjqLNxj09DbGc9lssAm+1F6L2RrXqxtcJzZ8qkD+ulExRjw/LPG/IIZm1tJomvAQ6a1Xrc60C4uCTPUbHWgLaR8IAzOyRJmgfOabU2hi027I1yFhF6JAcgfCrt4nQguDZZv2ncC5GIS1hlI/XytR9204hxmaJE43yUsAyaqzlrWAbpWttJtHOAGEEf13yGQoWARGYuz2toicdbykXdTgbAUIrYVfr0O6K9AtgbCs/BGizAeCHfV0Qj3NB33RJYM1HUQUSTOuiUSvPySXJd/02ehVngaoD8CDAj5nLVCvIWzZL+Kn9MD93FmmyjBekBh7jmNZc8hSwZqGN5hFIkz2mB14Vysm1L7a2Gd1bQA4SAWwH1AvIU4p/jyC5ME47lDGKNu59ieR39MRJ1giEHWCKTZiBJq71d07wmRJUMRBhbrvxFnJUfdlLZlYKzwrICsVGcDQhBwC8+LMwXGgwUgXiGnxyU9cZfximLgPFbkuF7eqZjnci7rBL/7YJ0BTM6bYeinfe/5yqCYDAguyoPWfBFsFwJg3/RpmwJEqRWeJfB/IFgmAgOBsgsxKWdGXzCjEaT0CMGpCxKIxvUWQYeyL7okPvZ+/P96LyBfKSbEpJt32+yuOP24VuqCatxmMKJ3UzsnG2TdklFgUth1Bd6xW+IR7DnHOiZbEJhATzloyCjU+DMm9DRzvbW9IbgIphL2maPDLtOUpzzEUHjOoPGWB549ex9UD7XFRfExz9vQXTwwWDIaO63Zo7hMtfhf0JLxA6El44cB0b/wbe6oMrAKDAAAAABJRU5ErkJggg=="

	// Decode logo from base64 to bytes
	logoBytes, _ := base64.StdEncoding.DecodeString(logo)
	// Load Logo into a byte-buffer
	logoBuffer := bytes.NewBuffer(logoBytes)

	// Loop over the Evaluation History from the report
	for _, history := range jsonReport.Data.EvaluationHistory {

		// Check if we have any results in the EvaluationDetails
		if len(history.Data.Evaluationdetails.IndicatorResults) != 0 {

			// Add the evaluation time to the TimesValue Array
			TimeValues = append(TimeValues, history.Time)

			// Loop over the Indicator Results
			for _, historyResult := range history.Data.Evaluationdetails.IndicatorResults {

				// Append results to the mYValues map like mYValues['response_time_95'][1.0123]
				mYValues[historyResult.Value.Metric] = append(mYValues[historyResult.Value.Metric], historyResult.Value.Value)
			}
		}
	}

	// Create a new Chart
	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:    20,
				Left:   20,
				Right:  20,
				Bottom: 100,
			},
		},
		// Leave the Series empty, we will add them later
		Series: []chart.Series{},
	}

	// Loop Over the mYValues
	var series []chart.Series
	// k is the key like response_time_95
	for k := range mYValues {
		series = append(series, chart.TimeSeries{
			Name:    k,
			XValues: TimeValues,
			YValues: mYValues[k],
		})
	}

	// Add Series to the chart
	graph.Series = series
	graph.Elements = []chart.Renderable{
		chart.LegendThin(&graph),
	}

	// Render to buffer to an image
	chartBuffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, chartBuffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Add a new PDF Page
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 16)

	// Add project, service and test-information
	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	pdf.CellFormat(207, 7, "Report", "0", 0, "CM", false, 0, "")
	pdf.SetFont("Courier", "", 10)

	pdf.SetXY(10, 25)

	pdf.CellFormat(190, 4, "Project.: "+jsonReport.Data.Project, "0", 1, "LM", false, 0, "")
	pdf.CellFormat(190, 4, "Service.: "+jsonReport.Data.Service, "0", 1, "LM", false, 0, "")
	pdf.CellFormat(190, 4, "Stage...: "+jsonReport.Data.Stage, "0", 1, "LM", false, 0, "")
	pdf.CellFormat(190, 4, "Strategy: "+jsonReport.Data.Teststrategy, "0", 1, "LM", false, 0, "")

	// Set the Text Color to green/red when the test was passed/failed
	if jsonReport.Data.Result == "pass" {
		pdf.SetTextColor(30, 125, 30)
	} else {
		pdf.SetTextColor(255, 0, 0)
	}
	pdf.CellFormat(190, 4, "Result..: "+jsonReport.Data.Result, "0", 1, "LM", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(190, 4, "Date....: "+jsonReport.Time.String(), "0", 1, "LM", false, 0, "")

	// Show Indicator Results
	pdf.SetXY(10, 55)
	for _, indicator := range jsonReport.Data.Evaluationdetails.IndicatorResults {
		var valueString string = strconv.FormatFloat(indicator.Value.Value, 'f', -1, 64)
		pdf.CellFormat(190, 4, indicator.Value.Metric+":\t "+valueString+"\t ("+indicator.Status+")",
			"0", 1, "LM", false, 0, "")

		for _, target := range indicator.Targets {
			var targetString string
			if target.Violated {
				targetString = "   Test criteria " + target.Criteria + " violated with " + strconv.FormatInt(int64(target.TargetValue), 10)
			} else {
				targetString = "   Test criteria " + target.Criteria + " passed with " + strconv.FormatInt(int64(target.TargetValue), 10)
			}
			pdf.CellFormat(190, 4, targetString,
				"0", 1, "LM", false, 0, "")
		}
	}

	// Add keptn logo image from buffer to the pdf
	var logoOptions gofpdf.ImageOptions
	logoOptions.ImageType = "png"
	logoOptions.AllowNegativePosition = true
	_ = pdf.RegisterImageOptionsReader("logo", logoOptions, logoBuffer)
	pdf.ImageOptions(
		"logo",
		77, 8,
		0, 0,
		false,
		logoOptions,
		0,
		"",
	)

	// Add chart image from buffer to the pdf
	// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	var chartOptions gofpdf.ImageOptions
	chartOptions.ImageType = "png"
	chartOptions.AllowNegativePosition = true
	_ = pdf.RegisterImageOptionsReader("chart", chartOptions, chartBuffer)

	pdf.ImageOptions(
		"chart",
		0, 100,
		200, 0,
		false,
		chartOptions,
		0,
		"",
	)

	// Save PDF Document
	fmt.Println("Saving to " + reportFile)
	err = pdf.OutputFileAndClose(reportFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	var isPipe bool
	var output []rune
	var jsonText string
	fmt.Println("keptn-report")

	// Parsing command line arguments
	sourceFile := flag.String("jsonfile", "", "keptn-evaluation json payload as file")
	reportFile := flag.String("out", "report.pdf", "report file name")
	flag.Parse()

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	// If we are getting the json data via pipe
	if info.Mode()&os.ModeNamedPipe != 0 {
		isPipe = true

		// Start a pipe reader
		reader := bufio.NewReader(os.Stdin)

		for {
			input, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				// When the pipe ends, break
				break
			}
			// Append data to rune
			output = append(output, input)
		}
	} else {
		isPipe = false
	}

	if !isPipe {
		// When the data is not from the pipe, use the file from the cli args
		content, err := ioutil.ReadFile(*sourceFile)
		if err != nil {
			log.Fatal(err)
		}
		jsonText = string(content)
	} else {
		// When the data was readed from the pipe, use this json data
		// Iterate over the rune and append the text as string
		for _, v := range output {
			jsonText += string(v)
		}
	}

	// Create the Report
	createReport(jsonText, *reportFile)
}
