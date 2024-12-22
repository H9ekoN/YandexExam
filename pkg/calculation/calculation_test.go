package calculation_test

import (
	"testing"

	"github.com/H9ekoN/YandexExam_go/pkg/calculation"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
		{
			name:           "final",
			expression:     "(2*(3+4)-5)/2",
			expectedResult: 4.5,
		},
	}

	testCasesFail := []struct {
		name       string
		expression string
		wantErr    string
	}{
		{
			name:       "input error",
			expression: "a",
			wantErr:    "неизвестный символ",
		},
		{
			name:       "input error",
			expression: "2+2//*()3",
			wantErr:    "неизвестный символ",
		},
		{
			name:       "input error",
			expression: "****4",
			wantErr:    "неизвестный символ",
		},
		{
			name:       "division by zero",
			expression: "7/0",
			wantErr:    "на ноль делить нельзя",
		},
		{
			name:       "input error",
			expression: "(5+5",
			wantErr:    "ошибочка в количестве скобок",
		},
		{
			name:       "overflow error",
			expression: "999999999999999999999999999999*999999999999999999999999999999*999999",
			wantErr:    "переполнение",
		},
	}

	for _, tc := range testCasesSuccess {
		t.Run(tc.name, func(t *testing.T) {
			val, err := calculation.Calc(tc.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error: %v", tc.expression, err)
			}
			if val != tc.expectedResult {
				t.Fatalf("%f should be equal %f", val, tc.expectedResult)
			}
		})
	}

	for _, tc := range testCasesFail {
		t.Run(tc.name, func(t *testing.T) {
			_, err := calculation.Calc(tc.expression)
			if err == nil {
				t.Fatalf("expected error for %s", tc.expression)
			}
			if err.Error() != tc.wantErr {
				t.Fatalf("got error %v, want %v", err, tc.wantErr)
			}
		})
	}
}
