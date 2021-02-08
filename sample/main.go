package sample

import (
	"fmt"
	"go_experiment/metric"
	"math/rand"
	"time"
)

func main() {
	err := metric.SetGlobal("127.0.0.1", "8125", metric.SetEnv("staging"), metric.ServiceName("biller_testing"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ch := make(chan int)
	n := 10000
	for i := 0; i < n; i++ {
		go func() {
			tracer, _ := metric.NewTrace("biller.histogram")
			defer func() {
				err := tracer.Stop(nil)
				if err != nil {
					fmt.Println(err.Error())
				}
			}()
			tracer.AddTags(metric.Tags{
				"url":         "http://localhost",
				"http_code":   200,
				"biller_name": "indosat",
			})
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			ch <- 1
		}()
	}
	for i := 0; i < n; i++ {
		<-ch
	}

	metric.Close()
}
