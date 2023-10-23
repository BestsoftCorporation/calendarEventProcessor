package usecase

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func startRecuringEvent() {
	ticker := time.NewTicker(time.Minute)
	done := make(chan bool)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				h, m, _ := time.Now().Clock()
				if m == 0 && (h == 00) {
					fmt.Printf("Doing the job")
				}
			}
		}
	}()

	<-ctx.Done()
	stop()
	done <- true
}
