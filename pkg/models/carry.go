package models

type Carry struct {
	ID          int    `json:"id"`
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  int    `json:"locality_id"`
}

func (v *Carry) AreFieldsValid() bool {
	if v.Cid == 0 || v.CompanyName == "" || v.Address == "" ||
		v.Telephone == "" || v.LocalityId == 0 {
		return false
	}
	return true
}
