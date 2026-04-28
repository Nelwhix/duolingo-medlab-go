package request_test

import (
	"encoding/json"
	"testing"

	"github.com/Nelwhix/duolingo-medlab-go/pkg/request"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserValidation(t *testing.T) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	tests := []struct {
		name        string
		jsonInput   string
		expectErr   bool
		validateErr bool
	}{
		{
			name: "Valid all true",
			jsonInput: `{
				"learning_goal_pass_exams": true,
				"learning_goal_refresh_knowledge": true,
				"learning_goal_practice_daily": true
			}`,
			expectErr:   false,
			validateErr: false,
		},
		{
			name: "Valid all false",
			jsonInput: `{
				"learning_goal_pass_exams": false,
				"learning_goal_refresh_knowledge": false,
				"learning_goal_practice_daily": false
			}`,
			expectErr:   false,
			validateErr: false,
		},
		{
			name: "Missing field",
			jsonInput: `{
				"learning_goal_pass_exams": true,
				"learning_goal_refresh_knowledge": true
			}`,
			expectErr:   false,
			validateErr: true, // learning_goal_practice_daily is required
		},
		{
			name: "Null value",
			jsonInput: `{
				"learning_goal_pass_exams": true,
				"learning_goal_refresh_knowledge": true,
				"learning_goal_practice_daily": null
			}`,
			expectErr:   false,
			validateErr: true,
		},
		{
			name: "Invalid type string",
			jsonInput: `{
				"learning_goal_pass_exams": "true",
				"learning_goal_refresh_knowledge": true,
				"learning_goal_practice_daily": true
			}`,
			expectErr:   true, // Unmarshal should fail
			validateErr: false,
		},
		{
			name: "Invalid type number",
			jsonInput: `{
				"learning_goal_pass_exams": 1,
				"learning_goal_refresh_knowledge": true,
				"learning_goal_practice_daily": true
			}`,
			expectErr:   true, // Unmarshal should fail
			validateErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req request.UpdateUser
			err := json.Unmarshal([]byte(tt.jsonInput), &req)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			err = validate.Struct(req)
			if tt.validateErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
