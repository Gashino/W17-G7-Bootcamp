package models

type Carry struct {
	ID          int    `json:"id"`
	Cid         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  int    `json:"locality_id"`
}

type CarryDoc struct {
	Cid         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  int    `json:"locality_id"`
}

func (v *CarryDoc) AreFieldsValid() bool {
	if v.Cid == "" || v.CompanyName == "" || v.Address == "" ||
		v.Telephone == "" || v.LocalityId == 0 {
		return false
	}
	return true
}
