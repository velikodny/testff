package main

import "strconv"

type Validator struct {
	err *Error
}

func (v *Validator) validCheckType(value string) bool {
	switch value {
	case "DEVICE", "BIOMETRIC", "COMBO":
		return true
	}

	v.err = &Error{
		Code:    0,
		Message: "fault to validate a checkType",
	}

	return false
}

func (v *Validator) validActivityType(value string) bool {
	switch value {
	case "SIGNUP", "LOGIN", "PAYMENT", "CONFIRMATION":
		return true
	}

	v.err = &Error{
		Code:    0,
		Message: "fault to validate a activityType",
	}

	return false
}

func (v *Validator) validActivityData(data []ActivityData) bool {
	valid := true
	errStr := &Error{
		Code:    0,
		Message: "the specified value of ActivityData not a correct",
	}
	keyStore := make(map[string]bool, len(data))

	for _, item := range data {
		if ok := keyStore[item.KvpKey]; ok {
			v.err = &Error{
				Message: "a duplicate of an activity data key",
			}
			valid = false
			break
		}
		keyStore[item.KvpKey] = true
		if item.KvpType == "general.string" {
			continue
		}
		switch item.KvpType {
		case "general.integer":
			_, err := strconv.ParseInt(item.KvpValue, 10, 64)
			if err != nil {
				v.err = errStr
				valid = false
				break
			}
		case "general.float":
			_, err := strconv.ParseFloat(item.KvpValue, 64)
			if err != nil {
				v.err = errStr
				valid = false
				break
			}
		case "general.bool":
			_, err := strconv.ParseBool(item.KvpValue)
			if err != nil {
				v.err = errStr
				valid = false
				break
			}
		default:
			v.err = &Error{
				Message: "specified an invalid type of an activity data key",
			}
			valid = false
		}
	}

	return valid
}
