package view

import (
	"encoding/json"
)

// func Convert(src, dst interface{}) error {
// 	buf := bytes.Buffer{}
// 	err := gob.NewEncoder(&buf).Encode(src)
// 	if err != nil {
// 		return err
// 	}
// 	return gob.NewDecoder(&buf).Decode(dst)
// }

// Convert converts src to dst
func Convert(src, dst interface{}) error {
	v, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(v, dst)
}
