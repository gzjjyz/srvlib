package main

import (
	"encoding/json"
	"fmt"
	"github.com/995933447/std-go/scan"
	"github.com/gzjjyz/srvlib/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"os"
	"path"
	"reflect"
	"strings"
)

type Cfg struct {
	CommonJsonPath string `json:"common_json_path"` // 通用表 JSON 所在文件夹
	CommonJsonName string `json:"common_json_name"` // 通用表 JSON 文件名
	OutputStPath   string `json:"output_st_path"`   // 生成结构体路径
	OutputStName   string `json:"output_st_name"`   // 生成结构体文件名
	PkgName        string `json:"pkg_name"`         // 包名
	TypeFlag       string `json:"type_flag"`        // 类型标识 - 没有表示不需要后端生成, 只支持组合:基本类型 + 数组
}

const (
	// 类型标识 - 没有表示不需要后端生成, 只支持组合:基本类型 + 数组
	Bool    = 1
	Int     = 2
	Int8    = 3
	Int16   = 4
	Int32   = 5
	Int64   = 6
	Uint    = 7
	Uint8   = 8
	Uint16  = 9
	Uint32  = 10
	Uint64  = 11
	Float32 = 12
	Float64 = 13
	Slice   = 14
	String  = 15
)

func genTyeStr(enum uint32) string {
	switch enum {
	case Bool:
		return reflect.Bool.String()
	case Int:
		return reflect.Int.String()
	case Int8:
		return reflect.Int8.String()
	case Int16:
		return reflect.Int16.String()
	case Int32:
		return reflect.Int32.String()
	case Int64:
		return reflect.Int64.String()
	case Uint:
		return reflect.Uint.String()
	case Uint8:
		return reflect.Uint8.String()
	case Uint16:
		return reflect.Uint16.String()
	case Uint32:
		return reflect.Uint32.String()
	case Uint64:
		return reflect.Uint64.String()
	case Float32:
		return reflect.Float32.String()
	case Float64:
		return reflect.Float64.String()
	case Slice:
		return reflect.Slice.String()
	case String:
		return reflect.String.String()
	}
	return ""
}

const (
	defaultPrefix = "cfg"
	defaultSuffix = "yaml"
)

var global *Cfg

func loadGlobal() error {
	viper.SetConfigName(defaultPrefix)
	viper.SetConfigType(defaultSuffix)
	viper.AddConfigPath(scan.OptStrDefault("p", ""))
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	val := viper.Get("cfg")
	var cfg Cfg
	err = json2St(val, &cfg)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return err
	}
	global = &cfg
	return nil
}

func json2St(re interface{}, out interface{}) error {
	marshal, err := json.Marshal(re)
	if err != nil {
		fmt.Printf("err is : %v\n", err)
		return err
	}

	err = json.Unmarshal(marshal, out)
	if err != nil {
		fmt.Printf("err is : %v\n", err)
		return err
	}
	return nil
}

type StGenTemplate struct {
	Pkg    string   `json:"pkg"`
	StName string   `json:"st_name"`
	Fs     []*Field `json:"fs"`
}

type Field struct {
	Name string `json:"name"`
	Typ  string `json:"typ"`
}

func main() {
	err := loadGlobal()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	tmp := make(map[string]map[string]interface{})
	data, err := os.ReadFile(path.Join(global.CommonJsonPath, global.CommonJsonName))
	if err != nil {
		fmt.Printf("未找到配置文件数据:[%s] %s\n", path.Join(global.CommonJsonPath, global.CommonJsonName), err)
		return
	}

	if err := jsoniter.Unmarshal(data, &tmp); err != nil {
		fmt.Printf("load %s Unmarshal json error:%s\n", path.Join(global.CommonJsonPath, global.CommonJsonName), err)
		return
	}

	genTemplate := &StGenTemplate{
		Pkg:    global.PkgName,
		StName: "CommonStConf",
	}

	for s, m := range tmp {
		f := Field{
			Name: s,
		}
		//for _, valFlag := range global.ValFlag {
		typFlag, ok := m[global.TypeFlag]
		if !ok || typFlag == nil || fmt.Sprintf("%v", typFlag) == "" {
			continue
		}
		split := strings.Split(fmt.Sprintf("%v", typFlag), ",")
		if len(split) == 2 {
			atoUint32 := utils.AtoUint32(split[0])
			str := genTyeStr(atoUint32)
			f.Typ = "[]" + str
		} else if len(split) == 1 {
			atoUint32 := utils.AtoUint32(split[0])
			f.Typ = genTyeStr(atoUint32)
		} else {
			continue
		}
		if len(f.Typ) == 0 {
			continue
		}
		genTemplate.Fs = append(genTemplate.Fs, &f)
	}
	outputGoFile(genTemplate)
}

