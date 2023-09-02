package worker

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wanderer69/tools/queue"
)

type Command struct {
	TaskID   string
	Cmd      string
	Sentence string
	Payload  interface{}
}

type CommandAnswer struct {
	TaskID  string
	Payload interface{}
	Result  string
	Err     error
}

type Process struct {
	env      interface{}
	fn       func(interface{}, interface{}) (interface{}, error)
	cmdCh    chan *Command
	answerCh chan *CommandAnswer
}

func NewProcess(env interface{}, fn func(interface{}, interface{}) (interface{}, error)) *Process {
	return &Process{
		env:      env,
		fn:       fn,
		cmdCh:    make(chan *Command),
		answerCh: make(chan *CommandAnswer),
	}
}

func (p *Process) Send(pl interface{}) (string, error) {
	c := &Command{
		Cmd:     "process",
		Payload: pl,
	}
	p.cmdCh <- c
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			return "", errors.New("timeout")
		case ca := <-p.answerCh:
			return ca.TaskID, nil
		}
	}
}

func (p *Process) Check(taskID string) (interface{}, string, error) {
	c := &Command{
		Cmd:    "check_process_finished",
		TaskID: taskID,
	}
	p.cmdCh <- c
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			return "", "error", errors.New("timeout")
		case ca := <-p.answerCh:
			if ca.Err != nil {
				return nil, ca.Result, ca.Err
			}
			if ca.Result == "wait" {
				return nil, ca.Result, nil
			}
			return ca.Payload, ca.Result, nil
		}
	}
}

func (p *Process) Stop() {
	c := &Command{
		Cmd: "exit",
	}
	p.cmdCh <- c
}

func (p *Process) Run() {
	go p.tasker()
}

func (p *Process) tasker() {
	channelIn := make(chan *Command)
	channelOut := make(chan *CommandAnswer)
	channelExit := make(chan bool)
	// обмен с однопоточным почтовым клиентом
	internalProc := func() {
		flag := false
		for {
			select {
			case sri := <-channelIn:
				fmt.Printf("sri %#v\r\n", sri)
				srai := &CommandAnswer{}
				srai.TaskID = sri.TaskID
				if p.fn != nil {
					pl, err := p.fn(p.env, sri.Payload)
					srai.Payload = pl
					srai.Err = err
				}
				channelOut <- srai
			case <-channelExit:
				flag = true
			}
			if flag {
				break
			}
		}
	}
	go internalProc()
	queueIn := queue.NewQueue()
	queueOut := queue.NewQueue() // очередь выполненных заданий
	// в цикле ожидаем  приход команды и интерпретируем ее
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	state := 0
	n := 0
	for {
		select {
		case <-ticker.C:
			n = n + 1
			if n == 60 {
				n = 0
			}
		case commandIn := <-p.cmdCh:
			switch commandIn.Cmd {
			case "process":
				// формируем установку в очередь отправки
				commandIn.TaskID = uuid.NewString()
				qi := &queue.QueueItem{Value: commandIn}
				queueIn.Push(qi)
				// формируем ответ
				ca := CommandAnswer{}
				ca.TaskID = commandIn.TaskID
				p.answerCh <- &ca
			case "check_process_finished":
				// проверяем, что почта отправилась. для этого проверяем состояние по идентификатору
				var ca *CommandAnswer
				qi, err := queueIn.Get()
				if err != nil {
					if !errors.Is(err, queue.ErrQueueEmpty) {
						ca = &CommandAnswer{}
						ca.Result = "error"
						ca.Err = err
						p.answerCh <- ca
						continue
					}
				}
				isAnswered := false
				if qi != nil {
					for {
						if commandIn.TaskID == qi.Value.(*Command).TaskID {
							ca = &CommandAnswer{}
							ca.TaskID = qi.Value.(*Command).TaskID
							ca.Result = "wait"
							isAnswered = true
							break
						}
						qi, err = queueIn.Next(qi)
						if err != nil {
							break
						}
						if qi == nil {
							break
						}
					}
				}
				if !isAnswered {
					// проверяем очередь
					qi, err = queueOut.Get()
					if err != nil {
						continue
					}
					if qi != nil {
						for {
							if commandIn.TaskID == qi.Value.(*CommandAnswer).TaskID {
								ca = qi.Value.(*CommandAnswer)
								queueOut.Delete(qi)
								isAnswered = true
								break
							}
							qi, err := queueOut.Next(qi)
							if err != nil {
								continue
							}
							if qi == nil {
								break
							}
						}
					}
				}
				// формируем ответ
				if !isAnswered {
					ca = &CommandAnswer{}
					ca.Result = "error"
					ca.Err = fmt.Errorf("query with id %v not found", commandIn.TaskID)
				}
				p.answerCh <- ca
			case "check":
				// метод проверки работы сервиса
				result := "Error"
				switch commandIn.TaskID {
				case "Ping":
					result = "Pong"
				default:
					result = "OK"
				}
				ca := &CommandAnswer{}
				ca.Result = result
				p.answerCh <- ca
			case "quit":
				channelExit <- true
			}
		case srai := <-channelOut:
			// ответ от почтового клиента
			state = 0
			qi, err := queueIn.Get()
			if err != nil {
				continue
			}
			if qi != nil {
				qi.Value = srai
				queueOut.Push(qi)
				queueIn.DeleteFirst()
				state = 0
			}
		}
		switch state {
		case 0:
			// нет загрузки почтового клиента. смотрим очередь
			qi, err := queueIn.Get()
			if err != nil {
				continue
			}
			if qi != nil {
				// есть, берем в работу
				sri := qi.Value.(*Command)
				channelIn <- sri
				state = -1
			}
		}
	}
}
