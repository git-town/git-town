package cmd

import "testing"

func TestAliasPreRunE(t *testing.T) {
	cases := []struct {
		name   string
		args   []string
		expErr bool
	}{
		{"TrueArgument", []string{"true"}, false},
		{"FalseArgument", []string{"false"}, false},
		{"NoArguments", []string{}, true},
		{"MutipleArguments", []string{"true", "false"}, true},
		{"InvalidArgument", []string{"f"}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := aliasCommand.PreRunE(nil, tc.args)
			if tc.expErr && err == nil {
				t.Error("expected PreRunE to return an error")
			}
			if !tc.expErr && err != nil {
				t.Errorf("expeced err to be nil; got: %v", err)
			}
		})
	}
}
