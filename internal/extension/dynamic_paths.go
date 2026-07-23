package extension

import (
	"os"
	"path/filepath"

	"github.com/bluenviron/mediamtx/internal/conf"
	"github.com/bluenviron/mediamtx/internal/conf/yamlwrapper"
	"github.com/goccy/go-yaml"
)

// DynamicPathsConfig는 동적으로 추가된 채널(paths) 데이터만을 YAML로 직렬화/역직렬화하기 위한 구조체입니다.
type DynamicPathsConfig struct {
	Paths map[string]*conf.OptionalPath `yaml:"paths" json:"paths"`
}

func getPathsFilePath(fpath string) string {
	pathsFile := filepath.Join("data", "paths.yml")
	if fpath != "" {
		pathsFile = filepath.Join(filepath.Dir(fpath), "data", "paths.yml")
	}
	return pathsFile
}

// staticPathKeys는 메인 설정 파일(mediamtx.yml)에서 기본적으로 불러온 채널들의 목록을 기억합니다.
// 이 목록에 있는 채널들은 동적 데이터 파일(paths.yml)에 저장되지 않도록 필터링됩니다.
var staticPathKeys = make(map[string]bool)

// LoadPaths는 별도의 YAML 파일(paths.yml)에서 동적 채널 정보(OptionalPaths)를 불러옵니다.
func LoadPaths(c *conf.Conf, fpath string) {
	// 현재 c.OptionalPaths에 있는 것들은 순수하게 mediamtx.yml에서 온 것들입니다.
	// 이를 기억해 둡니다.
	staticPathKeys = make(map[string]bool)
	for k := range c.OptionalPaths {
		staticPathKeys[k] = true
	}

	pathsFile := getPathsFilePath(fpath)

	byts, err := os.ReadFile(pathsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 파일이 없으면 기본 빈 파일을 생성해 둡니다.
			_ = SavePaths(c, fpath)
		}
		return
	}

	var d DynamicPathsConfig
	err = yamlwrapper.Unmarshal(byts, &d)
	if err != nil {
		return
	}

	if c.OptionalPaths == nil {
		c.OptionalPaths = make(map[string]*conf.OptionalPath)
	}
	for k, v := range d.Paths {
		c.OptionalPaths[k] = v
	}
}

// SavePaths는 메모리의 동적 채널 정보(OptionalPaths)를 별도의 YAML 파일(paths.yml)에 저장합니다.
func SavePaths(c *conf.Conf, fpath string) error {
	pathsFile := getPathsFilePath(fpath)

	// data 폴더가 없을 수 있으므로 생성
	if err := os.MkdirAll(filepath.Dir(pathsFile), 0755); err != nil {
		return err
	}

	// mediamtx.yml에 있는 기본 path들은 제외하고 순수하게 동적으로 추가된 것만 추려냅니다.
	dynamicOnly := make(map[string]any)
	for k, v := range c.OptionalPaths {
		if !staticPathKeys[k] {
			// v (*conf.OptionalPath) 대신 v.Values를 저장하여 불필요한 "values:" 래퍼를 제거합니다.
			dynamicOnly[k] = v.Values
		}
	}

	d := struct {
		Paths map[string]any `yaml:"paths"`
	}{
		Paths: dynamicOnly,
	}
	
	byts, err := yaml.Marshal(d)
	if err != nil {
		return err
	}
	return os.WriteFile(pathsFile, byts, 0644)
}
