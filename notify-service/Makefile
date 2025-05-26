run:
	go run ./cmd/cli/main.go
start-notify-order:
	curl http://localhost:5080/order/notify/start
stop-notify-order:
	curl http://localhost:5080/order/notify/stop
truss:
	truss proto/notifysvc.proto  --pbpkg github.com/cg917658910/fzkj-wallet/notify-service/proto --svcpkg github.com/cg917658910/fzkj-wallet/notify-service --svcout . -v 
cleanlog:
	rm -rf ./tmp/log	
consumer-send-num:
	grep -o "消费者收到消息" ./tmp/log/notify/consumer.log | wc -l
caller-received-num:
	grep -o "Caller Received msg" ./tmp/log/notify/caller.log | wc -l
caller-notify-successed-num:
	grep -o "Caller Notify success" ./tmp/log/notify/caller.log | wc -l
caller-notify-channel-close-error-num:
	grep -o "send notify result but channel closed" ./tmp/log/notify/caller.log | wc -l
producer-receive-num:
	grep -o "producer received" ./tmp/log/notify/producer.log | wc -l
producer-send-num:
	grep -o "生产者发送消息" ./tmp/log/notify/producer.log | wc -l
producer-send-failed-num:
	grep -o "生产者发送消息失败" ./tmp/log/notify/producer.log | wc -l
consumer-mark-sum:
	grep -o "消费者Mark消息" ./tmp/log/notify/consumer.log | wc -l
count-notify: consumer-send-num \
	caller-received-num \
	caller-notify-channel-close-error-num \
	caller-notify-successed-num \
	producer-receive-num \
	producer-send-num \
	producer-send-failed-num \
	consumer-mark-sum

