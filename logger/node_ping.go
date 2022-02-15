package logger

import (
	"context"
	"fmt"
)

type INodePing interface {
	LogNodeRequestSend(ctx context.Context, isHealthCheck bool)
	LogNodeRequestHandled(ctx context.Context, isHealthCheck bool)
}

type NodePing struct {
	ID, Address, Blockchain string
}

// LogNodeRequestSend logs a message about sending the request to a node
// it marks the log entry as a healthcheck if isHealthCheck is true
func (n *NodePing) LogNodeRequestSend(ctx context.Context, isHealthCheck bool) {
	logMessage := fmt.Sprintf("request send to %s-node %s on %s", n.Blockchain, n.ID, n.Address)
	if isHealthCheck {
		Log().Info("healthcheck: " + logMessage)
		return
	}
	LogWithContext(ctx).Info(logMessage)
}

// LogNodeRequestHandled logs a message about finishing handling the request to a node
// it marks the log entry as a healthcheck if isHealthCheck is true
func (n *NodePing) LogNodeRequestHandled(ctx context.Context, isHealthCheck bool) {
	logMessage := fmt.Sprintf("request handled by  %s-node %s on %s", n.Blockchain, n.ID, n.Address)
	if isHealthCheck {
		Log().Info("healthcheck: " + logMessage)
		return
	}
	LogWithContext(ctx).Info(logMessage)
}
