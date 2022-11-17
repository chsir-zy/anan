package framework

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"mime/multipart"

	"github.com/spf13/cast"
)

const defaultMultipartMemory = 32 << 20

//请求参数接口
type IRequest interface {
	//获取请求参数 /foo?a=1&b=2
	QueryInt(key string, def int) (int, bool)
	QueryInt64(key string, def int64) (int64, bool)
	QueryFloat64(key string, def float64) (float64, bool)
	QueryFloat32(key string, def float32) (float32, bool)
	QueryBool(key string, def bool) (bool, bool)
	QueryString(key string, def string) (string, bool)
	QueryStringSlice(key string, def []string) ([]string, bool)
	Query(key string) interface{}

	// 参数相关的 例如url:/user/:id
	ParamInt(key string, def int) (int, bool)
	ParamInt64(key string, def int64) (int64, bool)
	ParamFloat64(key string, def float64) (float64, bool)
	ParamFloat32(key string, def float32) (float32, bool)
	ParamBool(key string, def bool) (bool, bool)
	ParamString(key string, def string) (string, bool)
	Param(key string) interface{}

	// form表单相关的
	FormInt(key string, def int) (int, bool)
	FormInt64(key string, def int64) (int64, bool)
	FormFloat64(key string, def float64) (float64, bool)
	FormFloat32(key string, def float32) (float32, bool)
	FormBool(key string, def bool) (bool, bool)
	FormString(key string, def string) (string, bool)
	FormStringSlice(key string, def []string) ([]string, bool)
	FormFile(key string) (*multipart.FileHeader, error)
	Form(key string) interface{}

	// json格式
	BindJson(obj interface{}) error

	// xml格式
	BindXml(obj interface{}) error

	// 其它格式
	GetRawData() ([]byte, error)

	// header的基础数据
	Uri() string
	Method() string
	Host() string
	ClientIp() string

	Headers() map[string][]string
	Header(key string) (string, bool)

	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request.URL.Query() != nil {
		return map[string][]string(ctx.request.URL.Query())
	}

	return map[string][]string{}
}

func (ctx *Context) QueryInt(key string, def int) (int, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryInt64(key string, def int64) (int64, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryFloat64(key string, def float64) (float64, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryFloat32(key string, def float32) (float32, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryString(key string, def string) (string, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToString(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) QueryStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals, true
	}

	return def, false
}
func (ctx *Context) QueryBool(key string, def bool) (bool, bool) {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) Query(key string) interface{} {
	params := ctx.QueryAll()
	if vals, ok := params[key]; ok {
		return vals[0]
	}

	return nil
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		ctx.request.ParseForm()
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}
func (ctx *Context) FormInt(key string, def int) (int, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) FormInt64(key string, def int64) (int64, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToInt64(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) FormFloat64(key string, def float64) (float64, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat64(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) FormFloat32(key string, def float32) (float32, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToFloat32(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) FormBool(key string, def bool) (bool, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToBool(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) FormString(key string, def string) (string, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return cast.ToString(vals[0]), true
		}
	}

	return def, false
}
func (ctx *Context) FormStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		if len(vals) > 0 {
			return vals, true
		}
	}

	return def, false
}
func (ctx *Context) FormFile(key string) (*multipart.FileHeader, error) {
	if ctx.request.MultipartForm == nil {
		if err := ctx.request.ParseMultipartForm(defaultMultipartMemory); err != nil {
			return nil, err
		}
	}

	f, fh, err := ctx.request.FormFile(key)
	if err != nil {
		return nil, err
	}
	f.Close()

	return fh, nil
}
func (ctx *Context) Form(key string) interface{} {
	params := ctx.FormAll()
	if vals, ok := params[key]; ok {
		return vals[0]
	}
	return nil
}

func (ctx *Context) Param(key string) interface{} {
	if ctx.params == nil {
		return nil
	}

	if val, ok := ctx.params[key]; ok {
		return val
	}
	return nil
}

func (ctx *Context) ParamInt(key string, def int) (int, bool) {
	ret := ctx.Param(key)
	if ret == nil {
		return def, false
	}

	return cast.ToInt(ret), true
}
func (ctx *Context) ParamInt64(key string, def int64) (int64, bool) {
	ret := ctx.Param(key)
	if ret == nil {
		return def, false
	}

	return cast.ToInt64(ret), true
}
func (ctx *Context) ParamFloat64(key string, def float64) (float64, bool) {
	ret := ctx.Param(key)
	if ret == nil {
		return def, false
	}

	return cast.ToFloat64(ret), true
}
func (ctx *Context) ParamFloat32(key string, def float32) (float32, bool) {
	ret := ctx.Param(key)
	if ret == nil {
		return def, false
	}

	return cast.ToFloat32(ret), true
}
func (ctx *Context) ParamBool(key string, def bool) (bool, bool) {
	ret := ctx.Param(key)
	if ret == nil {
		return def, false
	}

	return cast.ToBool(ret), true
}
func (ctx *Context) ParamString(key string, def string) (string, bool) {
	ret := ctx.Param(key)
	if ret == nil {
		return def, false
	}

	return cast.ToString(ret), true
}

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		// 将body再次填充回ctx.request.Body
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

func (ctx *Context) BindXml(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		// 将body再次填充回ctx.request.Body
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = xml.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}
