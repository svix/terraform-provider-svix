package internal

import (
	"encoding/json"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func randStr32() *string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	res := string(b)
	return &res

}

// func spw(v any) {
// 	log.Println(spew.Sdump(v))
// }

func ptr[T any](value T) *T {
	return &value
}

// Marshal a map[string]T to a string
// if v is nil, return nil
func mapStringTToString[T any](d *diag.Diagnostics, v *map[string]T) *string {
	if v == nil {
		return nil
	}
	ret, err := json.Marshal(v)
	if err != nil {
		d.AddError("Failed to marshal a map[ring]any to a string", err.Error())
		return nil
	}
	strRet := string(ret)
	return &strRet

}

// Unmarshal a *string to a map[string]T
// if v is nil, return nil
func stringToMapStringT[T any](d *diag.Diagnostics, v *string) *map[string]T {
	if v == nil {
		return nil
	}
	var ret map[string]T
	err := json.Unmarshal([]byte(*v), &ret)
	if err != nil {
		d.AddError("Unable to convert a json string to a map[string]T", err.Error())
		return nil
	}
	return &ret

}

// if unknown return nil, else return value
func strOrNil(v types.String) *string {
	if v.IsUnknown() {
		return nil
	}
	return v.ValueStringPointer()
}
