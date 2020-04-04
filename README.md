# keptn-report

A small utility to render keptn Quality Gates evaluation results to a PDF Report.

![sample](https://i.imgur.com/g4lfo6u.png)

## Installation

Install with go

```bash
go install github.com/checkelmann/keptn-report
```

## Usage

You can either use the command line switches, or directly pipe your curl output into the report generator.

```bash
cat sampledata.json | ./keptn-report -out myreport.pdf
```

```bash
./keptn-report -h
keptn-report
Usage of ./keptn-report:
  -jsonfile string
        keptn-evaluation json payload as file
  -out string
        report file name (default "report.pdf")
```
