package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorReport 前端错误报告结构
type ErrorReport struct {
	Errors   []ErrorInfo `json:"errors"`
	Metadata Metadata    `json:"metadata"`
}

// ErrorInfo 单个错误信息
type ErrorInfo struct {
	ID           string      `json:"id"`
	Timestamp    string      `json:"timestamp"`
	Type         string      `json:"type"`
	Message      string      `json:"message"`
	URL          string      `json:"url"`
	UserAgent    string      `json:"userAgent"`
	Filename     string      `json:"filename,omitempty"`
	Lineno       int         `json:"lineno,omitempty"`
	Colno        int         `json:"colno,omitempty"`
	Stack        string      `json:"stack,omitempty"`
	Method       string      `json:"method,omitempty"`
	Status       int         `json:"status,omitempty"`
	StatusText   string      `json:"statusText,omitempty"`
	ErrorCode    int         `json:"errorCode,omitempty"`
	RequestData  interface{} `json:"requestData,omitempty"`
	ResponseData interface{} `json:"responseData,omitempty"`
}

// Metadata 错误报告元数据
type Metadata struct {
	AppName     string `json:"appName"`
	Environment string `json:"environment"`
	Timestamp   string `json:"timestamp"`
	TotalErrors int    `json:"totalErrors"`
}

// ManusNotification Manus通知请求结构
type ManusNotification struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ErrorReportHandler 处理前端错误报告
func ErrorReportHandler(c *gin.Context) {
	var report ErrorReport
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Invalid request body",
		})
		return
	}

	// 记录到日志
	for _, err := range report.Errors {
		fmt.Printf("[ERROR REPORT] Type: %s, Message: %s, URL: %s, Time: %s\n",
			err.Type, err.Message, err.URL, err.Timestamp)
	}

	// 发送到Manus通知系统
	go sendToManus(report)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Error report received",
		"data": gin.H{
			"received": len(report.Errors),
		},
	})
}

// sendToManus 发送错误通知到Manus
func sendToManus(report ErrorReport) {
	// 获取Manus API配置
	apiURL := os.Getenv("BUILT_IN_FORGE_API_URL")
	apiKey := os.Getenv("BUILT_IN_FORGE_API_KEY")

	if apiURL == "" || apiKey == "" {
		fmt.Println("[ERROR REPORT] Manus API not configured, skipping notification")
		return
	}

	// 构建通知内容
	title := fmt.Sprintf("🚨 [%s] 前端错误报告 (%d个错误)",
		report.Metadata.AppName, report.Metadata.TotalErrors)

	var content bytes.Buffer
	content.WriteString(fmt.Sprintf("**环境**: %s\n", report.Metadata.Environment))
	content.WriteString(fmt.Sprintf("**时间**: %s\n\n", report.Metadata.Timestamp))
	content.WriteString("---\n\n")

	for i, err := range report.Errors {
		if i >= 5 { // 最多显示5个错误
			content.WriteString(fmt.Sprintf("\n... 还有 %d 个错误未显示\n", len(report.Errors)-5))
			break
		}
		content.WriteString(fmt.Sprintf("### 错误 %d: %s\n", i+1, err.Type))
		content.WriteString(fmt.Sprintf("- **消息**: %s\n", err.Message))
		content.WriteString(fmt.Sprintf("- **页面**: %s\n", err.URL))
		if err.Filename != "" {
			content.WriteString(fmt.Sprintf("- **文件**: %s:%d:%d\n", err.Filename, err.Lineno, err.Colno))
		}
		if err.Method != "" {
			content.WriteString(fmt.Sprintf("- **API**: %s (状态: %d)\n", err.Method, err.Status))
		}
		content.WriteString("\n")
	}

	// 发送通知
	notification := ManusNotification{
		Title:   title,
		Content: content.String(),
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		fmt.Printf("[ERROR REPORT] Failed to marshal notification: %v\n", err)
		return
	}

	// Manus通知API端点
	notifyURL := fmt.Sprintf("%s/webdevtoken.v1.WebDevService/SendNotification", apiURL)
	req, err := http.NewRequest("POST", notifyURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[ERROR REPORT] Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Connect-Protocol-Version", "1")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR REPORT] Failed to send notification: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[ERROR REPORT] Manus notification failed: %d - %s\n", resp.StatusCode, string(body))
	} else {
		fmt.Printf("[ERROR REPORT] Manus notification sent successfully\n")
	}
}
