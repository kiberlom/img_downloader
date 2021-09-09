package parslink

import (
	"regexp"
)

func GetAllUrl(s string) ([]string, error) {

	rxg, err := regexp.Compile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	if err != nil {
		return nil, err
	}

	sr := rxg.FindAllString(s, -1)

	f := make(map[string]struct{})
	for _, v := range sr {
		if _, ok := f[v]; !ok {
			f[v] = struct{}{}
		}
	}

	var result []string

	for i := range f {
		result = append(result, i)
	}

	return result, nil
}
