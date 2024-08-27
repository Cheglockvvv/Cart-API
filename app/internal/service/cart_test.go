package service

import (
	"Cart-API/app/internal/models"
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateCart(t *testing.T) {
	tests := []struct {
		name          string
		mockCart      *models.Cart
		mockError     error
		expectedCart  *models.Cart
		expectedError error
	}{
		{
			name:          "Success Test",
			mockCart:      &models.Cart{ID: "1", Items: []models.CartItem{}},
			mockError:     nil,
			expectedCart:  &models.Cart{ID: "1", Items: []models.CartItem{}},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
		})
	}
}
