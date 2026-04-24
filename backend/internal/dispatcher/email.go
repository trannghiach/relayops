package dispatcher

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
)

type EmailDispatcher struct {
	RandomFailure bool
	FailureRate   float64
}

type EmailPayload struct {
	Template string `json:"template"`
}

func NewEmailDispatcher(randomFailure bool, failureRate float64) *EmailDispatcher {
	return &EmailDispatcher{
		RandomFailure: randomFailure,
		FailureRate:   failureRate,
	}
}

func (d *EmailDispatcher) SendEmailMock(payload []byte) error {
	var data EmailPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("failed to unmarshal email payload: %w", err)
	}

	log.Printf("sending email (mock) with template: %s", data.Template)

	if d.RandomFailure && rand.Float64() < d.FailureRate {
		return fmt.Errorf("simulated email sending failure")
	}

	return nil
}
