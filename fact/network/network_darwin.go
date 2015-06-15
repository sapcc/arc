package network

func (h Source) Facts() (map[string]interface{}, error) {

	return newFacts(), nil
}
