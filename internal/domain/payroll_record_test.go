package domain

import (
	"testing"
	"time"
)

func TestPayrollRecord_Validate(t *testing.T) {
	now := time.Now()
	valid := func() PayrollRecord {
		return PayrollRecord{
			EmployeeContractID: 1,
			PeriodStart:        now.AddDate(0, -1, 0),
			PeriodEnd:          now,
		}
	}

	tests := []struct {
		name    string
		record  PayrollRecord
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid payroll record - should pass",
			record:  valid(),
			wantErr: false,
		},
		{
			name: "zero contract id - should fail",
			record: func() PayrollRecord {
				pr := valid()
				pr.EmployeeContractID = 0
				return pr
			}(),
			wantErr: true,
			errMsg:  "contract is required",
		},
		{
			name: "zero period start - should fail",
			record: func() PayrollRecord {
				pr := valid()
				pr.PeriodStart = time.Time{}
				return pr
			}(),
			wantErr: true,
			errMsg:  "period is required",
		},
		{
			name: "zero period end - should fail",
			record: func() PayrollRecord {
				pr := valid()
				pr.PeriodEnd = time.Time{}
				return pr
			}(),
			wantErr: true,
			errMsg:  "period is required",
		},
		{
			name: "end before start - should fail",
			record: func() PayrollRecord {
				pr := valid()
				pr.PeriodStart = now
				pr.PeriodEnd = now.AddDate(0, -1, 0)
				return pr
			}(),
			wantErr: true,
			errMsg:  ErrInvalidPeriod.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.record.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Expected error '%s' but got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
