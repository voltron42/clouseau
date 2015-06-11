package common

func castAsError(obj interface{}) error {
	if obj == nil {
		return nil
	}
	out, ok := obj.(error)
	if !ok {
		panic("Cannot be cast to error.")
	}
	return out
}

func castAsBytes(obj interface{}) []byte {
	if obj == nil {
		return nil
	}
	out, ok := obj.([]byte)
	if !ok {
		panic("Cannot be cast to byte array.")
	}
	return out
}

func castAsStrings(obj interface{}) []string {
	if obj == nil {
		return nil
	}
	out, ok := obj.([]string)
	if !ok {
		panic("Cannot be cast to string array.")
	}
	return out
}
