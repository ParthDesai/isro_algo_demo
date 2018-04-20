package lib

import (
	"fmt"
	"strconv"
)

// Reporter interface is used to get event data from LS and isro
type Reporter interface {
	ReportSatteliteCountDown(lsID, satID, count int)
	ReportCountDown(lsID, count int)
	ReportControlTransfer(fromLS, toLS, numberOfSatelliteLaunched int)
	ReportCompletion(lsID, totalSatelliteLaunched, numberOfSatelliteLaunched int)
}

// PrintReporter is implementation of interface Reporter
type PrintReporter struct {
}

// ReportSatteliteCountDown reports satellite countdown with launch station Id,
// satellite id and count
func (r *PrintReporter) ReportSatteliteCountDown(lsID, satID, count int) {
	fmt.Println("LS"+strconv.Itoa(lsID), " ", "SAT"+strconv.Itoa(satID), " ", "CountDown:"+strconv.Itoa(count))
}

// ReportCountDown reports sync countdown of entire batch of satellites
func (r *PrintReporter) ReportCountDown(lsID, count int) {
	fmt.Println("LS"+strconv.Itoa(lsID), " ", "All Sattelites at:", count)
}

// ReportControlTransfer reports control transfer from current launch station to next launch station
func (r *PrintReporter) ReportControlTransfer(fromLS, toLS, numberOfSatelliteLaunched int) {
	fmt.Println("Current batch of satellites:"+strconv.Itoa(numberOfSatelliteLaunched)+" launched from: LS"+strconv.Itoa(fromLS)+",", "Handing over control to: LS"+strconv.Itoa(toLS))
}

// ReportCompletion reports completion of launch of all satellites
func (r *PrintReporter) ReportCompletion(lsID, totalSatelliteLaunched, numberOfSatelliteLaunched int) {
	fmt.Println("Completed: Launched", strconv.Itoa(totalSatelliteLaunched),
		"Last batch of ", strconv.Itoa(numberOfSatelliteLaunched), "was launched by:"+strconv.Itoa(lsID))
}
