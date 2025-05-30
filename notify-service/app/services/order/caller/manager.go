package caller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/enum"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/repo"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/types"
	"github.com/cg917658910/fzkj-wallet/notify-service/config"
)

var (
	_httpClient = &http.Client{
		Timeout: 5 * time.Second, // 总超时时间
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second, // TCP连接超时（包括DNS）
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   5,
			IdleConnTimeout:       60 * time.Second,
		},
	}
)

const (
	// defaultMaxRetries 默认最大重试次数
	_maxRetries        = 10
	minWorkerNum       = 10
	maxWorkerNum       = 3000
	retryDelayMaxTimeS = 60 //s
)

type CallerManger interface{}

type MyCallerManager struct {
	ctx                  context.Context
	workerNum            uint // 调用者数量 用于 执行 notify 的数量
	msgCh                <-chan *sarama.ConsumerMessage
	workerCh             chan *types.NotifyTask   //本地任务通道
	notifyResultCh       chan *types.NotifyResult // 调用者通道 用于发送调用结果至生产者
	notifyResultChClosed uint32
	maxRetries           uint // 最大重试次数
	retryDelayTimeS      time.Duration
	timers               []*time.Timer
	timerMutex           sync.Mutex
	notifyReop           *repo.NotifyResultRepo
}

func NewCallerManager(ctx context.Context, msgCh <-chan *sarama.ConsumerMessage, notifyResultCh chan *types.NotifyResult) *MyCallerManager {
	workerNum := min(max(config.Configs.OrderNotify.OrderNofifyCallerWorkerNum, minWorkerNum), maxWorkerNum)
	return &MyCallerManager{
		ctx:                  ctx,
		workerNum:            workerNum, //
		maxRetries:           min(config.Configs.OrderNotify.OrderNofifyRetryNum, _maxRetries),
		msgCh:                msgCh,
		notifyResultCh:       notifyResultCh,
		notifyResultChClosed: 0,
		notifyReop:           repo.NewNotifyResultRepo(context.Background()),
		retryDelayTimeS:      time.Second * time.Duration(min(config.Configs.OrderNotify.OrderNofifyRetryDelayTimeS, retryDelayMaxTimeS)),
		workerCh:             make(chan *types.NotifyTask, workerNum*5), //需要负责关闭
	}
}

func (m *MyCallerManager) Start() error {
	logger.Info("Starting Caller Manager...")
	atomic.StoreUint32(&m.notifyResultChClosed, 0)
	if err := m.startReceiveMsg(); err != nil {
		errLogger.Errorf("Failed to start receive msg: %v", err)
		return err
	}

	if err := m.setupWorker(); err != nil {
		errLogger.Errorf("Failed to setup caller: %v", err)
		return err
	}

	return nil
}

// startReceiveMsg 启动消息接收
// 1. 启动消息接收
// 2. 将消息放入 workerCh
func (m *MyCallerManager) startReceiveMsg() error {
	go func() {
		for msg := range m.msgCh {
			m.processMsg(msg)
			logger.Infof("Caller Received msg %s", string(msg.Value))
		}
		/* for {
			select {
			case msg := <-m.msgCh:
				m.processMsg(msg)
				logger.Infof("Caller Received msg %s", string(msg.Value))
			case <-m.ctx.Done():
				return
			}
		} */
	}()
	return nil
}

func (m *MyCallerManager) processMsg(msg *sarama.ConsumerMessage) error {
	// 1. 解析消息
	notifyTask, err := buildMsgToNotifyTask(msg)
	if err != nil {
		errLogger.Errorf("Failed to build notify task: %v", err)
		// send notifyResult invalid params
		notifyResult := &types.NotifyResult{
			NotifyTask: notifyTask,
			Status:     enum.NotifyResultStatusFailedInvalidParams,
			Msg:        enum.NotifyResultStatusFailedInvalidParams.String(),
		}
		m.sendNotifyResult(notifyResult)
		return err
	}
	// 2. 发送消息到 workerCh
	m.sendNotifyTask(notifyTask)

	return nil
}

func (m *MyCallerManager) sendNotifyTask(msg *types.NotifyTask) {
	// TODO: 检测通道是否关闭
	m.workerCh <- msg
}

func buildMsgToNotifyTask(msg *sarama.ConsumerMessage) (notifyTask *types.NotifyTask, err error) {
	if msg == nil {
		err = errors.New("msg is nil")
		return
	}
	notifyMsg := &types.NotifyMessage{
		RawMsg: msg,
	}
	notifyTask = &types.NotifyTask{
		NotifyMessage: notifyMsg,
		RetryCount:    0,
	}
	if err = json.Unmarshal(msg.Value, notifyMsg); err != nil {
		err = fmt.Errorf("failed to unmarshal message: %w", err)
		return
	}
	if err = notifyMsg.Check(); err != nil {
		err = fmt.Errorf("failed to check message: %w", err)
		return
	}

	return
}

// setupWorker 启动 Notify Worker
func (m *MyCallerManager) setupWorker() (err error) {
	for range m.workerNum {
		go func() {
			for task := range m.workerCh {
				m.processNotifyTask(task)
			}
			/* for {
				select {
				case task := <-m.workerCh:
					m.processNotifyTask(task)
				case <-m.ctx.Done():
					return
				}
			} */
		}()
	}
	return nil
}

