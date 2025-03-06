package taxescalc_test

import (
	"testing"

	taxes_calc "github.com/iobrasil/taxes-calc"
)

func TestInssTax(t *testing.T) {
	tests := []struct {
		name     string
		salary   float64
		expected float64
	}{
		{"Below first cut", 1000.00, 1000.00 * taxes_calc.InssFirstCutAliquot},
		{"At first cut", taxes_calc.InssFirstCutValue, taxes_calc.InssFirstCutValue * taxes_calc.InssFirstCutAliquot},
		{"Between first and second cut", 2500.00, (taxes_calc.InssFirstCutValue * taxes_calc.InssFirstCutAliquot) + ((2500.00 - taxes_calc.InssFirstCutValue) * taxes_calc.InssSecondCutAliquot)},
		{"Between second and third cut", 3000.00, (taxes_calc.InssFirstCutValue * taxes_calc.InssFirstCutAliquot) + ((taxes_calc.InssSecondCutValue - taxes_calc.InssFirstCutValue) * taxes_calc.InssSecondCutAliquot) + ((3000.00 - taxes_calc.InssSecondCutValue) * taxes_calc.InssThirdCutAliquot)},
		{"Between third and fourth cut", 5000.00, (taxes_calc.InssFirstCutValue * taxes_calc.InssFirstCutAliquot) + ((taxes_calc.InssSecondCutValue - taxes_calc.InssFirstCutValue) * taxes_calc.InssSecondCutAliquot) + ((taxes_calc.InssThirdCutValue - taxes_calc.InssSecondCutValue) * taxes_calc.InssThirdCutAliquot) + ((5000.00 - taxes_calc.InssThirdCutValue) * taxes_calc.InssForthCutAliquot)},
		{"Above fourth cut", 12000.00, (taxes_calc.InssFirstCutValue * taxes_calc.InssFirstCutAliquot) + ((taxes_calc.InssSecondCutValue - taxes_calc.InssFirstCutValue) * taxes_calc.InssSecondCutAliquot) + ((taxes_calc.InssThirdCutValue - taxes_calc.InssSecondCutValue) * taxes_calc.InssThirdCutAliquot) + ((taxes_calc.InssForthCutValue - taxes_calc.InssThirdCutValue) * taxes_calc.InssForthCutAliquot)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taxes_calc.InssTax(tt.salary)
			if result != taxes_calc.RoundFloat(tt.expected, 2) {
				t.Errorf("For salary %.2f, expected %.2f, but got %.2f", tt.salary, tt.expected, result)
			}
		})
	}
}

func TestIrrfSalaryTax(t *testing.T) {
	tests := []struct {
		name             string
		salary           float64
		foodPensionValue float64
		dependantsNumber int
		expected         float64
	}{
		{"Below first cut", 1000.00, 0.00, 0, 0.00},
		{"At first cut", taxes_calc.IrrfFirstCutValue, 0.00, 0, 0.00},
		{"Between first and second cut", 2900.00, 0.00, 0, 5.70},
		{"Between second and third cut", 3000.00, 0.00, 0, 13.20},
		{"Between third and fourth cut", 5000.00, 0.00, 0, 335.15},
		{"Above fourth cut", 8200.00, 0.00, 0, 1097.30},
		{"Above fourth cut", 12000.00, 0.00, 0, 2142.30},
		{"Between second and third cut with dependant", 3500.00, 0.00, 1, 58.84},
		{"Between second and third cut with dependant and pension", 3500.00, 100.00, 1, 53.11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taxes_calc.IrrfSalaryTax(tt.salary, tt.foodPensionValue, tt.dependantsNumber)
			if result != taxes_calc.RoundFloat(tt.expected, 2) {
				t.Errorf("For salary %.2f, expected %.2f, but got %.2f", tt.salary, tt.expected, result)
			}
		})
	}
}

func TestIrrfPlrTax(t *testing.T) {
	tests := []struct {
		name     string
		prlValue float64
		expected float64
	}{
		{"Below first cut", 1000.00, 0.00},
		{"At first cut", taxes_calc.IrrfPlrFirstCutValue, 0.00},
		{"Between first and second cut", 9000.00, 101.94},
		{"Between second and third cut", 10000.00, 182.77},
		{"Between third and fourth cut", 15000.00, 1070.24},
		{"Above fourth cut", 18200.00, 1881.23},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taxes_calc.IrrfPlrTax(tt.prlValue)
			if result != taxes_calc.RoundFloat(tt.expected, 2) {
				t.Errorf("For salary %.2f, expected %.2f, but got %.2f", tt.prlValue, tt.expected, result)
			}
		})
	}
}

