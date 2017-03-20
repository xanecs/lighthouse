package driver

import "errors"

type params map[string]interface{}

func (p params) getBool(key string) (bool, error) {
	val, ok := p[key]
	if !ok {
		return false, errors.New("Parameter '" + key + "' not found")
	}

	b, ok := val.(bool)
	if !ok {
		return false, errors.New("Parameter '" + key + "' invalid")
	}
	return b, nil
}

func (p params) getFloat(key string) (float64, error) {
	val, ok := p[key]
	if !ok {
		return 0, errors.New("Parameter '" + key + "' not found")
	}

	b, ok := val.(float64)
	if !ok {
		return 0, errors.New("Parameter '" + key + "' invalid")
	}
	return b, nil
}
