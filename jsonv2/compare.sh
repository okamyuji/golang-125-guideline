#!/usr/bin/env bash

echo "=========================================="
echo "JSON v1 vs v2 挙動比較スクリプト"
echo "=========================================="

echo
echo "=== JSON v1 (標準) での実行 ==="
go run .

echo
echo "=========================================="
echo "=== JSON v2 (実験的) での実行 ==="
echo "注意: Go 1.25以降が必要です"
echo "=========================================="

if GOEXPERIMENT=jsonv2 go run . 2>/dev/null; then
    echo "✅ jsonv2実行成功"
else
    echo "❌ jsonv2実行失敗 - Go 1.25以降またはjsonv2サポートが必要"
    echo
    echo "【重要な違い】"
    echo "v1ではtime.Timeのゼロ値が'0001-01-01T00:00:00Z'として出力される"
    echo "v2ではomitemptyでゼロ値が省略される可能性がある"
    echo "v2ではomitzeroタグでより厳密なゼロ値制御が可能"
    echo
    echo "【互換性への影響】"
    echo "- APIレスポンスの形式が変わる可能性"
    echo "- クライアント側でフィールドの存在チェックが必要"
    echo "- データベースとの同期で問題が発生する可能性"
fi

echo
echo "=========================================="
echo "【対策】"
echo "1. omitzeroタグを明示的に使用する"
echo "2. time.Timeのゼロ値チェックを実装する"
echo "3. 既存APIの互換性テストを実施する"
echo "4. 段階的移行計画の策定をする"
echo "=========================================="
