package gini

import "fmt"

func fileToMap(data []string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	for i := 0; i < len(data); i++ {
		if len(data[i]) <= 0 {
			continue
		}
		if data[i][0] == '[' {
			name := getSectionName(data[i])
			i += 1
			result[name] = parseFields(data, &i)
		}
	}
	return result
}

func parseFields(data []string, i *int) map[string]string {
	result := make(map[string]string)
	for ; *i < len(data); *i++ {
		if len(data[*i]) <= 0 {
			continue
		}
		if data[*i][0] == '[' {
			*i -= 1
			return result
		}
		result[getFieldName(data[*i])] = getFieldValue(data[*i])
	}
	return result
}

func getSectionName(data string) string {
	result := ""
	for i := range data {
		if data[i] != '[' && data[i] != ']' {
			result += string(data[i])
		}
	}
	return result
}
func getFieldName(data string) string {
	result := ""
	for i := range data {
		if data[i] == '=' {
			break
		}
		if data[i] != ' ' && data[i] != '\t' {
			result += string(data[i])
		}
	}
	return result
}

func getFieldValue(data string) string {
	result := ""
	isValue := false
	for i := range data {
		if data[i] == '=' {
			isValue = true
			continue
		}
		if isValue && data[i] != ' ' && data[i] != '\t' {
			result += string(data[i])
		}
	}
	return result
}

func mapToFile(data map[string]map[string]string) []string {
	result := make([]string, 0)
	for i := range data {
		result = append(result, fmt.Sprintf("[%s]", i))
		for j := range data[i] {
			result = append(result, fmt.Sprintf("%s\t=\t%s", j, data[i][j]))
		}
	}
	return result
}
