package service

import (
	"Cart-API/app/internal/errs"
	"Cart-API/app/internal/models"
	"Cart-API/app/internal/repository/mocks"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddItemToCart(t *testing.T) {
	repoError := errors.New("some error")

	tests := []struct {
		name                string
		mockCartIsAvailable bool
		mockItemID          string
		mockCartID          string
		mockProduct         string
		mockQuantity        int
		mockError           error
		expectedCartItem    models.CartItem
		expectedError       error
	}{
		{
			name:                "Success Test",
			mockCartIsAvailable: true,
			mockItemID:          "1",
			mockCartID:          "1",
			mockProduct:         "Apple",
			mockQuantity:        10,
			mockError:           nil,
			expectedCartItem: models.CartItem{
				ID:       "1",
				CartID:   "1",
				Product:  "Apple",
				Quantity: 10,
			},
			expectedError: nil,
		},
		{
			name:                "Failure: cart not found",
			mockCartIsAvailable: false,
			mockItemID:          "",
			mockCartID:          "15",
			mockProduct:         "Apple",
			mockQuantity:        10,
			mockError:           nil,
			expectedCartItem:    models.CartItem{},
			expectedError:       errs.ErrCartNotFound,
		},
		{
			name:                "Failure: repository error",
			mockCartIsAvailable: false,
			mockItemID:          "",
			mockCartID:          "15",
			mockProduct:         "Apple",
			mockQuantity:        10,
			mockError:           repoError,
			expectedCartItem:    models.CartItem{},
			expectedError:       repoError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockCartRepo := mocks.NewMockCartRepository(ctrl)
			mockCartRepo.EXPECT().CartIsAvailable(context.Background(),
				test.mockCartID).Return(test.mockCartIsAvailable, test.mockError)

			mockCartItemRepo := mocks.NewMockCartItemRepository(ctrl)

			if test.mockCartIsAvailable {
				mockCartItemRepo.EXPECT().AddItemToCart(context.Background(),
					test.mockCartID, test.mockProduct, test.mockQuantity).Return(test.mockItemID,
					test.mockError)
				mockCartItemRepo.EXPECT().GetItemByID(context.Background(),
					test.mockCartID).Return(test.expectedCartItem, test.mockError)
			}

			service := NewCartItem(mockCartRepo, mockCartItemRepo)
			item, err := service.AddItemToCart(context.Background(),
				test.mockCartID, test.mockProduct, test.mockQuantity)

			assert.Equal(t, test.expectedCartItem, item)
			assert.ErrorIs(t, err, test.expectedError)
		})
	}
}

func TestRemoveItemFromCart(t *testing.T) {
	repoError := errors.New("some error")

	tests := []struct {
		name                string
		mockCartIsAvailable bool
		mockItemIsAvailable bool
		mockCartID          string
		mockItemID          string
		mockError           error
		expectedError       error
	}{
		{
			name:                "Success Test",
			mockCartIsAvailable: true,
			mockItemIsAvailable: true,
			mockCartID:          "1",
			mockItemID:          "1",
			mockError:           nil,
			expectedError:       nil,
		},
		{
			name:                "Failure: cart not found",
			mockCartIsAvailable: false,
			mockItemIsAvailable: false,
			mockCartID:          "1",
			mockItemID:          "15",
			mockError:           nil,
			expectedError:       errs.ErrCartNotFound,
		},
		{
			name:                "Failure: item not found",
			mockCartIsAvailable: true,
			mockItemIsAvailable: false,
			mockCartID:          "1",
			mockItemID:          "15",
			mockError:           nil,
			expectedError:       errs.ErrItemNotFound,
		},
		{
			name:                "Failure: repository error",
			mockCartIsAvailable: false,
			mockItemIsAvailable: false,
			mockCartID:          "15",
			mockItemID:          "15",
			mockError:           repoError,
			expectedError:       repoError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockCartRepo := mocks.NewMockCartRepository(ctrl)
			mockCartRepo.EXPECT().CartIsAvailable(context.Background(),
				test.mockCartID).Return(test.mockCartIsAvailable, test.mockError)

			mockCartItemRepo := mocks.NewMockCartItemRepository(ctrl)
			if test.mockCartIsAvailable {
				mockCartItemRepo.EXPECT().ItemIsAvailable(context.Background(),
					test.mockItemID).Return(test.mockItemIsAvailable, test.mockError)
			}

			if test.mockItemIsAvailable {
				mockCartItemRepo.EXPECT().RemoveItemFromCart(context.Background(),
					test.mockCartID, test.mockItemID).Return(test.mockError)
			}

			service := NewCartItem(mockCartRepo, mockCartItemRepo)
			err := service.RemoveItemFromCart(context.Background(),
				test.mockCartID, test.mockItemID)

			assert.ErrorIs(t, err, test.expectedError)
		})
	}
}
