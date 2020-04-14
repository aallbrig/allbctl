package ansible

type GroupKeyValue struct {
	Key   string
	Value string
}
type GroupValues struct {
	Values []GroupKeyValue
}

var DefaultGroupValues = GroupValues{
	Values: []GroupKeyValue{
		{
			Key:   "key1",
			Value: "Value1",
		},
		{
			Key:   "key2",
			Value: "Value2",
		},
		{
			Key:   "key3",
			Value: "Value3",
		},
	},
}

const GroupVarFile = `---
{{ range .Values -}}
{{ .Key }}: "{{ .Value }}"
{{ end -}}
`
