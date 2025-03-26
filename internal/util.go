package internal

import (
	"encoding/json"
	"log"
	"math/rand"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func spw(v any) {
	log.Println(spew.Sdump(v))
}

func ptr[T any](value T) *T {
	return &value
}

func mapStringTToString[T any](d *diag.Diagnostics, v map[string]T) *string {
	ret, err := json.Marshal(v)
	if err != nil {
		d.AddError("Failed to marshal a map[ring]any to a string", err.Error())
		return nil
	}
	strRet := string(ret)
	return &strRet

}

func stringToMapStringT[T any](d *diag.Diagnostics, v string) *map[string]T {
	var ret map[string]T
	err := json.Unmarshal([]byte(v), &ret)
	if err != nil {
		d.AddError("Unable to convert a json string to a map[string]T", err.Error())
		return nil
	}
	return &ret

}
