package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidCheckType(t *testing.T) {
	validator := Validator{}

	assert.Equal(t, true, validator.validCheckType("BIOMETRIC"))
	assert.Equal(t, false, validator.validCheckType("SIG"))

}

func TestValidActivityType(t *testing.T) {
	validator := Validator{}

	assert.Equal(t, true, validator.validActivityType("LOGIN"))
	assert.Equal(t, false, validator.validActivityType("LOGIN_2"))
}

func TestValidActivityData(t *testing.T) {
	validator := Validator{}
	correctData := []ActivityData{
		ActivityData{
			KvpKey:   "Test 1",
			KvpValue: "1",
			KvpType:  "general.integer",
		},
		ActivityData{
			KvpKey:   "Test 2",
			KvpValue: "true",
			KvpType:  "general.bool",
		},
		ActivityData{
			KvpKey:   "Test 3",
			KvpValue: "2.33",
			KvpType:  "general.float",
		},
		ActivityData{
			KvpKey:   "Test 4",
			KvpValue: "2",
			KvpType:  "general.string",
		},
	}

	assert.Equal(t, true, validator.validActivityData(correctData))

	incorrectDataKey := []ActivityData{
		ActivityData{
			KvpKey:   "Test",
			KvpValue: "true",
			KvpType:  "general.bool",
		},
		ActivityData{
			KvpKey:   "Test",
			KvpValue: "true",
			KvpType:  "general.bool",
		},
	}

	assert.Equal(t, false, validator.validActivityData(incorrectDataKey))

	incorrectDataValue1 := []ActivityData{
		ActivityData{
			KvpKey:   "Test",
			KvpValue: "food",
			KvpType:  "general.bool",
		},
	}

	assert.Equal(t, false, validator.validActivityData(incorrectDataValue1))
}
