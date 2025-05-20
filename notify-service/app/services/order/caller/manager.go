package caller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/data"
)

type CallerManger interface{}

type MyCallerManager struct {
	ctx          context.Context
	callerNum    int // 调用者数量
	msgCh        chan *sarama.ConsumerMessage
	callResultCh chan *data.CallResult
}

func NewCallerManager(ctx context.Context, msgCh chan *sarama.ConsumerMessage, callResultCh chan *data.CallResult) *MyCallerManager {
	return &MyCallerManager{
		ctx:          ctx,
		callerNum:    5,
		msgCh:        msgCh,
		callResultCh: callResultCh,
	}
}

func (m *MyCallerManager) Start() error {
	logger.Info("Starting Caller Manager...")
	if err := m.setupCaller(); err != nil {
		logger.Errorf("Failed to setup caller: %v", err)
		return err
	}
	return nil
}

func (m *MyCallerManager) setupCaller() error {
	for range m.callerNum {
		go func() {
			for {
				select {
				case msg := <-m.msgCh:
					logger.Infof("Received message: %s", string(msg.Value))
					m.handleCall(msg)
				case <-m.ctx.Done():
					logger.Info("Stopping caller...")
					return
				}
			}
		}()
	}
	return nil
}

func (m *MyCallerManager) handleCall(msg *sarama.ConsumerMessage) {
	logger.Infof("Handling call for message: %s", string(msg.Value))
	params := &data.OrderNotifyMessage{}
	if err := json.Unmarshal(msg.Value, params); err != nil {
		logger.Errorf("Failed to unmarshal message: %v", err)
		return
	}
	if err := sendNotifyRequest(m.ctx, params); err != nil {
		logger.Errorf("Failed to send notify request: %v", err)
		return
	}

}

func sendNotifyRequest(ctx context.Context, params *data.OrderNotifyMessage) error {
	url := params.Data.NotifyUrl

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	logger.Info("Notify service response Body:", string(body))
	return nil
}

func (m *MyCallerManager) Stop() error {
	logger.Info("Stopping Caller Manager...")
	close(m.msgCh)
	close(m.callResultCh)
	logger.Info("Caller Manager stopped successfully")
	return nil
}
