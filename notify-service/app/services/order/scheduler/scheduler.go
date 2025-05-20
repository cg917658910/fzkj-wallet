package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
)

var logger = log.DLogger()

// Scheduler 代表调度器的接口类型。
type Scheduler interface {
	// Init 用于初始化调度器。
	Init() (err error)
	// Start 用于启动调度器并执行爬取流程。
	Start() (err error)
	// Stop 用于停止调度器的运行。
	// 所有处理模块执行的流程都会被中止。
	Stop() (err error)
	// Status 用于获取调度器的状态。
	Status() Status
	// ErrorChan 用于获得错误通道。
	ErrorChan() <-chan error
}

// NewScheduler 会创建一个调度器实例。
func NewScheduler() Scheduler {
	return &myScheduler{}
}

type myScheduler struct {
	// ctx 代表上下文，用于感知调度器的停止。
	ctx context.Context
	// cancelFunc 代表取消函数，用于停止调度器。
	cancelFunc context.CancelFunc
	// status 代表状态。
	status Status
	// statusLock 代表专用于状态的读写锁。
	statusLock sync.RWMutex
}

func (sched *myScheduler) Init() (err error) {
	// 检查状态。
	logger.Info("Check status for initialization...")
	var oldStatus Status
	oldStatus, err =
		sched.checkAndSetStatus(SCHED_STATUS_INITIALIZING)
	if err != nil {
		return
	}
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_INITIALIZED
		}
		sched.statusLock.Unlock()
	}()
	// TODO

	return nil
}

func (sched *myScheduler) Start() (err error) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("Fatal scheduler error: %s", p)
			logger.Fatal(errMsg)
			err = genError(errMsg)
		}
	}()
	logger.Info("Start scheduler...")
	// 检查状态。
	logger.Info("Check status for start...")
	var oldStatus Status
	oldStatus, err =
		sched.checkAndSetStatus(SCHED_STATUS_STARTING)
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_STARTED
		}
		sched.statusLock.Unlock()
	}()
	if err != nil {
		return
	}
	// TODO

	return nil
}
func (sched *myScheduler) Stop() (err error) {
	logger.Info("Stop scheduler...")
	// 检查状态。
	logger.Info("Check status for stop...")
	var oldStatus Status
	oldStatus, err =
		sched.checkAndSetStatus(SCHED_STATUS_STOPPING)
	defer func() {
		sched.statusLock.Lock()
		if err != nil {
			sched.status = oldStatus
		} else {
			sched.status = SCHED_STATUS_STOPPED
		}
		sched.statusLock.Unlock()
	}()
	if err != nil {
		return
	}
	sched.cancelFunc()
	// TODO
	logger.Info("Scheduler has been stopped.")
	return nil
}
func (sched *myScheduler) Status() Status {
	logger.Debug("Scheduler Status")
	var status Status
	sched.statusLock.RLock()
	status = sched.status
	sched.statusLock.RUnlock()
	return status
}
func (sched *myScheduler) ErrorChan() <-chan error {
	logger.Debug("Scheduler ErrorChan")
	errCh := make(chan error, 1)
	return errCh
}

// checkAndSetStatus 用于状态的检查，并在条件满足时设置状态。
func (sched *myScheduler) checkAndSetStatus(
	wantedStatus Status) (oldStatus Status, err error) {
	sched.statusLock.Lock()
	defer sched.statusLock.Unlock()
	oldStatus = sched.status
	err = checkStatus(oldStatus, wantedStatus, nil)
	if err == nil {
		sched.status = wantedStatus
	}
	return
}
