package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"

	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/caller"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/consumer"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/producer"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/types"
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
	//ErrorChan() <-chan error
}

type myScheduler struct {
	consumerManager *consumer.MyConsumerManager
	callerManager   *caller.MyCallerManager
	producerManager producer.ProducerManager
	consumerMsgCh   chan *sarama.ConsumerMessage // 消费者通道 用于接收kafka消息
	notifyResultCh  chan *types.NotifyResult     // 调用者通道 用于发送调用结果至生产者
	markMessageCh   chan *sarama.ConsumerMessage // 用于标记消息的通道
	// ctx 代表上下文，用于感知调度器的停止。
	ctx context.Context
	// cancelFunc 代表取消函数，用于停止调度器。
	cancelFunc context.CancelFunc
	// status 代表状态。
	status Status
	// statusLock 代表专用于状态的读写锁。
	statusLock sync.RWMutex
}

func NewScheduler() Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	consumerMsgCh := make(chan *sarama.ConsumerMessage, 20)
	markMessageCh := make(chan *types.MarkMessageParams, 50)
	notifyResultCh := make(chan *types.NotifyResult, 50) // TODO: close
	return &myScheduler{
		ctx:             ctx,
		cancelFunc:      cancel,
		consumerManager: consumer.NewConsumerManager(ctx, consumerMsgCh, markMessageCh),
		callerManager:   caller.NewCallerManager(ctx, consumerMsgCh, notifyResultCh),
		producerManager: producer.NewProducerManager(ctx, notifyResultCh, markMessageCh),
		consumerMsgCh:   consumerMsgCh,
		notifyResultCh:  notifyResultCh,
	}
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
		logger.Errorf("Failed to check status: %v", err)
		return
	}
	// TODO
	if err := sched.consumerManager.Start(); err != nil {
		logger.Errorf("Failed to start consumer manager: %v", err)
		return err
	}

	if err := sched.callerManager.Start(); err != nil {
		logger.Errorf("Failed to start caller manager: %v", err)
		return err
	}

	if err := sched.producerManager.Start(); err != nil {
		logger.Errorf("Failed to start producer manager: %v", err)
		return err
	}

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
	time.Sleep(time.Second * 1) // TODO:
	// TODO: check
	if err = sched.consumerManager.Stop(); err != nil {
		logger.Errorf("Failed to stop consumer manager: %v", err)
		return
	}
	if err = sched.callerManager.Stop(); err != nil {
		logger.Errorf("Failed to stop caller manager: %v", err)
		return
	}
	if err = sched.producerManager.Stop(); err != nil {
		logger.Errorf("Failed to stop producer manager: %v", err)
		return
	}
	if sched.consumerMsgCh != nil {
		close(sched.consumerMsgCh)
	}
	if sched.notifyResultCh != nil {
		close(sched.notifyResultCh)
	}
	if sched.markMessageCh != nil {
		close(sched.markMessageCh)
	}

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

/* func (sched *myScheduler) ErrorChan() <-chan error {
	logger.Debug("Scheduler ErrorChan")
	errCh := make(chan error, 1)
	return errCh
} */

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

// canceled 用于判断调度器的上下文是否已被取消。
func (sched *myScheduler) canceled() bool {
	select {
	case <-sched.ctx.Done():
		return true
	default:
		return false
	}
}
