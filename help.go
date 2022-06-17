package main

import "encoding/json"

func fatal(err string) []byte {
	jsonResp, _ := json.Marshal(Message{Message: err})
	return jsonResp
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
