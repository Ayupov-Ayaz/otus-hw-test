package sender

import (
	"fmt"
	"go.uber.org/zap"
)

func Start(done <-chan struct{}, consume <-chan []byte, logg *zap.Logger) {
	for {
		select {
		case <-done:
			logg.Info("stop consuming")
			return
		case event := <-consume:
			fmt.Println(string(event))
		}
	}
}
