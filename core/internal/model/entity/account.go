package entity

// Account defines the account entity with multi-tenancy support
type Account struct {
	AccountId     int64  `json:"account_id"     dc:"Account ID"`
	Username      string `json:"username"       dc:"Username"`
	Password      string `json:"-"       dc:"Password"`
	Email         string `json:"email"          dc:"Email address"`
	Status        int    `json:"status"         dc:"Status: 1-active, 0-disabled"`
	Language      string `json:"language"       dc:"Language setting"`
	DefaultTenantId int64 `json:"default_tenant_id" dc:"Default tenant ID for this account"`
	AccountType   string `json:"account_type"   dc:"Account type: super_admin, tenant_admin, user"`
	LastLoginTime int64  `json:"last_login_time" dc:"Last login time"`
	CreateTime    int64  `json:"create_time"    dc:"Creation time"`
	UpdateTime    int64  `json:"update_time"    dc:"Update time"`
}
