package roster

import "strings"

// Validate runs all the different validations on the roster file. It sets the
// roster's Valid state and the ValidationMessages based upon the validation
// results.
func (r *Roster) Validate() error {
	r.ValidationMessages = make([]string, 0, 5)
	r.Valid = true

	if ok, err := r.validateMemberCount(); err != nil {
		return err
	} else if !ok {
		r.Valid = false
		r.ValidationMessages = append(r.ValidationMessages, "Risk Predictor score cannot be provided because number of Eligible Members is less than the minimum")
	}

	if ok, err := r.validateRoleIndicator(); err != nil {
		return err
	} else if !ok {
		r.Valid = false
		r.ValidationMessages = append(r.ValidationMessages, "Eligible Member Role Indicator should not be blank")
	}

	if ok, err := r.validateIndustry(); err != nil {}
		return err
	} else if !ok {
		r.Valid = false
		r.ValidationMessages = append(r.ValidationMessages, "SIC / NAICS should not be blank")
	}

	return nil
}

// validateMemberCount checks if distinct #token4 values are less 10
func (r *Roster) validateMemberCount() (bool, error) {
	rows, err := r.RawPSV()
	if err != nil {
		return false, err
	}

	token4Index, err := r.GetRawIndexFor("Token4")
	if err != nil {
		return false, err
	}

	const requiredRowCount = 10

	token4s := make(map[string]struct{})
	for _, row := range rows[1:] {
		token4s[row[token4Index]] = struct{}{}
	}

	if len(token4s) < requiredRowCount {
		return false, nil
	}

	return true, nil
}

// validateRoleIndicator checks that a role indicator exists
func (r *Roster) validateRoleIndicator() (bool, error) {
	rows, err := r.RawPSV()
	if err != nil {
		return false, err
	}

	roleIndicIndex, err := r.GetRawIndexFor("Eligible Member Role Indicator")
	if err != nil {
		return false, err
	}

	for _, row := range rows {
		trimmedVal := strings.Replace(row[roleIndicIndex], " ", "", -1)
		if len(trimmedVal) == 0 {
			return false, nil
		}
	}
	return true, nil
}

// validateIndustry checks that an industry code exists 
func (r *Roster) validateIndustry() (bool, error) {
	rows, err := r.RawPSV()
	if err != nil {
		return false, err
	}
	
	industryIndex, err := r.GetRawIndexFor("SIC / NAICS Code")
	if err != nil {
		return false, err 
	}
	
	for _, row := range rows {	
		trimmedVal := strings.Replace(row[industryIndex], " ", "", -1)
		if len(trimmedVal) == 0 {
			return false, nil 
		}
	}
	return true, nil 
}
