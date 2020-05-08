package ansible

type KeyValue struct {
	Key   string
	Value string
}

type KeyValuePairs struct {
	Values []KeyValue
}

var DefaultKeyValue = KeyValuePairs{
	Values: []KeyValue{
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

