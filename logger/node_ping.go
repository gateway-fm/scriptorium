package logger

import (
	"context"
	"fmt"
)

type INodePing interface {
	LogNodeRequestSend(ctx context.Context, isHealthCheck bool, bc, id, addr string)
	LogNodeRequestHandled(ctx context.Context, isHealthCheck bool, bc, id, addr string)
}

type NodePing struct{}

// LogNodeRequestSend logs a message about sending the request to a node
// it marks the log entry as a healthcheck if isHealthCheck is true
func (n *NodePing) LogNodeRequestSend(ctx context.Context, isHealthCheck bool, bc, id, addr string) {
	logMessage := fmt.Sprintf("request send to %s-node %s on %s", bc, id, addr)
	LogWithContext(ctx).Info(logMessage)
}

// LogNodeRequestHandled logs a message about finishing handling the request to a node
// it marks the log entry as a healthcheck if isHealthCheck is true
func (n *NodePing) LogNodeRequestHandled(ctx context.Context, isHealthCheck bool, bc, id, addr string) {
	logMessage := fmt.Sprintf("request handled by  %s-node %s on %s", bc, id, addr)
	LogWithContext(ctx).Info(logMessage)
}
