package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// omitemptyの挙動変更を示すための構造体
type Event struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitzero"` // jsonv2の新機能
	Count       int       `json:"count,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Metadata    *Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Version string `json:"version,omitempty"`
	Author  string `json:"author,omitempty"`
}

func main() {
	fmt.Println("=== 基本的なJSON操作 ===")
	in := []byte(`{"id":1,"name":"gopher"}`)
	var u User
	if err := json.Unmarshal(in, &u); err != nil {
		panic(err)
	}
	out, _ := json.Marshal(u)
	fmt.Println("User:", string(out))

	fmt.Println("\n=== omitempty vs omitzero の挙動比較 ===")

	// ケース1: 完全にゼロ値の構造体
	event1 := Event{Name: "Empty Event"}
	data1, _ := json.Marshal(event1)
	fmt.Println("ゼロ値フィールド含む:", string(data1))

	// ケース2: 一部フィールドに値がある構造体
	event2 := Event{
		Name:        "Partial Event",
		Description: "説明あり",
		Count:       0,          // ゼロ値だがomitemptyで省略される
		Tags:        []string{}, // 空スライスでomitemptyで省略される
	}
	data2, _ := json.Marshal(event2)
	fmt.Println("一部値あり:", string(data2))

	// ケース3: time.Timeのゼロ値問題を示す
	event3 := Event{
		Name:      "Time Zero Value Event",
		StartTime: time.Time{}, // ゼロ値 - v1では省略されない、v2では省略される可能性
		EndTime:   time.Time{}, // ゼロ値 - omitzeroで確実に省略
	}
	data3, _ := json.Marshal(event3)
	fmt.Println("time.Timeゼロ値:", string(data3))

	// ケース4: 実際の時刻値
	now := time.Now()
	event4 := Event{
		Name:      "Actual Time Event",
		StartTime: now,
		EndTime:   now.Add(time.Hour),
		Count:     5,
		Tags:      []string{"important", "scheduled"},
		Metadata:  &Metadata{Version: "1.0", Author: "gopher"},
	}
	data4, _ := json.Marshal(event4)
	fmt.Println("実際の値:", string(data4))

	fmt.Println("\n=== 互換性問題の例 ===")
	fmt.Println("【重要】time.Timeのゼロ値の扱い:")
	fmt.Println("- v1: omitemptyでもゼロ値が出力される (0001-01-01T00:00:00Z)")
	fmt.Println("- v2: omitemptyでゼロ値が省略される可能性")
	fmt.Println("- 対策: omitzeroタグを明示的に使用してゼロ値制御")

	fmt.Println("\n【実際の問題シナリオ】")
	fmt.Println("1. APIレスポンスでtime.Timeフィールドが突然消える")
	fmt.Println("2. クライアントアプリがフィールド存在を前提とした処理でエラー")
	fmt.Println("3. データベース同期でNULL vs ゼロ値の判定が困難")
	fmt.Println("4. 既存のJSONスキーマ検証が失敗する可能性")

	fmt.Println("\n構造体ポインタのnil vs ゼロ値:")
	eventWithNil := Event{
		Name:     "Nil Metadata Event",
		Metadata: nil, // nilポインタ - 常に省略
	}
	dataWithNil, _ := json.Marshal(eventWithNil)
	fmt.Println("nilポインタ:", string(dataWithNil))

	eventWithEmpty := Event{
		Name:     "Empty Metadata Event",
		Metadata: &Metadata{}, // 空の構造体 - v2では省略される可能性
	}
	dataWithEmpty, _ := json.Marshal(eventWithEmpty)
	fmt.Println("空構造体:", string(dataWithEmpty))

	fmt.Println("\n=== 移行対策コード例 ===")

	// 安全なtime.Time処理の例
	type SafeEvent struct {
		Name      string     `json:"name"`
		StartTime *time.Time `json:"start_time,omitempty"` // ポインタ使用でnilチェック可能
		EndTime   *time.Time `json:"end_time,omitzero"`    // omitzeroで明示的制御
	}

	safeEvent := SafeEvent{Name: "Safe Event"}
	safeData, _ := json.Marshal(safeEvent)
	fmt.Println("安全な実装:", string(safeData))

	now2 := time.Now()
	safeEventWithTime := SafeEvent{
		Name:      "Safe Event With Time",
		StartTime: &now2,
		EndTime:   &now2,
	}
	safeDataWithTime, _ := json.Marshal(safeEventWithTime)
	fmt.Println("時刻あり:", string(safeDataWithTime))

	// 実験ON: GOEXPERIMENT=jsonv2 でビルド/実行
	// 例) Linux: GOEXPERIMENT=jsonv2 go run .
	// 参考: v2はAPI/実装が進化中（将来変更の可能性あり）
	fmt.Println("\n実験フラグ: GOEXPERIMENT=jsonv2 go run .")
	_, _ = os.Stdout.Write([]byte("ok\n"))
}
