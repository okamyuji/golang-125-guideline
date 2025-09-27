package main

import (
	"context"
	"fmt"
	"time"
)

// doWithDeadline は指定された遅延時間後に完了するか、コンテキストがキャンセルされるまで待機する
func doWithDeadline(ctx context.Context, delay time.Duration) error {
	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	fmt.Println("実行開始...")
	err := doWithDeadline(ctx, 200*time.Millisecond)

	if err != nil {
		fmt.Printf("エラー: %v\n", err)
	} else {
		fmt.Println("正常完了")
	}
}
