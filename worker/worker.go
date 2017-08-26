package worker

const (
	StartSignal = iota
	RestartSignal
	StopSignal
	KillSignal
)

type RoutineWorker interface {
	Start() error
	Stop()
	Name() string
}

func ManageWorker(inChan chan int, outChan chan interface{}, worker RoutineWorker) {
	innerChan := make(chan bool)
	innerStatus := false
	for {
		select {
		case signal := <-inChan:
			switch signal {
			case StartSignal:
				if !innerStatus {
					go wrapWithStatusManagement(worker.Start, outChan, innerChan)
					innerStatus = true
				}
			case RestartSignal:
				if innerStatus {
					worker.Stop()
					<-innerChan
				}
				go wrapWithStatusManagement(worker.Start, outChan, innerChan)
				innerStatus = true
			case StopSignal:
				if innerStatus {
					worker.Stop()
					<-innerChan
					innerStatus = false
				}
			case KillSignal:
				if innerStatus {
					worker.Stop()
					<-innerChan
					innerStatus = false
				}
				close(innerChan)
				close(inChan)
				if outChan != nil {
					close(outChan)
				}
				return
			}
		}
	}
}

func wrapWithStatusManagement(f func() error, outChan chan interface{}, ch chan bool) {
	if outChan != nil {
		outChan <- true
	}
	defer func() {
		if outChan != nil {
			outChan <- false
		}
		ch <- true
	}()
	err := f()
	if err != nil {
		outChan <- err
	}
}
