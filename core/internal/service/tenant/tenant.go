package tenant

import (
	"billionmail-core/internal/model/entity"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ITenant service interface
type ITenant interface {
	CreateTenant(ctx context.Context, tenant *entity.Tenant) error
	GetTenantBySlug(ctx context.Context, slug string) (*entity.Tenant, error)
	GetTenantById(ctx context.Context, tenantId int64) (*entity.Tenant, error)
	AddUserToTenant(ctx context.Context, tenantId, accountId int64, role string) error
	GetUserTenants(ctx context.Context, accountId int64) ([]*entity.Tenant, error)
	GetUserRoleInTenant(ctx context.Context, accountId, tenantId int64) (string, error)
	CheckTenantAccess(ctx context.Context, accountId, tenantId int64) bool
	GetTenantDatabase(ctx context.Context, tenantId int64) (gdb.DB, error)
}

type sTenant struct{}

var iTenantService ITenant

func New() *sTenant {
	return &sTenant{}
}

func RegisterTenant(i ITenant) {
	iTenantService = i
}

func Tenant() ITenant {
	if iTenantService == nil {
		panic("implement not found for interface ITenant, forgot register?")
	}
	return iTenantService
}

func init() {
	RegisterTenant(New())
}

// CreateTenant creates a new tenant with isolated database schema
func (s *sTenant) CreateTenant(ctx context.Context, tenant *entity.Tenant) error {
	// Generate unique tenant slug if not provided
	if tenant.TenantSlug == "" {
		tenant.TenantSlug = s.generateTenantSlug(tenant.TenantName)
	}

	// Set default values
	if tenant.TenantType == "" {
		tenant.TenantType = "free"
	}
	if tenant.Status == 0 {
		tenant.Status = 1
	}
	if tenant.MaxUsers == 0 {
		tenant.MaxUsers = s.getDefaultMaxUsers(tenant.TenantType)
	}
	if tenant.MaxEmails == 0 {
		tenant.MaxEmails = s.getDefaultMaxEmails(tenant.TenantType)
	}
	if tenant.IpPoolType == "" {
		tenant.IpPoolType = s.getDefaultIpPoolType(tenant.TenantType)
	}

	// Generate database schema name
	tenant.DatabaseName = fmt.Sprintf("tenant_%d", tenant.TenantId)

	// Insert tenant record
	_, err := g.DB().Model("tenant").Ctx(ctx).Insert(tenant)
	if err != nil {
		return err
	}

	// Create tenant-specific database schema
	err = s.createTenantSchema(ctx, tenant.DatabaseName)
	if err != nil {
		return err
	}

	return nil
}

// GetTenantBySlug retrieves tenant by slug
func (s *sTenant) GetTenantBySlug(ctx context.Context, slug string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := g.DB().Model("tenant").Ctx(ctx).Where("tenant_slug", slug).Scan(&tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetTenantById retrieves tenant by ID
func (s *sTenant) GetTenantById(ctx context.Context, tenantId int64) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := g.DB().Model("tenant").Ctx(ctx).Where("tenant_id", tenantId).Scan(&tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// AddUserToTenant adds a user to a tenant with specified role
func (s *sTenant) AddUserToTenant(ctx context.Context, tenantId, accountId int64, role string) error {
	tenantUser := &entity.TenantUser{
		TenantId:   tenantId,
		AccountId:  accountId,
		Role:       role,
		Status:     1,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}

	_, err := g.DB().Model("tenant_user").Ctx(ctx).Insert(tenantUser)
	return err
}

// GetUserTenants retrieves all tenants for a user
func (s *sTenant) GetUserTenants(ctx context.Context, accountId int64) ([]*entity.Tenant, error) {
	var tenants []*entity.Tenant
	err := g.DB().Model("tenant t").
		LeftJoin("tenant_user tu", "t.tenant_id = tu.tenant_id").
		Where("tu.account_id", accountId).
		Where("tu.status", 1).
		Where("t.status", 1).
		Ctx(ctx).
		Scan(&tenants)
	return tenants, err
}

// GetUserRoleInTenant gets user's role within a specific tenant
func (s *sTenant) GetUserRoleInTenant(ctx context.Context, accountId, tenantId int64) (string, error) {
	var role string
	err := g.DB().Model("tenant_user").
		Fields("role").
		Where("account_id", accountId).
		Where("tenant_id", tenantId).
		Where("status", 1).
		Ctx(ctx).
		Scan(&role)
	return role, err
}

// CheckTenantAccess verifies if user has access to tenant
func (s *sTenant) CheckTenantAccess(ctx context.Context, accountId, tenantId int64) bool {
	count, err := g.DB().Model("tenant_user").
		Where("account_id", accountId).
		Where("tenant_id", tenantId).
		Where("status", 1).
		Ctx(ctx).
		Count()
	return err == nil && count > 0
}

// GetTenantDatabase returns database connection for tenant
func (s *sTenant) GetTenantDatabase(ctx context.Context, tenantId int64) (gdb.DB, error) {
	tenant, err := s.GetTenantById(ctx, tenantId)
	if err != nil {
		return nil, err
	}

	// Switch to tenant-specific schema
	db := g.DB().Schema(tenant.DatabaseName)
	return db, nil
}

// Helper functions
func (s *sTenant) generateTenantSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Add timestamp to ensure uniqueness
	return fmt.Sprintf("%s-%d", slug, time.Now().Unix())
}

func (s *sTenant) getDefaultMaxUsers(tenantType string) int {
	switch tenantType {
	case "free":
		return 1
	case "pro":
		return 10
	case "enterprise":
		return -1 // unlimited
	default:
		return 1
	}
}

func (s *sTenant) getDefaultMaxEmails(tenantType string) int64 {
	switch tenantType {
	case "free":
		return 500
	case "pro":
		return 10000
	case "enterprise":
		return -1 // unlimited
	default:
		return 500
	}
}

func (s *sTenant) getDefaultIpPoolType(tenantType string) string {
	switch tenantType {
	case "free":
		return "shared"
	case "pro", "enterprise":
		return "dedicated"
	default:
		return "shared"
	}
}

func (s *sTenant) createTenantSchema(ctx context.Context, schemaName string) error {
	// Create schema for tenant isolation
	_, err := g.DB().Exec(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName))
	if err != nil {
		return err
	}

	// Copy base tables to tenant schema
	tables := []string{
		"batch_mail", "mailbox", "domain", "bm_relay_config",
		"operation_log", "letsencrypt", "alias", "alias_domain",
		"bm_bcc", "bm_relay", "bm_relay_domain_mapping",
	}

	for _, table := range tables {
		createTableSQL := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s.%s (LIKE public.%s INCLUDING ALL);
		`, schemaName, table, table)
		
		_, err = g.DB().Exec(ctx, createTableSQL)
		if err != nil {
			return fmt.Errorf("failed to create table %s in schema %s: %v", table, schemaName, err)
		}
	}

	return nil
}