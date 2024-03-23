package service

import (
	"context"
	"fmt"
)

// formatJSONService struct
type formatJSONService struct {
	ctx context.Context
}

// NewFormatJSONService creates a new formatJSONService application struct
func NewFormatJSONService() *formatJSONService {
	return &formatJSONService{}
}

// Startup is called when the service starts. The context is saved
// so we can call the runtime methods
func (a *formatJSONService) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *formatJSONService) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
