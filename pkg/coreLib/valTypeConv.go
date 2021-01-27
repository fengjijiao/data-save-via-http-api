package coreLib

import (
	"errors"
	"strconv"
)

//值类型转换
func ValTypeConv(valStr string, valType int) (interface{}, error) {
	if valType == TYPE_STRING {
		return valStr, nil
	}else if valType == TYPE_INT {
		res, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return -1, err
		}
		return res, nil
	}else if valType == TYPE_FLOAT {
		res, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return -1, err
		}
		return res, nil 
	}else if valType == TYPE_BOOl {
		res, err := strconv.ParseBool(valStr)
		if err != nil {
			return false, err
		}
		return res, nil 
	}else {
		return nil, errors.New("this value type was no supported.")
	}
}