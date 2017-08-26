package worker

const (
	StartSignal = iota
	RestartSignal
	StopSignal
	KillSignal
	FlushSignal
)

type RoutineWorker interface {
	Start()
	Stop()
}

func ManageWorker(inChan chan int, status *bool, worker RoutineWorker) {
	outChan := make(chan bool)
	innerStatus := false
	for {
		select {
		case signal := <-inChan:
			switch signal {
			case StartSignal:
				if !innerStatus {
					go wrapWithStatusManagement(worker.Start, status, outChan)
					innerStatus = true
				}
			case RestartSignal:
				if innerStatus {
					worker.Stop()
					<-outChan
				}
				go wrapWithStatusManagement(worker.Start, status, outChan)
				innerStatus = true
			case StopSignal:
				if innerStatus {
					worker.Stop()
					<-outChan
					innerStatus = false
				}
			case KillSignal:
				if innerStatus {
					worker.Stop()
					<-outChan
					innerStatus = false
				}
				close(outChan)
				close(inChan)
				return
			case FlushSignal:
			}
		}
	}
}

func wrapWithStatusManagement(f func(), status *bool, ch chan bool) {
	*status = true
	defer func() {
		ch <- true
		*status = false
	}()
	f()
}
