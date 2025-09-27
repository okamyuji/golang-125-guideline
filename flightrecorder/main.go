package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/trace"
	"sync"
	"syscall"
	"time"
)

// より実用的なスナップショット管理
type SnapshotManager struct {
	cooldown time.Duration
	lastSnap time.Time
	mu       sync.Mutex
	fr       *trace.FlightRecorder
}

func NewSnapshotManager(fr *trace.FlightRecorder, cooldown time.Duration) *SnapshotManager {
	return &SnapshotManager{
		cooldown: cooldown,
		fr:       fr,
	}
}

func (sm *SnapshotManager) TrySnapshot() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if time.Since(sm.lastSnap) < sm.cooldown {
		return false // クールダウン中
	}

	f, err := os.CreateTemp("", "snap-*.ctrace")
	if err != nil {
		log.Printf("temp create: %v", err)
		return false
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("file close: %v", err)
		}
	}()

	if _, err := sm.fr.WriteTo(f); err != nil {
		log.Printf("trace write: %v", err)
		return false
	}

	sm.lastSnap = time.Now()
	log.Printf("trace snapshot: %s", f.Name())
	return true
}

func main() {
	// 直前5s保持、最大16MiB目安（ヒント）
	fr := trace.NewFlightRecorder(trace.FlightRecorderConfig{
		MinAge:   5 * time.Second,
		MaxBytes: 16 << 20,
	})
	if err := fr.Start(); err != nil {
		log.Fatalf("flight recorder start: %v", err)
	}
	defer fr.Stop()

	// クールダウン期間を10秒に設定（連続スナップショットを防ぐ）
	sm := NewSnapshotManager(fr, 10*time.Second)

	// 手動トリガ（SIGUSR1）、終了時（SIGTERM）でスナップショット
	// SIGINT（Ctrl-C）は通常終了として扱う
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGUSR1, syscall.SIGTERM)
	go func() {
		for s := range sigc {
			log.Printf("signal %v -> attempting snapshot", s)
			if sm.TrySnapshot() {
				log.Printf("snapshot completed")
			} else {
				log.Printf("snapshot skipped (cooldown)")
			}

			// SIGTERMの場合は終了
			if s == syscall.SIGTERM {
				log.Println("received SIGTERM, shutting down...")
				os.Exit(0)
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// 疑似遅延
		time.Sleep(200 * time.Millisecond)
		if _, err := fmt.Fprintln(w, "hello"); err != nil {
			log.Printf("response write: %v", err)
		}
		if d := time.Since(start); d > 100*time.Millisecond {
			go func() {
				if sm.TrySnapshot() {
					log.Printf("auto snapshot triggered by slow request (%v)", d)
				}
			}()
		}
	})

	log.Println("listen :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
