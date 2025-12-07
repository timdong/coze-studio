/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hertz-contrib/sse"
)

// GinSSESender 是 Gin 版本的 SSE 发送器，实现 sse.SSenderInterface 接口
type GinSSESender struct {
	writer http.ResponseWriter
	ctx    *gin.Context
}

// NewGinSSESender 创建新的 Gin SSE 发送器
func NewGinSSESender(c *gin.Context) *GinSSESender {
	return &GinSSESender{
		writer: c.Writer,
		ctx:    c,
	}
}

// Send 发送 SSE 事件，兼容 hertz-contrib/sse.Event
func (g *GinSSESender) Send(ctx context.Context, event *sse.Event) error {
	if event == nil {
		return nil
	}

	// 写入 event 字段
	if event.Event != "" {
		_, err := fmt.Fprintf(g.writer, "event: %s\n", event.Event)
		if err != nil {
			return err
		}
	}

	// 写入 id 字段
	if event.ID != "" {
		_, err := fmt.Fprintf(g.writer, "id: %s\n", event.ID)
		if err != nil {
			return err
		}
	}

	// 写入 retry 字段
	if event.Retry > 0 {
		_, err := fmt.Fprintf(g.writer, "retry: %d\n", event.Retry)
		if err != nil {
			return err
		}
	}

	// 写入 data 字段
	if len(event.Data) > 0 {
		// SSE 格式要求每行数据前加 "data: "
		// 将数据按行分割
		dataStr := string(event.Data)
		lines := splitSSELines(dataStr)
		for _, line := range lines {
			_, err := fmt.Fprintf(g.writer, "data: %s\n", line)
			if err != nil {
				return err
			}
		}
	}

	// 写入空行表示事件结束
	_, err := fmt.Fprint(g.writer, "\n")
	if err != nil {
		return err
	}

	// 刷新缓冲区
	if flusher, ok := g.writer.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

// splitSSELines 将数据按行分割，处理多行数据
func splitSSELines(data string) []string {
	if data == "" {
		return []string{""}
	}

	var lines []string
	var current []rune

	for _, r := range data {
		if r == '\n' {
			lines = append(lines, string(current))
			current = nil
		} else {
			current = append(current, r)
		}
	}

	// 添加最后一行（如果有）
	if len(current) > 0 {
		lines = append(lines, string(current))
	} else if len(lines) == 0 {
		// 如果没有换行符，整个字符串作为一行
		lines = append(lines, data)
	}

	return lines
}

