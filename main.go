package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/slok/goresilience"
	"github.com/slok/goresilience/bulkhead"
	"github.com/slok/goresilience/circuitbreaker"
	"github.com/slok/goresilience/retry"
	"github.com/slok/goresilience/timeout"
)

func main() {
	// Create our execution chain.
	cmd := goresilience.RunnerChain(
		bulkhead.NewMiddleware(bulkhead.Config{}),
		retry.NewMiddleware(retry.Config{}),
		timeout.NewMiddleware(timeout.Config{}),
		circuitbreaker.NewMiddleware(circuitbreaker.Config{
			//ErrorPercentThresholdToOpen:        50,
			//MinimumRequestToOpen:               20,
			//SuccessfulRequiredOnHalfOpen:       1,
			//WaitDurationInOpenState:            5 * time.Second,
			//MetricsSlidingWindowBucketQuantity: 10,
			//MetricsBucketDuration:              1 * time.Second,
		}),
	)

	result := ""
	err := cmd.Run(context.Background(), func(_ context.Context) error {

		resp, err := http.Get("http://localhost:1080/test")
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 500 {
			return errors.New("received status code 500 from server")
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		result = "all ok"
		return nil
	})

	if err != nil {
		result = "not ok, but fallback"
	}

	fmt.Printf("result: %s", result)
}
