package main

import (
	"context"
	"log"
	"time"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_data_accesslog_v3 "github.com/envoyproxy/go-control-plane/envoy/data/accesslog/v3"
	pb "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	"google.golang.org/grpc"
)

func sendTCPAccessLogMessage(stream pb.AccessLogService_StreamAccessLogsClient) {
	log.Println("Sending TCP logs")
	if err := stream.Send(
		&pb.StreamAccessLogsMessage{
			Identifier: &pb.StreamAccessLogsMessage_Identifier{
				Node: &envoy_config_core_v3.Node{
					Id: "test-client",
				},
			},
			LogEntries: &pb.StreamAccessLogsMessage_TcpLogs{
				TcpLogs: &pb.StreamAccessLogsMessage_TCPAccessLogEntries{
					LogEntry: []*envoy_data_accesslog_v3.TCPAccessLogEntry{
						{
							ConnectionProperties: &envoy_data_accesslog_v3.ConnectionProperties{
								ReceivedBytes: 1024,
								SentBytes:     1024,
							},
						},
					},
				},
			},
		},
	); err != nil {
		log.Fatal(err)
	}
}

func sendHTTPAccessLogMessage(stream pb.AccessLogService_StreamAccessLogsClient) {
	log.Println("Sending HTTP Logs")
	if err := stream.Send(&pb.StreamAccessLogsMessage{
		Identifier: &pb.StreamAccessLogsMessage_Identifier{
			Node: &envoy_config_core_v3.Node{
				Id: "test-client",
			},
		},
		LogEntries: &pb.StreamAccessLogsMessage_HttpLogs{
			HttpLogs: &pb.StreamAccessLogsMessage_HTTPAccessLogEntries{
				LogEntry: []*envoy_data_accesslog_v3.HTTPAccessLogEntry{
					{
						Request: &envoy_data_accesslog_v3.HTTPRequestProperties{
							Scheme: "https",
						},
					},
				},
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(":8080", opts...)
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewAccessLogServiceClient(conn)
	stream, err := client.StreamAccessLogs(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for true {
		sendHTTPAccessLogMessage(stream)
		sendTCPAccessLogMessage(stream)
		time.Sleep(2 * time.Second)
	}
	stream.CloseAndRecv()
}
