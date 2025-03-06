package taxescalc

import "math"

const (
	InssFirstCutValue  = 1518.00
	InssSecondCutValue = 2793.88
	InssThirdCutValue  = 4190.83
	InssForthCutValue  = 8157.41

	InssFirstCutAliquot  = 0.075
	InssSecondCutAliquot = 0.09
	InssThirdCutAliquot  = 0.12
	InssForthCutAliquot  = 0.14

	IrrfFirstCutValue  = 2259.20
	IrrfSecondCutValue = 2826.65
	IrrfThirdCutValue  = 3751.05
	IrrfForthCutValue  = 4664.68

	IrrfFirstCutAliquot  = 0.075
	IrrfSecondCutAliquot = 0.15
	IrrfThirdCutAliquot  = 0.225
	IrrfForthCutAliquot  = 0.275

	DiscountSimplified   = 564.80
	DiscountPerDependant = 189.59

	IrrfPlrFirstCutValue  = 7640.80
	IrrfPlrSecondCutValue = 9922.28
	IrrfPlrThirdCutValue  = 13167.00
	IrrfPlrForthCutValue  = 16380.38

	IrrfPlrFirstCutAliquot  = 0.075
	IrrfPlrSecondCutAliquot = 0.15
	IrrfPlrThirdCutAliquot  = 0.225
	IrrfPlrForthCutAliquot  = 0.275
)

func NetSalary(salary, benefitValue, foodPensionValue, otherDiscounts, extraPercentage float64, dependantsNumber, extraHours, daysWorked, hoursPerDay int) float64 {
	if extraHours > 0 {
		salary += (SalaryPerHour(salary, daysWorked, hoursPerDay) * float64(extraHours)) * extraPercentage
	}
	inssValue := InssTax(salary)
	irrfValue := IrrfSalaryTax(salary, foodPensionValue, dependantsNumber)

	return RoundFloat(salary-inssValue-irrfValue-benefitValue-otherDiscounts, 2)
}

// Inss is the salary with progressive aliquot
// 1.401,99 -> 5%
// 1.402,00 - 2.666,68 -> 9%
// 2.666,69 - 4.000,03 -> 12%
// 4.000,04 - 7.786,02 -> 14%
func InssTax(salary float64) float64 {
	var tax float64
	if salary <= InssFirstCutValue {
		tax = salary * InssFirstCutAliquot
		return RoundFloat(tax, 2)
	} else {
		tax = InssFirstCutValue * InssFirstCutAliquot
	}

	if salary <= InssSecondCutValue {
		tax += (salary - InssFirstCutValue) * InssSecondCutAliquot
		return RoundFloat(tax, 2)
	} else {
		tax += (InssSecondCutValue - InssFirstCutValue) * InssSecondCutAliquot
	}

	if salary <= InssThirdCutValue {
		tax += (salary - InssSecondCutValue) * InssThirdCutAliquot
		return RoundFloat(tax, 2)
	} else {
		tax += (InssThirdCutValue - InssSecondCutValue) * InssThirdCutAliquot
	}

	if salary <= InssForthCutValue {
		tax += (salary - InssThirdCutValue) * InssForthCutAliquot
	} else {
		tax += (InssForthCutValue - InssThirdCutValue) * InssForthCutAliquot
	}

	return RoundFloat(tax, 2)
}

// Irrf is the salary with progressive aliquot
// 2.259,20 -> 0%
// 2.259,21 - 2.826,65 -> 7,5%
// 2.826,66 - 3.751,05 -> 15%
// 3.751,06 - 4.664,68 -> 22,5%
// 4.664,69 -> 27,5%
func IrrfSalaryTax(salary, foodPensionValue float64, dependantsNumber int) float64 {
	totalDeduction := InssTax(salary) + foodPensionValue + (float64(dependantsNumber) * DiscountPerDependant)
	if totalDeduction < DiscountSimplified {
		totalDeduction = DiscountSimplified
	}

	deductedSalary := salary - totalDeduction
	var tax float64
	if deductedSalary <= IrrfFirstCutValue {
		return RoundFloat(tax, 2)
	}

	if deductedSalary <= IrrfSecondCutValue {
		tax += (deductedSalary - IrrfFirstCutValue) * IrrfFirstCutAliquot
		return RoundFloat(tax, 2)
	} else {
		tax += (IrrfSecondCutValue - IrrfFirstCutValue) * IrrfFirstCutAliquot
	}

	if deductedSalary <= IrrfThirdCutValue {
		tax += (deductedSalary - IrrfSecondCutValue) * IrrfSecondCutAliquot
		return RoundFloat(tax, 2)
	} else {
		tax += (IrrfThirdCutValue - IrrfSecondCutValue) * IrrfSecondCutAliquot
	}

	if deductedSalary <= IrrfForthCutValue {
		tax += (deductedSalary - IrrfThirdCutValue) * IrrfThirdCutAliquot
		return RoundFloat(tax, 2)
	} else {
		tax += (IrrfForthCutValue - IrrfThirdCutValue) * IrrfThirdCutAliquot
	}

	tax += (deductedSalary - IrrfForthCutValue) * IrrfForthCutAliquot

	return RoundFloat(tax, 2)
}

