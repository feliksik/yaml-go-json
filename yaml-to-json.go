package main


// transformData courtesy https://github.com/bronze1man/yaml2json/
//
import (
	"errors"
	"fmt"
	"strconv"
)

func transformData(in interface{}) (out interface{}, err error) {
	switch in.(type) {
	case map[string]interface{}: // added by Eric
		o := make(map[string]interface{})
		for k, v := range in.(map[string]interface{}) {
			v, err = transformData(v)
			if err != nil {
				return nil, err
			}
			o[k] = v
		}
		return o, nil
	case map[interface{}]interface{}:
		o := make(map[string]interface{})
		for k, v := range in.(map[interface{}]interface{}) {
			sk := ""
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				return nil, errors.New(
					fmt.Sprintf("type not match: expect map key string or int get: %T", k))
			}
			v, err = transformData(v)
			if err != nil {
				return nil, err
			}
			o[sk] = v
		}
		return o, nil
	case []interface{}:
		in1 := in.([]interface{})
		len1 := len(in1)
		o := make([]interface{}, len1)
		for i := 0; i < len1; i++ {
			o[i], err = transformData(in1[i])
			if err != nil {
				return nil, err
			}
		}
		return o, nil
	default:
		return in, nil
	}
	return in, nil
}
