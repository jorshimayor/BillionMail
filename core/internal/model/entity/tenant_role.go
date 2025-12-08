package entity

// TenantRole defines tenant-specific roles for the R3tain model
type TenantRole struct {
	TenantRoleId int64  `json:"tenant_role_id" dc:"Tenant Role ID"`
	TenantId     int64  `json:"tenant_id"      dc:"Tenant ID"`
	RoleName     string `json:"role_name"      dc:"Role name: admin, marketer, viewer"`
	DisplayName  string `json:"display_name"   dc:"Display name for the role"`
	Description  string `json:"description"    dc:"Role description"`
	Permissions  string `json:"permissions"    dc:"JSON array of permissions"`
	IsDefault    int    `json:"is_default"     dc:"Is default role: 1-yes, 0-no"`
	Status       int    `json:"status"         dc:"Status: 1-active, 0-disabled"`
	CreateTime   int64  `json:"create_time"    dc:"Creation time"`
	UpdateTime   int64  `json:"update_time"    dc:"Update time"`
}

// TenantUserRole defines the relationship between tenant users and their roles
type TenantUserRole struct {
	Id           int64 `json:"id"             dc:"Relation ID"`
	TenantId     int64 `json:"tenant_id"      dc:"Tenant ID"`
	AccountId    int64 `json:"account_id"     dc:"Account ID"`
	TenantRoleId int64 `json:"tenant_role_id" dc:"Tenant Role ID"`
	CreateTime   int64 `json:"create_time"    dc:"Creation time"`
	UpdateTime   int64 `json:"update_time"    dc:"Update time"`
}

// JWTToken defines JWT token storage for authentication
type JWTToken struct {
	TokenId     int64  `json:"token_id"     dc:"Token ID"`
	AccountId   int64  `json:"account_id"   dc:"Account ID"`
	TenantId    int64  `json:"tenant_id"    dc:"Tenant ID"`
	TokenHash   string `json:"token_hash"   dc:"Token hash"`
	TokenType   string `json:"token_type"   dc:"Token type: access, refresh"`
	ExpiresAt   int64  `json:"expires_at"   dc:"Expiration timestamp"`
	IsRevoked   int    `json:"is_revoked"   dc:"Is revoked: 1-yes, 0-no"`
	UserAgent   string `json:"user_agent"   dc:"User agent"`
	IpAddress   string `json:"ip_address"   dc:"IP address"`
	CreateTime  int64  `json:"create_time"  dc:"Creation time"`
	UpdateTime  int64  `json:"update_time"  dc:"Update time"`
}