// pkg/mcp/tools_notification.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/notification"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- notification_show ---

type NotificationShowInput struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Subtitle string `json:"subtitle,omitempty"`
}
type NotificationShowOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) notificationShow(_ context.Context, _ *mcp.CallToolRequest, input NotificationShowInput) (*mcp.CallToolResult, NotificationShowOutput, resultFailure) {
	result := s.core.Action("notification.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: notification.TaskSend{Options: notification.NotificationOptions{
			Title:    input.Title,
			Message:  input.Message,
			Subtitle: input.Subtitle,
		}}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, NotificationShowOutput{}, err
		}
		return nil, NotificationShowOutput{}, nil
	}
	return nil, NotificationShowOutput{Success: true}, nil
}

// --- notification_permission_request ---

type NotificationPermissionRequestInput struct{}
type NotificationPermissionRequestOutput struct {
	Granted bool `json:"granted"`
}

func (s *Subsystem) notificationPermissionRequest(_ context.Context, _ *mcp.CallToolRequest, _ NotificationPermissionRequestInput) (*mcp.CallToolResult, NotificationPermissionRequestOutput, resultFailure) {
	result := s.core.Action("notification.request_permission").Run(context.Background(), core.NewOptions())
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, NotificationPermissionRequestOutput{}, err
		}
		return nil, NotificationPermissionRequestOutput{}, nil
	}
	granted, ok := result.Value.(bool)
	if !ok {
		return nil, NotificationPermissionRequestOutput{}, core.E("mcp.notificationPermissionRequest", "unexpected result type", nil)
	}
	return nil, NotificationPermissionRequestOutput{Granted: granted}, nil
}

// --- notification_permission_check ---

type NotificationPermissionCheckInput struct{}
type NotificationPermissionCheckOutput struct {
	Granted bool `json:"granted"`
}

func (s *Subsystem) notificationPermissionCheck(_ context.Context, _ *mcp.CallToolRequest, _ NotificationPermissionCheckInput) (*mcp.CallToolResult, NotificationPermissionCheckOutput, resultFailure) {
	result := s.core.QUERY(notification.QueryPermission{})
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, NotificationPermissionCheckOutput{}, err
		}
		return nil, NotificationPermissionCheckOutput{}, nil
	}
	status, ok := result.Value.(notification.PermissionStatus)
	if !ok {
		return nil, NotificationPermissionCheckOutput{}, core.E("mcp.notificationPermissionCheck", "unexpected result type", nil)
	}
	return nil, NotificationPermissionCheckOutput{Granted: status.Granted}, nil
}

// --- notification_clear ---

type NotificationClearInput struct {
	ID string `json:"id,omitempty"`
}

type NotificationClearOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) notificationClear(_ context.Context, _ *mcp.CallToolRequest, input NotificationClearInput) (*mcp.CallToolResult, NotificationClearOutput, resultFailure) {
	result := s.core.Action("notification.clear").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: notification.TaskClear{ID: input.ID}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, NotificationClearOutput{}, err
		}
		return nil, NotificationClearOutput{}, nil
	}
	return nil, NotificationClearOutput{Success: true}, nil
}

// --- notification_with_actions ---

type NotificationWithActionsInput struct {
	Title      string                            `json:"title"`
	Message    string                            `json:"message"`
	Subtitle   string                            `json:"subtitle,omitempty"`
	CategoryID string                            `json:"category_id,omitempty"`
	Actions    []notification.NotificationAction `json:"actions,omitempty"`
}

type NotificationWithActionsOutput struct {
	Success bool `json:"success"`
}

func (s *Subsystem) notificationWithActions(_ context.Context, _ *mcp.CallToolRequest, input NotificationWithActionsInput) (*mcp.CallToolResult, NotificationWithActionsOutput, resultFailure) {
	result := s.core.Action("notification.send").Run(context.Background(), core.NewOptions(
		core.Option{Key: "task", Value: notification.TaskSend{Options: notification.NotificationOptions{
			Title:      input.Title,
			Message:    input.Message,
			Subtitle:   input.Subtitle,
			CategoryID: input.CategoryID,
			Actions:    input.Actions,
		}}},
	))
	if !result.OK {
		if err, ok := result.Value.(error); ok {
			return nil, NotificationWithActionsOutput{}, err
		}
		return nil, NotificationWithActionsOutput{}, nil
	}
	return nil, NotificationWithActionsOutput{Success: true}, nil
}

// --- Registration ---

func (s *Subsystem) registerNotificationTools(server *mcp.Server) {
	addTool(s, server, &mcp.Tool{Name: "notification_show", Description: "Show a desktop notification"}, s.notificationShow)
	addTool(s, server, &mcp.Tool{Name: "notification_permission_request", Description: "Request notification permission"}, s.notificationPermissionRequest)
	addTool(s, server, &mcp.Tool{Name: "notification_permission_check", Description: "Check notification permission status"}, s.notificationPermissionCheck)
	addTool(s, server, &mcp.Tool{
		Name:        "notification_clear",
		Description: `Clear a notification by id or clear all notifications. Example: {"id":"core-123"}`,
	}, s.notificationClear)
	addTool(s, server, &mcp.Tool{
		Name:        "notification_with_actions",
		Description: `Show an interactive desktop notification with action buttons. Example: {"title":"Deploy","message":"Start deployment?","actions":[{"id":"confirm","label":"Deploy"}]}`,
	}, s.notificationWithActions)
}