func TestNetSalary(t *testing.T) {
	tests := []struct {
		name             string
		salary           float64
		benefitValue     float64
		foodPensionValue float64
		otherDiscounts   float64
		extraPercentage  float64
		dependantsNumber int
		extraHours       int
		daysWorked       int
		hoursPerDay      int
		expected         float64
	}{
		{"Net Salary", 8000.00, 0.00, 0.00, 0.00, 0.00, 0, 0, 0, 0, 6022.04},
		{"Net Salary with extra time", 8000.00, 0.00, 0.00, 0.00, 1.50, 0, 8, 8, 21, 6420.36},
		{"Net Salary with food pension and extra time", 8000.00, 0.00, 100.00, 0.00, 1.50, 0, 8, 8, 21, 6447.86},
		{"Net Salary with benefits, food pension and extra time", 8000.00, 100.00, 100.00, 0.00, 1.50, 0, 8, 8, 21, 6347.86},
		{"Net Salary with food pension and dependants", 8000.00, 0.00, 100.00, 0.00, 0, 2, 0, 0, 0, 6153.82},
		{"Net Salary with benefits, food pension, other discounts and extra time", 8000.00, 100.00, 100.00, 100.00, 1.50, 0, 8, 8, 21, 6247.86},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taxes_calc.NetSalary(
				tt.salary,
				tt.benefitValue,
				tt.foodPensionValue,
				tt.otherDiscounts,
				tt.extraPercentage,
				tt.dependantsNumber,
				tt.extraHours,
				tt.daysWorked,
				tt.hoursPerDay,
			)
			if result != taxes_calc.RoundFloat(tt.expected, 2) {
				t.Errorf("For salary %.2f, expected %.2f, but got %.2f", tt.salary, tt.expected, result)
			}
		})
	}
}

func TestSalaryPerHour(t *testing.T) {
	tests := []struct {
		name        string
		salary      float64
		daysWorked  int
		hoursPerDay int
		expected    float64
	}{
		{"Validate", 8000.00, 21, 8, 47.62},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taxes_calc.SalaryPerHour(tt.salary, tt.daysWorked, tt.hoursPerDay)
			if result != taxes_calc.RoundFloat(tt.expected, 2) {
				t.Errorf("For salary %.2f, expected %.2f, but got %.2f", tt.salary, tt.expected, result)
			}
		})
	}
}

func TestFgtsPerMonth(t *testing.T) {
	tests := []struct {
		name     string
		salary   float64
		expected float64
	}{
		{"Validate", 10000.00, 800.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taxes_calc.FgtsPerMonth(tt.salary)
			if result != taxes_calc.RoundFloat(tt.expected, 2) {
				t.Errorf("For salary %.2f, expected %.2f, but got %.2f", tt.salary, tt.expected, result)
			}
		})
	}
}

func TestVacationValue(t *testing.T) {
	tests := []struct {
		name             string
		salary           float64
		benefitValue     float64
		foodPensionValue float64
		otherDiscounts   float64
		dependantsNumber int
		expected         float64
	}{
		{"Validate vacation salary", 4500.00, 0.00, 0.00, 0.00, 0, 4775.04},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taxes_calc.VacationSalary(tt.salary, tt.benefitValue, tt.foodPensionValue, tt.otherDiscounts, tt.dependantsNumber)
			if result != taxes_calc.RoundFloat(tt.expected, 2) {
				t.Errorf("For salary %.2f, expected %.2f, but got %.2f", tt.salary, tt.expected, result)
			}
		})
	}
}

func TestVacationFracionatedSalaryValue(t *testing.T) {
	tests := []struct {
		name             string
		salary           float64
		benefitValue     float64
		foodPensionValue float64
		otherDiscounts   float64
		dependantsNumber int
		requestedDays    int
		expected         float64
	}{
		{"Validate vacation salary", 3500.00, 0.00, 0.00, 0.00, 0, 20, 2822.83},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taxes_calc.VacationFracionatedSalary(
				tt.salary,
				tt.benefitValue,
				tt.foodPensionValue,
				tt.otherDiscounts,
				tt.dependantsNumber,
				tt.requestedDays,
			)
			if result != taxes_calc.RoundFloat(tt.expected, 2) {
				t.Errorf("For salary %.2f, expected %.2f, but got %.2f", tt.salary, tt.expected, result)
			}
		})
	}
}

func TestThirteenthSalary(t *testing.T) {
	tests := []struct {
		name             string
		salary           float64
		dependantsNumber int
		workedMonths     int
		expected         float64
	}{
		{"Validate thirteenth salary", 3500.00, 0, 12, 3127.75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := taxes_calc.ThirteenthSalary(
				tt.salary,
				tt.dependantsNumber,
				tt.workedMonths,
			)
			if result != taxes_calc.RoundFloat(tt.expected, 2) {
				t.Errorf("For salary %.2f, expected %.2f, but got %.2f", tt.salary, tt.expected, result)
			}
		})
	}
}
