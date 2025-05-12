package models

type Inquiry struct {
	ID                       uint    `gorm:"primaryKey;autoIncrement;not null"`
	PartnerReferenceNo       string  `json:"partner_reference_no" gorm:"type:varchar(255);unique"`
	PartnerID                string  `json:"partner_id" gorm:"type:varchar(255)"`
	BeneficiaryAccountNumber string  `json:"beneficiary_account_number" gorm:"type:varchar(255)"`
	AmountValue              float64 `json:"amount_value,string" gorm:"type:double"`
	AmountCurrency           string  `json:"amount_currency" gorm:"type:varchar(255)"`
	BeneficiaryBankCode      string  `json:"beneficiary_bank_code" gorm:"type:varchar(255)"`
}

func (Inquiry) TableName() string {
	return "inquiries"
}
