package dispatcher

import (
	"encoding/json"
	"fmt"
	"log"
)

type EmailPayload struct {
	Template string `json:"template"`
}

func SendEmailMock(payload []byte) error {
	var data EmailPayload
	if err := json.Unmarshal(payload, &data); err != nil {
		return fmt.Errorf("failed to unmarshal email payload: %w", err)
	}

	log.Printf("sending email (mock) with template: %s", data.Template)

	// uncomment to simulate failure
	// return fmt.Errorf("simulated email sending failure")

	return nil
}
