package retry

import (
	"context"
	"go.uber.org/zap"
	"log"
	"strconv"
	"test/pkg/logger"
	"time"
)

// 实现定时任务的通用框架
// RetryConfig 重试配置
type RetryConfig struct {
	Interval     time.Duration // 固定重试间隔（固定间隔模式用）
	BaseInterval time.Duration // 基础间隔（指数退避模式用）
	MaxInterval  time.Duration // 最大间隔（指数退避模式用）
	MaxRetries   int           // 最大重试次数（0表示无限重试）
	RetryMode    string        // 重试模式：fixed（固定间隔）/exponential（指数退避）
}

func RetryTask(ctx context.Context, config RetryConfig,
	task func() bool, taskName string) {
	retryCount := 0
	interval := config.Interval
	if config.RetryMode == "exponential" {
		interval = config.BaseInterval
	}
	logger.Info("启动定时重试任务", zap.String("taskName", taskName))
	for {
		select {
		case <-ctx.Done():
			//上下文取消
			logger.Info("上下文取消")
			return
		default:
			success := task()
			if success {
				logger.Info("任务重试成功", zap.String("taskname", taskName), zap.String("重试次数", strconv.Itoa(retryCount)))
				return
			}
		}
		retryCount++
		// 检查是否达到最大重试次数
		if config.MaxRetries > 0 && retryCount >= config.MaxRetries {
			log.Printf("任务 %s 重试达到最大次数（%d次），停止重试", taskName, config.MaxRetries)
			return
		}
		time.Sleep(interval)
		//TODO: 指数退避模式：间隔翻倍，不超过最大值
	}

}
