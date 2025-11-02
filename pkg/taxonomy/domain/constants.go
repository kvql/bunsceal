package domain

// Use this to test compatibility in clients
const ApiVersion = "v1beta1"

// Define sensitivity levels
var SensitivityLevels = map[string]string{
	"A": "High",
	"B": "Medium",
	"C": "Low",
	"D": "N/A",
}

// Define Criticality levels
var CriticalityLevels map[string]string = map[string]string{
	"1": "Critical",
	"2": "High",
	"3": "Medium",
	"4": "Low",
	"5": "N/A",
}

// create slice of risk levels to track order
var SenseOrder = []string{"A", "B", "C", "D"}
var CritOrder = []string{"1", "2", "3", "4", "5"}
