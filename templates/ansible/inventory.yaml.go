package ansible

type InventoryKeyValue struct{}
type InventoryValues struct {
	Values []InventoryKeyValue
}

var DefaultInventoryValues = InventoryValues{
	Values: []InventoryKeyValue{},
}

const InventoryFile = `---
# Tip: Use group_vars/group_name.yml to add in group variables for group_name
# Tip: Use host_vars/192.168.3.155.yml to add some host variables for the host with IP address 192.168.3.155.yml
# group_name:
#   hosts:
#     192.168.[0:8].[150:200]:
#   children:
#     other_group_names:
`