// Irrf Plr tax is the prlValue with progressive aliquot
// 7.640,79 -> 0%
// 7.640,80 - 9.922,27 -> 7,5%
// 9.922,28 - 13.166,99 -> 15%
// 13.167,00 - 16.380,37 -> 22,5%
// 16.380,38 -> 27,5%
func IrrfPlrTax(prlValue float64) float64 {
	var tax float64
	if prlValue <= IrrfPlrFirstCutValue {
		return RoundFloat(tax, 2)
	}

	if prlValue <= IrrfPlrSecondCutValue {
		tax += (prlValue - IrrfPlrFirstCutValue) * IrrfPlrFirstCutAliquot
		return RoundFloat(tax, 2)
	} else {
		tax += (IrrfPlrSecondCutValue - IrrfPlrFirstCutValue) * IrrfPlrFirstCutAliquot
	}

	if prlValue <= IrrfPlrThirdCutValue {
		tax += (prlValue - IrrfPlrSecondCutValue) * IrrfPlrSecondCutAliquot
		return RoundFloat(tax, 2)
	} else {
		tax += (IrrfPlrThirdCutValue - IrrfPlrSecondCutValue) * IrrfPlrSecondCutAliquot
	}

	if prlValue <= IrrfPlrForthCutValue {
		tax += (prlValue - IrrfPlrThirdCutValue) * IrrfPlrThirdCutAliquot
		return RoundFloat(tax, 2)
	} else {
		tax += (IrrfPlrForthCutValue - IrrfPlrThirdCutValue) * IrrfPlrThirdCutAliquot
	}

	tax += (prlValue - IrrfPlrForthCutValue) * IrrfPlrForthCutAliquot

	return RoundFloat(tax, 2)
}

// (gross salary + ⅓ of gross salary)/30 x Number of Days – discounts of IRRF and INSS
func VacationFracionatedSalary(salary, benefitValue, foodPensionValue, otherDiscounts float64, dependantsNumber, requestedDays int) float64 {
	salary += RoundFloat(salary/3, 2)
	salary = (salary / 30) * float64(requestedDays)
	return RoundFloat(salary-InssTax(salary)-IrrfSalaryTax(salary, foodPensionValue, dependantsNumber)-benefitValue, 2)
}

// (gross salary + ⅓ of gross salary) – discounts of IRRF and INSS
func VacationSalary(salary, benefitValue, foodPensionValue, otherDiscounts float64, dependantsNumber int) float64 {
	salary += RoundFloat(salary/3, 2)
	return RoundFloat(salary-InssTax(salary)-IrrfSalaryTax(salary, foodPensionValue, dependantsNumber)-benefitValue, 2)
}

// gross salary – discounts of IRRF and INSS / 12 * workedMonths
func ThirteenthSalary(salary float64, dependantsNumber, workedMonths int) float64 {
	if workedMonths == 0 {
		workedMonths = 12
	}
	return RoundFloat((salary/12)*float64(workedMonths)-InssTax(salary)-IrrfSalaryTax(salary, 0.00, dependantsNumber), 2)
}

func FgtsPerMonth(salary float64) float64 {
	return RoundFloat((salary * 0.08), 2)
}

func SalaryPerHour(salary float64, daysWorked, hoursPerDay int) float64 {
	return RoundFloat((salary/float64(daysWorked))/float64(hoursPerDay), 2)
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
