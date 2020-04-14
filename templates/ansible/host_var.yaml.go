package ansible

type HostKeyValue struct {
	Key   string
	Value string
}
type HostValues struct {
	Values []HostKeyValue
}

var DefaultHostValues = HostValues{
	Values: []HostKeyValue{
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

const HostVarFile = `---
{{ range .Values -}}
{{ .Key }}: "{{ .Value }}"
{{ end -}}
`