func outputGoFile(st *StGenTemplate) {
	var strs strings.Builder
	strs.WriteString(fmt.Sprintf("package %s", st.Pkg))
	strs.WriteString("\n")
	strs.WriteString("\n")
	strs.WriteString(importStr)
	strs.WriteString("\n")
	strs.WriteString("\n")
	strs.WriteString(fmt.Sprintf("type %s struct {", st.StName))
	strs.WriteString("\n")
	for _, f := range st.Fs {
		strs.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`", toUpFirstStr(f.Name), f.Typ, f.Name))
		strs.WriteString("\n")
	}
	strs.WriteString("}")
	strs.WriteString("\n")
	strs.WriteString(loadTempStr)
	strs.WriteString("\n")

	var outputFile = path.Join(global.OutputStPath, global.OutputStName)
	err := createAndWriteFile(outputFile, strs.String())
	if err != nil {
		fmt.Printf("err is %v\n", err)
		return
	}
}

func createAndWriteFile(absTargetFilePath, content string) error {
	f, err := os.Create(absTargetFilePath)
	if err != nil {
		fmt.Printf("err is %v\n", err)
		return err
	}
	if _, err := f.Write([]byte(content)); err != nil {
		fmt.Printf("err is %v\n", err)
		return err
	}
	return nil
}

func toUpFirstStr(str string) string {
	return strings.ToUpper(str[:1]) + str[1:]
}

var importStr = `
import (
	"fmt"
	"github.com/gzjjyz/srvlib/utils"
	jsoniter "github.com/json-iterator/go"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)
`

var loadTempStr = `

var CommonStConfMgr *CommonStConf

// 增强配制表结构
func LoadCommonStConf() bool {
	tmp := new(CommonStConf)
	filePath := utils.GetCurrentDir() + "config/" + strings.ToLower("CommonConfig") + ".json"

	jsonTmp := make(map[string]map[string]interface{})
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("未找到配置文件数据:[%s] %s\n", filePath, err)
		return false
	}
	if err := jsoniter.Unmarshal(data, &jsonTmp); err != nil {
		fmt.Printf("load %s Unmarshal json error:%s\n", filePath, err)
		return false
	}

	initJsonDataToSt(jsonTmp, tmp)

	CommonStConfMgr = tmp
	return true
}

func GetCommonStConf() (*CommonStConf, bool) {
	return CommonStConfMgr, CommonStConfMgr != nil
}

// 手动 将 json 数据赋值
func initJsonDataToSt(jsonTmp map[string]map[string]interface{}, tmp *CommonStConf) {
	valueOf := reflect.ValueOf(tmp).Elem()
	elem := reflect.TypeOf(tmp).Elem()
	for i := 0; i < elem.NumField(); i++ {
		tStField := elem.Field(i)
		vStField := valueOf.Field(i)
		jsonFieldKey := tStField.Tag.Get("json")
		split := strings.Split(jsonFieldKey, ",")
		if len(split) > 1 {
			jsonFieldKey = split[0]
		}

		m, ok := jsonTmp[jsonFieldKey]
		if !ok {
			continue
		}

		for k, v := range m {
			if k == "key" || k == "type" {
				continue
			}
			switch tStField.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if k == "u32" {
					var atoUint32 = utils.AtoUint32(fmt.Sprintf("%v", v))
					vStField.Set(reflect.ValueOf(atoUint32))
				}
			case reflect.String:
				if k == "str" {
					vStField.Set(reflect.ValueOf(v))
				}
			case reflect.Float32, reflect.Float64:
				if k == "f32" {
					vStField.Set(reflect.ValueOf(v))
				}
			case reflect.Slice:
				// 目前只有 uint32
				if k == "u32Vec" {
					array := strings.Split(fmt.Sprintf("%v", v), ",")
					intList := make([]uint32, 0, len(array))
					for i := 0; i < len(array); i++ {
						intVal, err := strconv.ParseInt(array[i], 10, 32)
						if err != nil {
							continue
						}
						intList = append(intList, uint32(intVal))
					}
					vStField.Set(reflect.ValueOf(intList))
				}
			}
		}

	}
}
`
