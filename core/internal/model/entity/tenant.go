package entity

// Tenant defines the tenant entity for multi-tenancy support
type Tenant struct {
	TenantId     int64  `json:"tenant_id"     dc:"Tenant ID"`
	TenantName   string `json:"tenant_name"   dc:"Tenant name/company name"`
	TenantSlug   string `json:"tenant_slug"   dc:"Unique tenant identifier/slug"`
	TenantType   string `json:"tenant_type"   dc:"Tenant type: free, pro, enterprise"`
	Status       int    `json:"status"        dc:"Status: 1-active, 0-disabled, 2-suspended"`
	MaxUsers     int    `json:"max_users"     dc:"Maximum users allowed"`
	MaxEmails    int64  `json:"max_emails"    dc:"Maximum emails per month"`
	IpPoolType   string `json:"ip_pool_type"  dc:"IP pool type: shared, dedicated"`
	DedicatedIps string `json:"dedicated_ips" dc:"Comma-separated dedicated IP addresses"`
	DatabaseName string `json:"database_name" dc:"Tenant-specific database/schema name"`
	Settings     string `json:"settings"      dc:"JSON settings for tenant configuration"`
	CreateTime   int64  `json:"create_time"   dc:"Creation time"`
	UpdateTime   int64  `json:"update_time"   dc:"Update time"`
}

// TenantUser defines the relationship between tenants and users
type TenantUser struct {
	Id         int64  `json:"id"          dc:"ID"`
	TenantId   int64  `json:"tenant_id"   dc:"Tenant ID"`
	AccountId  int64  `json:"account_id"  dc:"Account ID"`
	Role       string `json:"role"        dc:"User role within tenant: admin, marketer, viewer"`
	Status     int    `json:"status"      dc:"Status: 1-active, 0-disabled"`
	CreateTime int64  `json:"create_time" dc:"Creation time"`
	UpdateTime int64  `json:"update_time" dc:"Update time"`
}

// TenantWorkspace defines workspaces for agency management
type TenantWorkspace struct {
	WorkspaceId   int64  `json:"workspace_id"   dc:"Workspace ID"`
	TenantId      int64  `json:"tenant_id"      dc:"Parent tenant ID"`
	WorkspaceName string `json:"workspace_name" dc:"Workspace name"`
	WorkspaceSlug string `json:"workspace_slug" dc:"Unique workspace identifier"`
	ClientName    string `json:"client_name"    dc:"Client/brand name"`
	Status        int    `json:"status"         dc:"Status: 1-active, 0-disabled"`
	Settings      string `json:"settings"       dc:"JSON settings for workspace"`
	CreateTime    int64  `json:"create_time"    dc:"Creation time"`
	UpdateTime    int64  `json:"update_time"    dc:"Update time"`
}