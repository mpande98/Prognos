package roster

import "strings"

// Validate runs all the different validations on the roster file. It sets the
// roster's Valid state and the ValidationMessages based upon the validation
// results.
func (r *Roster) Validate() error {
	r.ValidationMessages = make([]string, 0, 5)
	r.Valid = true
	
	if ok, err := r.validateID(); err != nil {
		return err
	} else if !ok {
		r.Valid = false
		r.ValidationMessages = append(r.ValidationMessages, "Roster ID Number should not be blank")
	}
	
	if ok, err := r.validateIDUnique(); err != nil {
		return err 
	}else if !ok {
		r.Valid = false
		r.ValidationMessages = append(r.ValidationMessages, "We cannot provide a Risk Predictor Score at this time because it appears that there is more than one Roster ID Number in the file.  Please submit a single unique census roster of Eligible Members in one employer group.  For every Eligible Member in the file, the Roster ID Number must be the same")
	}
	
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
	
	if ok, err := r.validateIndustryUnique(); err != nil {
		return err
	} else if !ok {
		r.Valid = false
		r.ValidationMessages = append(r.ValidationMessages, "We cannot provide a Risk Predictor Score at this time because it appears that there is more than one SIC/NAICS in the file")
	}
	
	if ok, err := r.validateStateCodeUnique(); err != nil {
		return err
	} else if !ok {
		r.Valid = false
		r.ValidationMessages = append(r.ValidationMessages, "We cannot provide a Risk Predictor Score at this time because it appears that there is more than one Employer State Code in the file")
	}

	if ok, err := r.validateZipCodeUnique(); err != nil {
		return err 
	}else if !ok {
		r.Valid = false 
		r.ValidationMessages = append(r.ValidationMessages, "We cannot provide a Risk Predictor Score at this time because it appears that there is more than one Zip Code in the file" )
	}

	return nil
}


// validateID checks if Roster ID Number is present in input roster file 
func (r *Roster) validateID() (bool, error) {
	rows, err := r.RawPSV()

	if err != nil{
		return false, err 
	}

	rosterIDIndex, err := r.GetRawIndexFor("Roster ID Number")
	for _, row := range rows {
		if err != nil {
			return false, err 
		}
		
		trimmedVal := strings.Replace(row[rosterIDIndex]," ", "", -1)
		if len(trimmedVal) == 0 {
			return false, nil 
		}
	}
	return true, nil 
}

// validateIDUnique checks that for every eligible member in the file, the roster ID number is unique
func (r *Roster) validateIDUnique() (bool, error) {
	rows, err := r.RawPSV()
	if err != nil {
		return false, err
	}
		
	rosterIDIndex, err := r.GetRawIndexFor("Roster ID Number")
	if err != nil {
		return false, err 
	}
		
	const uniqueID = 1 //change the constant to be the number of unique eligible members 
		 
	rosterIDs := make(map[string]struct{}) 
	for _, row := range rows[1:]{
		rosterIDs[row[rosterIDIndex]] = struct{}{}
	}
		
		
	if len(rosterIDs) > uniqueID{
		return false, nil 
	}
	return true, nil 
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

// validateZipCodeUique checks that zip codes are unique 
func (r *Roster) validateZipCodeUnique() (bool, error){
	rows, err := r.RawPSV()
	if err != nil {
		return false, err
	}

	zipIndex, err := r.GetRawIndexFor("Employer Zip Code")
	if err != nil {
		return false, err
	}
	
	const uniqueZip = 1 

	zipCode := make(map[string]struct{}) 
	for _, row := range rows[1:]{
		zipCode[row[zipIndex]] = struct{}{}
	}

	if len(zipCode) > uniqueZip{
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

// validateIndustryUnique checks if there are multiple SIC/NAICS codes in one roster 
func (r *Roster) validateIndustryUnique() (bool,error) {
	rows, err := r.RawPSV()
	if err != nil {
		return false, err
	}
	
	industryIndex, err := r.GetRawIndexFor("SIC / NAICS Code")
	if err != nil {
		return false, err
	}

	const uniqueIndustry = 1

	industryCodes := make(map[string]struct{}) 
	for _, row := range rows[1:]{
		industryCodes[row[industryIndex]] = struct{}{}
	}

	if len(industryCodes) > uniqueIndustry{
		return false, nil 
	}	
	
	return true, nil 
}

// validateStateCodeUnique checks that state codes are unique 
func (r *Roster) validateStateCodeUnique() (bool, error){
	rows, err := r.RawPSV()
	if err != nil {
		return false, err
	}

	stateIndex, err := r.GetRawIndexFor("Employer State Code")
	if err != nil {
		return false, err
	}

	const uniqueState = 1

	stateCode := make(map[string]struct{}) 
	for _, row := range rows[1:]{
		stateCode[row[stateIndex]] = struct{}{}
	}

	if len(stateCode) > uniqueState{
		return false, nil 
	}	
	
	return true, nil 
}
