package service

import (
	"Cart-API/app/internal/errs"
	"Cart-API/app/internal/models"
	"Cart-API/app/internal/repository/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"testing"
)

func TestCreateCart(t *testing.T) {
	repoError := errors.New("some error")

	tests := []struct {
		name           string
		mockCartID     string
		mockError      error
		expectedCartID string
		expectedError  error
	}{
		{
			name:           "Success Test",
			mockCartID:     "1",
			mockError:      nil,
			expectedCartID: "1",
			expectedError:  nil,
		},
		{
			name:           "Failure: repository error",
			mockCartID:     "",
			mockError:      repoError,
			expectedCartID: "",
			expectedError:  repoError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockCartRepository(ctrl)
			mockRepo.EXPECT().CreateCart(context.Background()).Return(test.mockCartID,
				test.mockError)

			service := NewCart(mockRepo)
			cartID, err := service.CreateCart(context.Background())

			assert.Equal(t, test.expectedCartID, cartID)
			assert.ErrorIs(t, err, test.expectedError)
		})
	}
}

func TestGetCartByID(t *testing.T) {
	repoError := errors.New("some error")

	tests := []struct {
		name                string
		inputCartID         string
		mockCartIsAvailable bool
		mockCart            models.Cart
		mockError           error
		expectedCart        models.Cart
		expectedError       error
	}{
		{
			name:                "Success Test: empty cart",
			inputCartID:         "1",
			mockCartIsAvailable: true,
			mockCart:            models.Cart{ID: "1", Items: []models.CartItem{}},
			mockError:           nil,
			expectedCart:        models.Cart{ID: "1", Items: []models.CartItem{}},
			expectedError:       nil,
		},
		{
			name:                "Success Test: cart with data",
			inputCartID:         "1",
			mockCartIsAvailable: true,
			mockCart: models.Cart{ID: "1", Items: []models.CartItem{
				{ID: "1", CartID: "1", Product: "Apples", Quantity: 2},
				{ID: "2", CartID: "1", Product: "Bananas", Quantity: 10},
			}},
			mockError: nil,
			expectedCart: models.Cart{ID: "1", Items: []models.CartItem{
				{ID: "1", CartID: "1", Product: "Apples", Quantity: 2},
				{ID: "2", CartID: "1", Product: "Bananas", Quantity: 10},
			}},
			expectedError: nil,
		},
		{
			name:                "Failure: cart not found",
			inputCartID:         "10",
			mockCartIsAvailable: false,
			mockCart:            models.Cart{},
			mockError:           nil,
			expectedCart:        models.Cart{},
			expectedError:       errs.ErrCartNotFound,
		},
		{
			name:                "Failure: repository error",
			inputCartID:         "1",
			mockCartIsAvailable: false,
			mockCart:            models.Cart{},
			mockError:           repoError,
			expectedCart:        models.Cart{},
			expectedError:       repoError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := mocks.NewMockCartRepository(ctrl)
			mockRepo.EXPECT().CartIsAvailable(context.Background(), test.inputCartID).
				Return(test.mockCartIsAvailable, test.mockError)

			if test.mockCartIsAvailable {
				mockRepo.EXPECT().GetCartByID(context.Background(), test.inputCartID).
					Return(test.mockCart, test.mockError)
			}

			service := NewCart(mockRepo)
			cart, err := service.GetCartByID(context.Background(), test.inputCartID)

			assert.Equal(t, test.expectedCart, cart)
			assert.ErrorIs(t, err, test.expectedError)
		})
	}
}
