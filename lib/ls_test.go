package lib

import (
	"fmt"
	"sync"
	"testing"
)

type FailInfo struct {
	Failed                 bool
	FailedForActualValue   string
	FailedForExpectedValue string
}

type InspectionReporter struct {
	ExpectedReportingSequence []string
	T                         *testing.T

	CheckSatelliteCountDown bool
	CheckCountDown          bool
	CheckControlTransfer    bool
	CheckCompletion         bool

	currentIndex int

	FailInfo
}

func (r *InspectionReporter) handleExpectation(actualValue string) {
	if r.Failed {
		return
	}

	if r.currentIndex >= len(r.ExpectedReportingSequence) {
		r.Failed = true
		r.FailedForActualValue = actualValue
		r.FailedForExpectedValue = ""
		return
	}

	if actualValue != r.ExpectedReportingSequence[r.currentIndex] {
		r.Failed = true
		r.FailedForActualValue = actualValue
		r.FailedForExpectedValue = r.ExpectedReportingSequence[r.currentIndex]
	} else {
		r.currentIndex++
	}
}

func (r *InspectionReporter) ReportSatteliteCountDown(lsID, satID, count int) {
	if !r.CheckSatelliteCountDown {
		return
	}

	actualValue := fmt.Sprintf("ReportSatteliteCountDown:%d,%d,%d", lsID, satID, count)
	r.handleExpectation(actualValue)
}

func (r *InspectionReporter) ReportCountDown(lsID, count int) {
	if !r.CheckCountDown {
		return
	}
	actualValue := fmt.Sprintf("ReportCountDown:%d,%d", lsID, count)
	r.handleExpectation(actualValue)
}

func (r *InspectionReporter) ReportControlTransfer(fromLS, toLS, numberOfSatelliteLaunched int) {
	if !r.CheckControlTransfer {
		return
	}

	actualValue := fmt.Sprintf("ReportControlTransfer:%d,%d,%d", fromLS, toLS, numberOfSatelliteLaunched)
	r.handleExpectation(actualValue)
}

func (r *InspectionReporter) ReportCompletion(lsID, totalSatelliteLaunched, numberOfSatelliteLaunched int) {
	if !r.CheckCompletion {
		return
	}

	actualValue := fmt.Sprintf("ReportCompletion:%d,%d,%d", lsID, totalSatelliteLaunched, numberOfSatelliteLaunched)
	r.handleExpectation(actualValue)
}

func TestISRO(t *testing.T) {
	t.Run("CompletionTest", func(t *testing.T) {
		isro := ISRO{}
		expectedSequence := []string{
			"ReportCompletion:31,123,3",
		}

		inspectionReporter := &InspectionReporter{ExpectedReportingSequence: expectedSequence, CheckCompletion: true}

		isro.Init(119, 4, inspectionReporter)
		isro.StartLaunch(123)

		if inspectionReporter.Failed {
			t.Errorf("Expected: %s, Received: %s", inspectionReporter.FailedForExpectedValue, inspectionReporter.FailedForActualValue)
		}

	})
}

func TestLS(t *testing.T) {
	t.Run("ControlTransfer_1", func(t *testing.T) {
		ls := LS{}
		nextLS := LS{}
		completionWaitGroup := &sync.WaitGroup{}
		expectedSequence := []string{
			"ReportControlTransfer:0,1,4",
			"ReportControlTransfer:1,0,4",
			"ReportControlTransfer:0,1,4",
			"ReportControlTransfer:1,0,4",
			"ReportCompletion:0,20,4",
		}
		inspectionReporter := InspectionReporter{T: t, ExpectedReportingSequence: expectedSequence, CheckControlTransfer: true, CheckCompletion: true}

		ls.Init(0, &nextLS, 10, 4, &inspectionReporter, completionWaitGroup)
		nextLS.Init(1, &ls, 10, 4, &inspectionReporter, completionWaitGroup)

		completionWaitGroup.Add(1)
		ls.LaunchBatch(0, 20)
		completionWaitGroup.Wait()

		if inspectionReporter.Failed {
			t.Errorf("Expected: %s, Received: %s", inspectionReporter.FailedForExpectedValue, inspectionReporter.FailedForActualValue)
		}
	})

	t.Run("ControlTransfer_2", func(t *testing.T) {
		ls := LS{}
		nextLS := LS{}
		completionWaitGroup := &sync.WaitGroup{}
		expectedSequence := []string{
			"ReportControlTransfer:0,1,4",
			"ReportControlTransfer:1,0,4",
			"ReportControlTransfer:0,1,4",
			"ReportControlTransfer:1,0,4",
			"ReportControlTransfer:0,1,4",
			"ReportCompletion:1,21,1",
		}
		inspectionReporter := InspectionReporter{T: t, ExpectedReportingSequence: expectedSequence, CheckControlTransfer: true, CheckCompletion: true}

		ls.Init(0, &nextLS, 10, 4, &inspectionReporter, completionWaitGroup)
		nextLS.Init(1, &ls, 10, 4, &inspectionReporter, completionWaitGroup)

		completionWaitGroup.Add(1)
		ls.LaunchBatch(0, 21)
		completionWaitGroup.Wait()

		if inspectionReporter.Failed {
			t.Errorf("Expected: %s, Received: %s", inspectionReporter.FailedForExpectedValue, inspectionReporter.FailedForActualValue)
		}
	})

	t.Run("CountDown", func(t *testing.T) {
		ls := LS{}
		completionWaitGroup := &sync.WaitGroup{}
		expectedSequence := []string{
			"ReportCountDown:0,5",
			"ReportCountDown:0,4",
			"ReportCountDown:0,3",
			"ReportCountDown:0,2",
			"ReportCountDown:0,1",
			"ReportCountDown:0,0",
			"ReportCompletion:0,1,1",
		}
		inspectionReporter := InspectionReporter{T: t, ExpectedReportingSequence: expectedSequence, CheckCountDown: true, CheckCompletion: true}

		ls.Init(0, nil, 5, 4, &inspectionReporter, completionWaitGroup)

		completionWaitGroup.Add(1)
		ls.LaunchBatch(0, 1)
		completionWaitGroup.Wait()

		if inspectionReporter.Failed {
			t.Errorf("Expected: %s, Received: %s", inspectionReporter.FailedForExpectedValue, inspectionReporter.FailedForActualValue)
		}
	})
}
