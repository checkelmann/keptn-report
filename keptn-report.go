package main

import (
	"bufio"
	"bytes"
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

	var jsonReport KeptnReport

	json.Unmarshal([]byte(sourceFile), &jsonReport)

	var evaluationIterations = 0
	var TimeValues []time.Time
	//var ResultValues []float64
	//var ResultMetrics []string
	//xValues := make([]time.Time, numValues)
	mYValues := make(map[string][]float64)
	//mYValues["Test"] = append(mYValues["Test"], 1)
	//mYValues["Test"] = append(mYValues["Test"], 2)
	/*for _, indicator := range jsonReport.Data.Evaluationdetails.IndicatorResults {
		TimeValues = append(TimeValues, jsonReport.Time)
		ResultValues = append(ResultValues, indicator.Value.Value)
		ResultMetrics = append(ResultMetrics, indicator.Value.Metric)
	}*/

	for _, history := range jsonReport.Data.EvaluationHistory {
		//fmt.Println(len(history.Data.Evaluationdetails.IndicatorResults))
		if len(history.Data.Evaluationdetails.IndicatorResults) != 0 {
			evaluationIterations++
			TimeValues = append(TimeValues, history.Time)
			for _, historyResult := range history.Data.Evaluationdetails.IndicatorResults {
				//ResultValues = append(ResultValues, historyResult.Value.Value)
				//fmt.Printf("%s - %v %s\n", history.Time.String(), evaluationIterations, historyResult.Value.Metric)
				mYValues[historyResult.Value.Metric] = append(mYValues[historyResult.Value.Metric], historyResult.Value.Value)
			}
		}
	}

	graph := chart.Chart{
		Background: chart.Style{
			Padding: chart.Box{
				Top:    20,
				Left:   20,
				Right:  20,
				Bottom: 100,
			},
		},
		Series: []chart.Series{},
	}

	var series []chart.Series
	for k := range mYValues {
		series = append(series, chart.TimeSeries{
			Name:    k,
			XValues: TimeValues,
			YValues: mYValues[k],
		})
	}
	graph.Series = series
	graph.Elements = []chart.Renderable{
		chart.LegendThin(&graph),
	}

	//f, _ := os.Create(reportFile)
	//defer f.Close()

	// Render to buffer
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 16)

	// CellFormat(width, height, text, border, position after, align, fill, link, linkStr)
	pdf.CellFormat(190, 7, "keptn Report", "0", 0, "CM", false, 0, "")
	pdf.SetFont("Courier", "", 10)

	pdf.SetXY(10, 25)
	//pdf.Cell(40, 10, "Project: "+jsonReport.Data.Project)
	//pdf.Cell(40, 10, "Service: "+jsonReport.Data.Service)
	pdf.CellFormat(190, 4, "Project.: "+jsonReport.Data.Project, "0", 1, "LM", false, 0, "")
	pdf.CellFormat(190, 4, "Service.: "+jsonReport.Data.Service, "0", 1, "LM", false, 0, "")
	pdf.CellFormat(190, 4, "Stage...: "+jsonReport.Data.Stage, "0", 1, "LM", false, 0, "")
	pdf.CellFormat(190, 4, "Strategy: "+jsonReport.Data.Teststrategy, "0", 1, "LM", false, 0, "")
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

	// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	var opt gofpdf.ImageOptions
	opt.ImageType = "png"
	opt.AllowNegativePosition = true
	//opt.ReadDpi = true
	_ = pdf.RegisterImageOptionsReader("chart", opt, buffer)

	pdf.ImageOptions(
		"chart",
		0, 100,
		200, 0,
		false,
		opt,
		0,
		"",
	)
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

	sourceFile := flag.String("jsonfile", "", "keptn-evaluation json payload as file")
	reportFile := flag.String("out", "report.pdf", "report file name")
	flag.Parse()

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeNamedPipe != 0 {
		isPipe = true
		reader := bufio.NewReader(os.Stdin)

		for {
			input, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			output = append(output, input)
		}
	} else {
		isPipe = false
	}

	//fmt.Println("Using:", *sourceFile)
	//fmt.Println("Report File:", *reportFile)

	if !isPipe {
		content, err := ioutil.ReadFile(*sourceFile)
		if err != nil {
			log.Fatal(err)
		}
		jsonText = string(content)
	} else {
		for _, v := range output {
			jsonText += string(v)
		}
	}

	createReport(jsonText, *reportFile)
}
