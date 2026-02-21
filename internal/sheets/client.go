package sheets

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func NewClient(credentialsPath string) (*sheets.Service, error) {
	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("reading credentials file: %w", err)
	}

	srv, err := sheets.NewService(context.Background(), option.WithCredentialsJSON(data))
	if err != nil {
		return nil, fmt.Errorf("creating sheets service: %w", err)
	}

	return srv, nil
}
