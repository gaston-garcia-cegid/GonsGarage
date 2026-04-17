package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCar_Validate_Success(t *testing.T) {
	// Arrange
	car := &Car{
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC-123",
		Color:        "Blue",
		Mileage:      5000,
		OwnerID:      uuid.New(),
	}

	// Act
	err := car.Validate()

	// Assert
	assert.NoError(t, err)
}

func TestCar_Validate_RequiredFields(t *testing.T) {
	testCases := []struct {
		name        string
		car         Car
		expectError bool
	}{
		{
			name: "empty make",
			car: Car{
				Model:        "Camry",
				Year:         2023,
				LicensePlate: "ABC-123",
				Color:        "Blue",
			},
			expectError: true,
		},
		{
			name: "empty model",
			car: Car{
				Make:         "Toyota",
				Year:         2023,
				LicensePlate: "ABC-123",
				Color:        "Blue",
			},
			expectError: true,
		},
		{
			name: "invalid year too old",
			car: Car{
				Make:         "Toyota",
				Model:        "Camry",
				Year:         1800,
				LicensePlate: "ABC-123",
				Color:        "Blue",
			},
			expectError: true,
		},
		{
			name: "invalid year too future",
			car: Car{
				Make:         "Toyota",
				Model:        "Camry",
				Year:         2050,
				LicensePlate: "ABC-123",
				Color:        "Blue",
			},
			expectError: true,
		},
		{
			name: "empty license plate",
			car: Car{
				Make:  "Toyota",
				Model: "Camry",
				Year:  2023,
				Color: "Blue",
			},
			expectError: true,
		},
		{
			name: "empty color",
			car: Car{
				Make:         "Toyota",
				Model:        "Camry",
				Year:         2023,
				LicensePlate: "ABC-123",
			},
			expectError: true,
		},
		{
			name: "negative mileage",
			car: Car{
				Make:         "Toyota",
				Model:        "Camry",
				Year:         2023,
				LicensePlate: "ABC-123",
				Color:        "Blue",
				Mileage:      -100,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := tc.car.Validate()

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidCarData, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCar_IsOwnedBy(t *testing.T) {
	// Arrange
	ownerID := uuid.New()
	differentOwnerID := uuid.New()
	car := &Car{
		OwnerID: ownerID,
	}

	testCases := []struct {
		name     string
		userID   uuid.UUID
		expected bool
	}{
		{
			name:     "owned by user",
			userID:   ownerID,
			expected: true,
		},
		{
			name:     "not owned by user",
			userID:   differentOwnerID,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := car.IsOwnedBy(tc.userID)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCar_TableName(t *testing.T) {
	// Arrange
	car := Car{}

	// Act
	tableName := car.TableName()

	// Assert
	assert.Equal(t, "cars", tableName)
}
