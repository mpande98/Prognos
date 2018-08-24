package roster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		rawPSVFile       string
		expectedValid    bool
		expectedMessages []string
	}{
		{"testdata/valid_roster.psv", true, []string{}},
		{"testdata/invalid_roster_validateMemberCount.psv", false, []string{"Risk Predictor score cannot be provided because number of Eligible Members is less than the minimum"}},
		{"testdata/invalid_roster_validateRoleIndicator.psv", false, []string{"Eligible Member Role Indicator should not be blank"}},
		{"testdata/invalid_roster_validateZipCodeUnique.psv", false, []string{"We cannot provide a Risk Predictor Score at this time because it appears that there is more than one Zip Code in the file"}},
	}

	for _, c := range cases {
		r := Roster{}
		r.Load(NewFileLoader(c.rawPSVFile))

		err := r.Validate()
		assert.NoError(t, err)
		assert.Equal(t, c.expectedValid, r.Valid)
		assert.Equal(t, c.expectedMessages, r.ValidationMessages)
	}
}

func TestValidateMemberCount(t *testing.T) {
	cases := []struct {
		rawPSV   string
		expected bool
	}{
		{"Token4|A|B\nt1|X|Y\nt2|X|Y\n", false},
		{"Token4\n1\n2\n3\n4\n5\n6\n7\n8\n9\n0", true},
		{"Token4\n1\n2\n3\n4\n5\n6\n7\n8\n9\n0\nA", true},
		{"Token4\n1\n2\n3\n4\n5\n6\n7\n8\n1\n1\n1", false},
	}

	for _, c := range cases {
		r := Roster{Raw: c.rawPSV}

		actual, err := r.validateMemberCount()
		assert.NoError(t, err)
		assert.Equal(t, c.expected, actual, c.rawPSV)
	}
}

func TestValidateMemberCount_errorToken4Index(t *testing.T) {
	r := Roster{Raw: "A|B\nX|Y\nX|Y\n"}

	_, err := r.validateMemberCount()
	assert.Error(t, err)
}

func TestValidateMemberCount_errorPSV(t *testing.T) {
	r := Roster{Raw: "A|B\nX\n"}

	_, err := r.validateMemberCount()
	assert.Error(t, err)
}

func TestValidateRoleIndicator(t *testing.T) {
	cases := []struct {
		rawPSV   string
		expected bool
	}{
		{"Eligible Member Role Indicator||\n||\n", false},
		{"Eligible Member Role Indicator\n1\n1\n2", true},
		{"Eligible Member Role Indicator\n1", true},
		{"Eligible Member Role Indicator||\n||", false},
	}

	for _, c := range cases {
		r := Roster{Raw: c.rawPSV}

		actual, err := r.validateRoleIndicator()
		assert.NoError(t, err)
		assert.Equal(t, c.expected, actual, c.rawPSV)
	}
}

func TestValidateRoleIndicator_errorRoleIndicIndex(t *testing.T) {
	r := Roster{Raw: "A|B\nX|Y\nX|Y\n"}

	_, err := r.validateRoleIndicator()
	assert.Error(t, err)
}

func TestValidateRoleIndicator_errorPSV(t *testing.T) {
	r := Roster{Raw: "A|B\nX\n"}

	_, err := r.validateRoleIndicator()
	assert.Error(t, err)
}

func TestValidateZipCodeUnique(t *testing.T) {
	cases := []struct {
		rawPSV   string
		expected bool
	}{
		{"Employer Zip Code\n1\n1\n2\n1\n1\n1\n1\n3\n1\n2", false},
		{"Employer Zip Code\n1\n1\n1\n1", true},
		{"Employer Zip Code\n1\n2", false},
		{"Employer Zip Code\n1", true},
	}
 	for _, c := range cases {
		r := Roster{Raw: c.rawPSV}
		
		actual, err := r.validateZipCodeUnique()
		assert.NoError(t, err)
		assert.Equal(t, c.expected, actual, c.rawPSV)
	}
}

func TestValidateZipCodeUnique_errorZipCodeIndex(t *testing.T) {
	r := Roster{Raw: "A|B\nX|Y\nX|Y\n"}
	
	_, err := r.validateZipCodeUnique()
	assert.Error(t, err)
}

func TestValidateZipCodeUnique_errorPSV(t *testing.T) {
	r := Roster{Raw: "A|B\nX\n"}
	
	_, err := r.validateZipCodeUnique()
	assert.Error(t, err)
}
