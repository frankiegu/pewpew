package pewpew

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	color "github.com/fatih/color"
)

//creates nice readable summary of entire stress test
func createTextSummary(reqStatSummary requestStatSummary) string {
	summary := "\n"

	summary = summary + "Runtime Statistics:\n"
	summary = summary + "Total time:  " + fmt.Sprintf("%d", reqStatSummary.endTime.Sub(reqStatSummary.startTime).Nanoseconds()/1000000) + " ms\n"
	summary = summary + "Mean RPS:    " + fmt.Sprintf("%.2f", reqStatSummary.avgRPS*1000000000) + " req/sec\n"

	summary = summary + "\nQuery Statistics\n"
	summary = summary + "Mean query:     " + fmt.Sprintf("%d", reqStatSummary.avgDuration/1000000) + " ms\n"
	summary = summary + "Fastest query:  " + fmt.Sprintf("%d", reqStatSummary.minDuration/1000000) + " ms\n"
	summary = summary + "Slowest query:  " + fmt.Sprintf("%d", reqStatSummary.maxDuration/1000000) + " ms\n"

	summary = summary + "Total Data Transferred: " + fmt.Sprintf("%d", reqStatSummary.totalDataTransferred) + " bytes\n"
	summary = summary + "Average Data Transferred:  " + fmt.Sprintf("%d", reqStatSummary.avgDataTransferred) + " bytes\n"

	summary = summary + "\nResponse Codes\n"
	//sort the status codes
	var codes []int
	for key := range reqStatSummary.statusCodes {
		codes = append(codes, key)
	}
	sort.Ints(codes)
	for _, code := range codes {
		if code == 0 {
			continue
		}
		summary = summary + fmt.Sprintf("%d", code) + ": " + fmt.Sprintf("%d", reqStatSummary.statusCodes[code]) + " responses\n"
	}
	if reqStatSummary.statusCodes[0] > 0 {
		summary = summary + "Failed: " + fmt.Sprintf("%d", reqStatSummary.statusCodes[0]) + " requests\n"
	}
	return summary
}

//print colored single line stats per requestStat
func printStat(stat requestStat) {
	if stat.Error != nil {
		color.Set(color.FgRed)
		fmt.Println("Failed to make request: " + stat.Error.Error())
		color.Unset()
	} else {
		if stat.StatusCode >= 100 && stat.StatusCode < 200 {
			color.Set(color.FgBlue)
		} else if stat.StatusCode >= 200 && stat.StatusCode < 300 {
			color.Set(color.FgGreen)
		} else if stat.StatusCode >= 300 && stat.StatusCode < 400 {
			color.Set(color.FgCyan)
		} else if stat.StatusCode >= 400 && stat.StatusCode < 500 {
			color.Set(color.FgMagenta)
		} else {
			color.Set(color.FgRed)
		}
		fmt.Printf("%s %d\t%d bytes\t%d ms\t-> %s %s\n",
			stat.Proto,
			stat.StatusCode,
			stat.DataTransferred,
			stat.Duration.Nanoseconds()/1000000,
			stat.Method,
			stat.URL)
		color.Unset()
	}
}

//print tons of info about the request, response and response body
func printVerbose(req *http.Request, response *http.Response) {
	var requestInfo string
	//request details
	requestInfo = requestInfo + fmt.Sprintf("Request:\n%+v\n\n", &req)

	//reponse metadata
	requestInfo = requestInfo + fmt.Sprintf("Response:\n%+v\n\n", response)

	//reponse body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		requestInfo = requestInfo + fmt.Sprintf("Failed to read response body: %s\n", err.Error())
	} else {
		requestInfo = requestInfo + fmt.Sprintf("Body:\n%s\n\n", body)
		response.Body.Close()
	}
	fmt.Println(requestInfo)
}
