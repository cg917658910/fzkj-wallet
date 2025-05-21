package caller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/enum"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/types"
)

type CallerManger interface{}

type MyCallerManager struct {
	ctx            context.Context
	workerNum      int // 调用者数量 用于 执行 notify 的数量
	msgCh          <-chan *sarama.ConsumerMessage
	workerCh       chan *types.NotifyTask   //本地任务通道
	notifyResultCh chan *types.NotifyResult // 调用者通道 用于发送调用结果至生产者
	maxRetries     int                      // 最大重试次数
}

func NewCallerManager(ctx context.Context, msgCh <-chan *sarama.ConsumerMessage, notifyResultCh chan *types.NotifyResult) *MyCallerManager {
	return &MyCallerManager{
		ctx:            ctx,
		workerNum:      100, //
		maxRetries:     5,
		msgCh:          msgCh,
		notifyResultCh: notifyResultCh,
		workerCh:       make(chan *types.NotifyTask, 100), //需要负责关闭
	}
}

func (m *MyCallerManager) Start() error {
	logger.Info("Starting Caller Manager...")

	if err := m.startReciveMsg(); err != nil {
		logger.Errorf("Failed to start receive msg: %v", err)
		return err
	}

	if err := m.setupWorker(); err != nil {
		logger.Errorf("Failed to setup caller: %v", err)
		return err
	}

	return nil
}

// startReciveMsg 启动消息接收
// 1. 启动消息接收
// 2. 将消息放入 workerCh
func (m *MyCallerManager) startReciveMsg() error {
	go func() {
		for {
			select {
			case msg := <-m.msgCh:
				m.processMsg(msg)
			case <-m.ctx.Done():
				logger.Infoln("Stopping caller...")
				return
			}
		}
	}()
	return nil
}

func (m *MyCallerManager) processMsg(msg *sarama.ConsumerMessage) error {
	// 1. 解析消息
	notifyTask, err := buildMsgToNotifyTask(msg)
	if err != nil {
		logger.Errorf("Failed to build notify task: %v", err)
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

func buildMsgToNotifyTask(msg *sarama.ConsumerMessage) (*types.NotifyTask, error) {
	if msg == nil {
		return nil, errors.New("message is nil")
	}
	notifyMsg := &types.NotifyMessage{
		RawMsg: msg,
	}

	if err := json.Unmarshal(msg.Value, notifyMsg); err != nil {
		logger.Errorf("Failed to unmarshal message: %v", err)
		return nil, err
	}
	notifyTask := &types.NotifyTask{
		NotifyMessage: notifyMsg,
		RetryCount:    0,
	}
	return notifyTask, nil
}

// setupWorker 启动 Notify Worker
func (m *MyCallerManager) setupWorker() error {
	for range m.workerNum {
		go func() {
			for {
				select {
				case task := <-m.workerCh:
					m.processNotifyTask(task)
				case <-m.ctx.Done():
					logger.Info("Stopping caller...")
					return
				}
			}
		}()
	}
	return nil
}

// processNotifyTask 处理 Notify 任务
func (m *MyCallerManager) processNotifyTask(task *types.NotifyTask) error {

	if task == nil {
		return errors.New("task is nil")
	}

	// send notify request
	params := &types.NotifyRequestParams{
		NotifyUrl:  task.Data.NotifyUrl,
		NotifyData: task.Data.Info,
	}
	notifyResp, err := sendNotifyRequest(m.ctx, params)
	if err != nil {
		// TODO: 判断条件是否加入重试队列
		logger.Errorf("Caller Failed to send notify request url: %s err: %v", params.NotifyUrl, err)
		return err
	}
	logger.Infof("Caller Notify Result url: %s, status: %v", params.NotifyUrl, notifyResp.Body)

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

// sendNotifyResult 发送 Notify结果至生产者
func (m *MyCallerManager) sendNotifyResult(msg *types.NotifyResult) {
	// TODO: 判断通道是否关闭
	m.notifyResultCh <- msg
}

var _httpClient = &http.Client{
	Timeout: 5 * time.Second, // 总超时时间
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   2 * time.Second, // TCP连接超时（包括DNS）
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   2 * time.Second,
		ResponseHeaderTimeout: 2 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
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
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
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

	close(m.notifyResultCh)
	logger.Info("Caller Manager stopped successfully")
	return nil
}
