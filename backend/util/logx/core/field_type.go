package core

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type FieldType uint32

const (
	BoolType   FieldType = 1
	StringType FieldType = 2
	Int32Type  FieldType = 3
	Int64Type  FieldType = 4
	FloatType  FieldType = 5
	DoubleType FieldType = 6
)

func tryBool(val interface{}) (bool, error) {
	switch x := val.(type) {
	case string:
		{
			temp := strings.ToLower(val.(string))
			switch temp {
			case "false":
				return false, nil
			case "true":
				return true, nil
			default:
				return false, errors.New("Invalid Type")
			}
		}
	case bool:
		{
			return x, nil
		}
	case uint:
		{
			if x == uint(1) {
				return true, nil
			}
			if x == uint(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	case int:
		{
			if x == int(1) {
				return true, nil
			}
			if x == int(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	case uint8:
		{
			if x == uint8(1) {
				return true, nil
			}
			if x == uint8(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	case int8:
		{
			if x == int8(1) {
				return true, nil
			}
			if x == int8(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	case uint16:
		{
			if x == uint16(1) {
				return true, nil
			}
			if x == uint16(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	case int16:
		{
			if x == int16(1) {
				return true, nil
			}
			if x == int16(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	case uint32:
		{
			if x == uint32(1) {
				return true, nil
			}
			if x == uint32(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	case int32:
		{
			if x == int32(1) {
				return true, nil
			}
			if x == int32(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	case uint64:
		{
			if x == uint64(1) {
				return true, nil
			}
			if x == uint64(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	case int64:
		{
			if x == int64(1) {
				return true, nil
			}
			if x == int64(0) {
				return false, nil
			}
			return false, errors.New("Invalid Type")
		}
	default:
		{
			return false, errors.New("Invalid Type")
		}
	}
}

func tryString(val interface{}) (string, error) {
	return fmt.Sprintf("%v", val), nil
}

func tryFloat(val interface{}) (interface{}, error) {
	switch x := val.(type) {
	case float32:
		return x, nil
	case float64:
		return float32(x), nil
	case int8:
		return float32(x), nil
	case uint8:
		return float32(x), nil
	case int16:
		return float32(x), nil
	case uint16:
		return float32(x), nil
	case int32:
		return float32(x), nil
	case uint32:
		return float32(x), nil
	case int64:
		return float32(x), nil
	case uint64:
		return float32(x), nil
	case int:
		return float32(x), nil
	case uint:
		return float32(x), nil
	case string:
		{
			res, err := strconv.ParseFloat(x, 32)
			if err != nil {
				return nil, err
			}
			return float32(res), nil
		}
	default:
		{
			return nil, errors.New("Invlaid Type")
		}
	}
}

func tryInt32(val interface{}) (interface{}, error) {
	switch x := val.(type) {
	case int8, uint8, int16, uint16, int32, uint32:
		return x, nil
	case string:
		{
			res, err := strconv.ParseInt(x, 10, 32)
			if err != nil {
				return nil, err
			}
			return int32(res), nil
		}
	default:
		{
			return nil, errors.New("Invlaid Type")
		}
	}
}

func tryInt64(val interface{}) (interface{}, error) {
	switch x := val.(type) {
	case int8, uint8, int16, uint16, int32, uint32, int64, uint64, int, uint:
		return x, nil
	case string:
		{
			res, err := strconv.ParseInt(x, 10, 64)
			if err != nil {
				return nil, err
			}
			return int64(res), nil
		}
	default:
		{
			return nil, errors.New("Invlaid Type")
		}
	}
}

func tryDouble(val interface{}) (interface{}, error) {
	switch x := val.(type) {
	case float32:
		return float64(x), nil
	case float64:
		return x, nil
	case int8:
		return float64(x), nil
	case uint8:
		return float64(x), nil
	case int16:
		return float64(x), nil
	case uint16:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case uint32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case uint64:
		return float64(x), nil
	case int:
		return float64(x), nil
	case uint:
		return float64(x), nil
	case string:
		{
			res, err := strconv.ParseFloat(x, 64)
			if err != nil {
				return nil, err
			}
			return float64(res), nil
		}
	default:
		{
			return nil, errors.New("Invlaid Type")
		}
	}
}
