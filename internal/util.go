package internal

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"
	"regexp"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	saneStringRegexOnce  sync.Once
	saneStringRegexValue *regexp.Regexp
)

func saneStringRegex() *regexp.Regexp {
	saneStringRegexOnce.Do(func() {
		var err error
		saneStringRegexValue, err = regexp.Compile(`^[a-zA-Z0-9\-_.]+$`)
		if err != nil {
			// we should never reach this code
			log.Panic(err)
		}
	})
	return saneStringRegexValue
}

func randStr32() *string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	res := string(b)
	return &res

}

func Spw(v any) {
	log.Println(spew.Sdump(v))
}

func ptr[T any](value T) *T {
	return &value
}

// Marshal a map[string]T to a string
//
// if v is nil or if Marshaling failed, return nil
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
//
// if v is nil or if Unmarshaling failed, return nil
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

// if unknown return nil, else return value
func boolOrNil(v types.Bool) *bool {
	if v.IsUnknown() {
		return nil
	}
	return v.ValueBoolPointer()
}

// wrapper function around `resp.Diagnostics.Append(resp.State.SetAttribute())` for *CreateResponse
func setCreateState(ctx context.Context, resp *resource.CreateResponse, rootPath string, val any) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(rootPath), val)...)
}

// wrapper function around `resp.Diagnostics.Append(resp.State.SetAttribute())` for *ReadResponse
func setReadState(ctx context.Context, resp *resource.ReadResponse, rootPath string, val any) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(rootPath), val)...)
}

// wrapper function around `resp.Diagnostics.Append(resp.State.SetAttribute())` for *UpdateResponse
func setUpdateState(ctx context.Context, resp *resource.UpdateResponse, rootPath string, val any) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(rootPath), val)...)
}
