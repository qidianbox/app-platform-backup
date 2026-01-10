package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DB 是gorm.DB的别名，方便外部包使用
type DB = gorm.DB

// TransactionFunc 事务函数类型
type TransactionFunc func(tx *gorm.DB) error

// WithTransaction 执行事务
// 如果fn返回错误，事务将回滚；否则提交
func WithTransaction(fn TransactionFunc) error {
	return WithTransactionContext(context.Background(), fn)
}

// WithTransactionContext 带上下文的事务执行
func WithTransactionContext(ctx context.Context, fn TransactionFunc) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // 重新抛出panic
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback failed: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithTransactionRetry 带重试的事务执行
func WithTransactionRetry(fn TransactionFunc, maxRetries int) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := WithTransaction(fn); err != nil {
			lastErr = err
			// 简单的退避策略
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}
		return nil
	}
	return fmt.Errorf("transaction failed after %d retries: %w", maxRetries, lastErr)
}

// BatchInsert 批量插入（带事务）
func BatchInsert[T any](items []T, batchSize int) error {
	if len(items) == 0 {
		return nil
	}

	return WithTransaction(func(tx *gorm.DB) error {
		for i := 0; i < len(items); i += batchSize {
			end := i + batchSize
			if end > len(items) {
				end = len(items)
			}
			if err := tx.Create(items[i:end]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchUpdate 批量更新（带事务）
func BatchUpdate[T any](items []T, updateFn func(tx *gorm.DB, item T) error) error {
	if len(items) == 0 {
		return nil
	}

	return WithTransaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := updateFn(tx, item); err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchDelete 批量删除（带事务）
func BatchDelete[T any](ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return WithTransaction(func(tx *gorm.DB) error {
		var model T
		return tx.Delete(&model, ids).Error
	})
}
