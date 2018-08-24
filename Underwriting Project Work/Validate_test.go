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
		{"testdata/invalid_roster_validateID.psv", false, []string{"Roster ID Number should not be blank"}},
		{"testdata/invalid_roster_validateIDUnique.psv", false, []string{"We cannot provide a Risk Predictor Score at this time because it appears that there is more than one Roster ID Number in the file.  Please submit a single unique census roster of Eligible Members in one employer group.  For every Eligible Member in the file, the Roster ID Number must be the same"}},
		{"testdata/invalid_roster_validateMemberCount.psv", false, []string{"Risk Predictor score cannot be provided because number of Eligible Members is less than the minimum"}},
		{"testdata/invalid_roster_validateRoleIndicator.psv", false, []string{"Eligible Member Role Indicator should not be blank"}},
		{"testdata/invalid_roster_validateIndustry.psv", false, []string{"SIC / NAICS should not be blank"}},
		{"testdata/invalid_roster_validateIndustryUnique.psv", false, []string{"We cannot provide a Risk Predictor Score at this time because it appears that there is more than one SIC/NAICS in the file"}},
		{"testdata/invalid_roster_validateStateCode.psv", false, []string{"We cannot provide a Risk Predictor Score at this time because it appears that there is more than one Employer State Code in the file"}},
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

// Unit test: TestValidateID 
func TestValidateID(t *testing.T) {
	cases := []struct {
		rawPSV   string
		expected bool
	}{	
		{"Roster ID Number||\n||\n", false },
		{"Roster ID Number\n1\n1\n2\n1\n1\n1\n1\n3\n1", true},
		{"Roster ID Number|A|B\nt1|X|Y\nt2|X|Y\n", true},
		{"Roster ID Number||\n||\n1||\n", false },
		{"Roster ID Number\n1", true},
	}
 	for _, c := range cases {
		r := Roster{Raw: c.rawPSV}
		
		actual, err := r.validateID()
		assert.NoError(t, err)
		assert.Equal(t, c.expected, actual, c.rawPSV)
	}
}

func TestValidateID_errorIDIndex(t *testing.T) {
	r := Roster{Raw: "A|B\nX|Y\nX|Y\n"}

	_, err := r.validateID()
	assert.Error(t, err)
}

func TestValidateID_errorPSV(t *testing.T) {
	r := Roster{Raw: "A|B\nX\n"}
	
	_, err := r.validateID()
	assert.Error(t, err)
}

// Unit test: TestValidateIDUnique
func TestValidateIDUnique(t *testing.T) {
	cases := []struct {
		rawPSV   string
		expected bool
	}{
		{"Roster ID Number\n1\n1\n2\n1\n1\n1\n1\n3\n1\n2", false },
		{"Roster ID Number\n1\n1\n2\n1\n1\n1\n1\n3\n1", false },
		{"Roster ID Number\n1\n1", true },
		{"Roster ID Number\n1", true },
		{"Roster ID Number\n1\n2", false },
	}
 	for _, c := range cases {
		r := Roster{Raw: c.rawPSV}
		actual, err := r.validateIDUnique()
		assert.NoError(t, err)
		assert.Equal(t, c.expected, actual, c.rawPSV)
	}
}

func TestValidateUniqueID_errorRosterIDIndex(t *testing.T) {
	r := Roster{Raw: "A|B\nX|Y\nX|Y\n"}
	
	_, err := r.validateIDUnique()
	assert.Error(t, err)
}

func TestValidateUniqueID_errorPSV(t *testing.T) {
	r := Roster{Raw: "A|B\nX\n"}
	
	_, err := r.validateIDUnique()
	assert.Error(t, err)
}

// Unit test: TestValidateMemberCount 
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

// Unit test: TestValidateRoleIndicator 
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

// Unit test: TestValidateRoleIndicator 
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

// Unit test: TestValidateIndustry
func TestValidateIndustry(t *testing.T) {
	cases := []struct {
		rawPSV   string
		expected bool
	}{
		{"SIC / NAICS Code||\n||\n", false },
		{"SIC / NAICS Code\n1\n2\n3\n4\n5\n6\n7", true},
		{"SIC / NAICS Code\n1", true},
		{"SIC / NAICS Code||\n||", false },
	}
 	for _, c := range cases {
		r := Roster{Raw: c.rawPSV}
		actual, err := r.validateIndustry()
		assert.NoError(t, err)
		assert.Equal(t, c.expected, actual, c.rawPSV)
	}
}

func TestValidateIndustry_errorIndustryIndex(t *testing.T) {
	r := Roster{Raw: "A|B\nX|Y\nX|Y\n"}
	
	_, err := r.validateIndustry()
	assert.Error(t, err)
}

func TestValidateIndustry_errorPSV(t *testing.T) {
	r := Roster{Raw: "A|B\nX\n"}
	
	_, err := r.validateIndustry()
	assert.Error(t, err)
}

// Unit test: TestValidateIndustryUnique 
func TestValidateIndustryUnique(t *testing.T) {
	cases := []struct {
		rawPSV   string
		expected bool
	}{
		{"SIC / NAICS Code\n1\n1\n2\n1\n1\n1\n1\n3\n1\n2", false},
		{"SIC / NAICS Code\n1\n1\n2\n1\n1\n1\n1\n3\n1", false},
		{"SIC / NAICS Code\n1\n1\n1\n1\n1\n1\n1\n1\n1", true},
		{"SIC / NAICS Code\n1", true},
	}
 	for _, c := range cases {
		r := Roster{Raw: c.rawPSV}
		
		actual, err := r.validateIndustryUnique()
		assert.NoError(t, err)
		assert.Equal(t, c.expected, actual, c.rawPSV)
	}
}

func TestValidateIndustryUnique_errorIndustryIndex(t *testing.T) {
	r := Roster{Raw: "A|B\nX|Y\nX|Y\n"}

	_, err := r.validateIndustryUnique()
	assert.Error(t, err)
}

func TestValidateIndustryUnique_errorPSV(t *testing.T) {
	r := Roster{Raw: "A|B\nX\n"}
	
	_, err := r.validateIndustryUnique()
	assert.Error(t, err)
}

// Unit test: TestValidateStateCodeUnique 
func TestValidateStateCodeUnique(t *testing.T) {
	cases := []struct {
		rawPSV   string
		expected bool
	}{
		{"Employer State Code\n1\n1\n2\n1\n1\n1\n1\n3\n1\n2", false},
		{"Employer State Code\n1\n1\n1\n1", true},
		{"Employer State Code\n1", true},
		{"Employer State Code\n1\n2", false},
	}
 	for _, c := range cases {
		r := Roster{Raw: c.rawPSV}
		actual, err := r.validateStateCodeUnique()
		assert.NoError(t, err)
		assert.Equal(t, c.expected, actual, c.rawPSV)
	}
}

func TestValidateStateCodeUnique_errorStateCodeIndex(t *testing.T) {
	r := Roster{Raw: "A|B\nX|Y\nX|Y\n"}

	_, err := r.validateStateCodeUnique()
	assert.Error(t, err)
}

func TestValidateStateCodeUnique_errorPSV(t *testing.T) {
	r := Roster{Raw: "A|B\nX\n"}

	_, err := r.validateStateCodeUnique()
	assert.Error(t, err)
}

// Unit test: TestValidateZipCodeUnique 
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