// processNotifyTask 处理 Notify 任务
func (m *MyCallerManager) processNotifyTask(task *types.NotifyTask) error {

	if task == nil {
		return errors.New("task is nil")
	}
	// 检查是否已经通知过了
	var isSend bool
	if m.notifyReop != nil {
		//isSend, _ = m.notifyReop.ExistsProduceResutl(task.MsgId)
	}
	// 如果已发送
	if isSend {
		notifyResult := &types.NotifyResult{
			NotifyTask: task,
			Status:     enum.NotifyResultStatusAlreadyNotified,
			Msg:        enum.NotifyResultStatusAlreadyNotified.String(),
		}
		m.sendNotifyResult(notifyResult)
		return nil
	}
	// send notify request
	params := &types.NotifyRequestParams{
		NotifyUrl:  task.Data.NotifyUrl,
		NotifyData: task.Data.Info,
	}
	// 发送通知请求
	logger.Infof("Caller send notify request url=%s ", params.NotifyUrl)
	notifyResp, err := sendNotifyRequest(context.Background(), params)
	if err != nil {
		logger.Warnf("Caller Failed to send notify request url=%s err=%v RetryCount=%d MaxRetriesNum=%d", params.NotifyUrl, err, task.RetryCount, m.maxRetries)
		// TODO: 判断条件是否加入重试队列
		if task.RetryCount < m.maxRetries {
			logger.Warnf("Caller send notify request failed url=%s join retry queue retryNum=%d", params.NotifyUrl, task.RetryCount)
			timers := time.AfterFunc(m.retryDelayTimeS, func() {
				task.RetryCount++
				logger.Warnf("Caller send notify request failed url: %s err: %v, retry num %d", params.NotifyUrl, err, task.RetryCount)
				m.sendNotifyTask(task)
			})
			m.addTimers(timers)
			return nil
		}
		notifyResult := &types.NotifyResult{
			NotifyTask: task,
			Status:     enum.NotifyResultStatusFailedMaxRetry,
			Msg:        enum.NotifyResultStatusFailedMaxRetry.String(),
		}
		m.sendNotifyResult(notifyResult)
		return err
	}
	logger.Infof("Caller Notify success url: %s, status: %v", params.NotifyUrl, notifyResp.Body)
	// 写入 notifyResultCh
	notifyResult := &types.NotifyResult{
		NotifyTask: task,
		Result:     notifyResp.Body,
		Status:     enum.NotifyResultStatusSuccessed,
		Msg:        enum.NotifyResultStatusSuccessed.String(),
	}
	m.sendNotifyResult(notifyResult)
	return nil
}

func (m *MyCallerManager) addTimers(timer *time.Timer) {
	m.timerMutex.Lock()
	defer m.timerMutex.Unlock()
	m.timers = append(m.timers, timer)
}

func (m *MyCallerManager) CleanupTimers() {
	m.timerMutex.Lock()
	defer m.timerMutex.Unlock()
	for _, timer := range m.timers {
		timer.Stop()
	}
}

// sendNotifyResult 发送 Notify结果至生产者
func (m *MyCallerManager) sendNotifyResult(msg *types.NotifyResult) {
	if m.resultChClosed() {
		logger.Infof("send notify result but channel closed")
		return
	}
	m.notifyResultCh <- msg
}

func (m *MyCallerManager) canceled() bool {
	select {
	case <-m.ctx.Done():
		logger.Infoln("MyCallerManager ctx Done")
		return true
	default:
		return false
	}
}

func sendNotifyRequest(ctx context.Context, params *types.NotifyRequestParams) (notifyResp *types.NotifyResponse, err error) {

	if params == nil {
		return nil, errors.New("params is nil")
	}
	url := params.NotifyUrl
	if url == "" {
		return nil, errors.New("notify url is empty")
	}
	payload, err := json.Marshal(params.NotifyData)
	if err != nil {
		return nil, errors.New("failed to marshal notify data")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := _httpClient.Do(req) // TODO: 复用 http.Client
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	notifyResp = &types.NotifyResponse{
		HTTPStatus: resp.StatusCode,
		Body:       string(body),
	}
	return notifyResp, nil
}

func (m *MyCallerManager) Stop() error {
	logger.Info("Stopping Caller Manager...")
	//m.CleanupTimers()
	if err := m.closeWorkerCh(); err != nil {
		errLogger.Errorf("Failed to close worker channel: %v", err)
		return err
	}
	if err := m.closeResultCh(); err != nil {
		errLogger.Errorf("Failed to close notify result channel: %v", err)
		return err
	}
	logger.Info("Caller Manager stopped successfully")
	return nil
}

func (m *MyCallerManager) closeWorkerCh() error {
	logger.Infoln("Caller Manager close worker channel...")
	sync.OnceFunc(func() {
		close(m.workerCh)
	})
	logger.Infoln("Caller Manager close worker channel successfully")

	return nil
}

func (m *MyCallerManager) closeResultCh() error {
	logger.Infoln("Caller Manager close result channel...")
	sync.OnceFunc(func() {
		if atomic.CompareAndSwapUint32(&m.notifyResultChClosed, 0, 2) {
			close(m.notifyResultCh)
		}
	})
	logger.Infoln("Caller Manager close result channel successfully")
	return nil
}
func (m *MyCallerManager) resultChClosed() bool {
	logger.Infof("atomic.LoadUint32(&m.notifyResultChClosed)=%d", atomic.LoadUint32(&m.notifyResultChClosed))
	if atomic.LoadUint32(&m.notifyResultChClosed) == 0 {
		return false
	}
	return true
}
