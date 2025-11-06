package domain

// ApiVersion specifies the API version for compatibility testing in clients.
const ApiVersion = "v1beta1"

// SensitivityLevels defines the available sensitivity classification levels.
var SensitivityLevels = map[string]string{
	"A": "High",
	"B": "Medium",
	"C": "Low",
	"D": "N/A",
}

// CriticalityLevels defines the available criticality classification levels.
var CriticalityLevels = map[string]string{
	"1": "Critical",
	"2": "High",
	"3": "Medium",
	"4": "Low",
	"5": "N/A",
}

// SenseOrder defines the ordered list of sensitivity levels.
var SenseOrder = []string{"A", "B", "C", "D"}

// CritOrder defines the ordered list of criticality levels.
var CritOrder = []string{"1", "2", "3", "4", "5"}
