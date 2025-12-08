package rbac

import (
	"billionmail-core/internal/model/entity"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gogf/gf/v2/frame/g"
)

// ITenantRBAC defines the interface for tenant-aware RBAC operations
type ITenantRBAC interface {
	// Role management
	CreateTenantRole(ctx context.Context, tenantId int64, roleName, displayName, description string, permissions []string) error
	GetTenantRoles(ctx context.Context, tenantId int64) ([]*entity.TenantRole, error)
	UpdateTenantRole(ctx context.Context, tenantRoleId int64, displayName, description string, permissions []string) error
	DeleteTenantRole(ctx context.Context, tenantRoleId int64) error
	
	// User role assignment
	AssignUserRole(ctx context.Context, tenantId, accountId, tenantRoleId int64) error
	RemoveUserRole(ctx context.Context, tenantId, accountId int64) error
	GetUserRole(ctx context.Context, tenantId, accountId int64) (*entity.TenantRole, error)
	GetUserPermissions(ctx context.Context, tenantId, accountId int64) ([]string, error)
	
	// Permission checking
	HasPermission(ctx context.Context, tenantId, accountId int64, permission string) (bool, error)
	CheckAccess(ctx context.Context, tenantId, accountId int64, resource, action string) (bool, error)
	
	// JWT authentication
	GenerateTokens(ctx context.Context, accountId, tenantId int64, userAgent, ipAddress string) (accessToken, refreshToken string, err error)
	ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error)
	RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error)
	RevokeToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, accountId int64) error
}

// JWTClaims defines the JWT claims structure
type JWTClaims struct {
	AccountId   int64    `json:"account_id"`
	TenantId    int64    `json:"tenant_id"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"`
	jwt.RegisteredClaims
}

type sTenantRBAC struct {
	jwtSecret []byte
}

var iTenantRBACService ITenantRBAC

func NewTenantRBAC() *sTenantRBAC {
	return &sTenantRBAC{
		jwtSecret: []byte("r3tain-jwt-secret-key"), // Should be from config
	}
}

func RegisterTenantRBAC(i ITenantRBAC) {
	iTenantRBACService = i
}

func TenantRBAC() ITenantRBAC {
	if iTenantRBACService == nil {
		iTenantRBACService = NewTenantRBAC()
	}
	return iTenantRBACService
}

func init() {
	RegisterTenantRBAC(NewTenantRBAC())
}

// CreateTenantRole creates a new tenant-specific role
func (s *sTenantRBAC) CreateTenantRole(ctx context.Context, tenantId int64, roleName, displayName, description string, permissions []string) error {
	permissionsJSON, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	_, err = g.DB().Model("tenant_role").Data(g.Map{
		"tenant_id":    tenantId,
		"role_name":    roleName,
		"display_name": displayName,
		"description":  description,
		"permissions":  string(permissionsJSON),
		"is_default":   0,
		"status":       1,
		"create_time":  time.Now().Unix(),
		"update_time":  time.Now().Unix(),
	}).Insert()

	return err
}

// GetTenantRoles retrieves all roles for a tenant
func (s *sTenantRBAC) GetTenantRoles(ctx context.Context, tenantId int64) ([]*entity.TenantRole, error) {
	var roles []*entity.TenantRole
	err := g.DB().Model("tenant_role").Where("tenant_id = ? AND status = 1", tenantId).Scan(&roles)
	return roles, err
}

// UpdateTenantRole updates a tenant role
func (s *sTenantRBAC) UpdateTenantRole(ctx context.Context, tenantRoleId int64, displayName, description string, permissions []string) error {
	permissionsJSON, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	_, err = g.DB().Model("tenant_role").Where("tenant_role_id = ?", tenantRoleId).Data(g.Map{
		"display_name": displayName,
		"description":  description,
		"permissions":  string(permissionsJSON),
		"update_time":  time.Now().Unix(),
	}).Update()

	return err
}

// DeleteTenantRole soft deletes a tenant role
func (s *sTenantRBAC) DeleteTenantRole(ctx context.Context, tenantRoleId int64) error {
	_, err := g.DB().Model("tenant_role").Where("tenant_role_id = ?", tenantRoleId).Data(g.Map{
		"status":      0,
		"update_time": time.Now().Unix(),
	}).Update()

	return err
}

// AssignUserRole assigns a role to a user within a tenant
func (s *sTenantRBAC) AssignUserRole(ctx context.Context, tenantId, accountId, tenantRoleId int64) error {
	// Remove existing role assignment
	_, err := g.DB().Model("tenant_user_role").Where("tenant_id = ? AND account_id = ?", tenantId, accountId).Delete()
	if err != nil {
		return err
	}

	// Add new role assignment
	_, err = g.DB().Model("tenant_user_role").Data(g.Map{
		"tenant_id":      tenantId,
		"account_id":     accountId,
		"tenant_role_id": tenantRoleId,
		"create_time":    time.Now().Unix(),
		"update_time":    time.Now().Unix(),
	}).Insert()

	return err
}

// RemoveUserRole removes a user's role within a tenant
func (s *sTenantRBAC) RemoveUserRole(ctx context.Context, tenantId, accountId int64) error {
	_, err := g.DB().Model("tenant_user_role").Where("tenant_id = ? AND account_id = ?", tenantId, accountId).Delete()
	return err
}

// GetUserRole retrieves a user's role within a tenant
func (s *sTenantRBAC) GetUserRole(ctx context.Context, tenantId, accountId int64) (*entity.TenantRole, error) {
	var role entity.TenantRole
	err := g.DB().Model("tenant_role tr").
		LeftJoin("tenant_user_role tur", "tr.tenant_role_id = tur.tenant_role_id").
		Where("tur.tenant_id = ? AND tur.account_id = ? AND tr.status = 1", tenantId, accountId).
		Scan(&role)
	
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetUserPermissions retrieves all permissions for a user within a tenant
func (s *sTenantRBAC) GetUserPermissions(ctx context.Context, tenantId, accountId int64) ([]string, error) {
	role, err := s.GetUserRole(ctx, tenantId, accountId)
	if err != nil {
		return nil, err
	}

	var permissions []string
	if role.Permissions != "" {
		err = json.Unmarshal([]byte(role.Permissions), &permissions)
		if err != nil {
			return nil, err
		}
	}

	return permissions, nil
}

// HasPermission checks if a user has a specific permission within a tenant
func (s *sTenantRBAC) HasPermission(ctx context.Context, tenantId, accountId int64, permission string) (bool, error) {
	permissions, err := s.GetUserPermissions(ctx, tenantId, accountId)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p == permission || p == "*" {
			return true, nil
		}
	}

	return false, nil
}

// CheckAccess checks if a user can perform an action on a resource within a tenant
func (s *sTenantRBAC) CheckAccess(ctx context.Context, tenantId, accountId int64, resource, action string) (bool, error) {
	permission := fmt.Sprintf("%s:%s", resource, action)
	return s.HasPermission(ctx, tenantId, accountId, permission)
}

// GenerateTokens generates access and refresh JWT tokens
func (s *sTenantRBAC) GenerateTokens(ctx context.Context, accountId, tenantId int64, userAgent, ipAddress string) (accessToken, refreshToken string, err error) {
	// Get user permissions
	permissions, err := s.GetUserPermissions(ctx, tenantId, accountId)
	if err != nil {
		return "", "", err
	}

	// Get user role
	role, err := s.GetUserRole(ctx, tenantId, accountId)
	if err != nil {
		return "", "", err
	}

	// Generate access token (15 minutes)
	accessClaims := &JWTClaims{
		AccountId:   accountId,
		TenantId:    tenantId,
		Role:        role.RoleName,
		Permissions: permissions,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "r3tain",
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (7 days)
	refreshClaims := &JWTClaims{
		AccountId: accountId,
		TenantId:  tenantId,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "r3tain",
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Store tokens in database
	accessTokenHash := s.hashToken(accessToken)
	refreshTokenHash := s.hashToken(refreshToken)

	// Store access token
	_, err = g.DB().Model("jwt_token").Data(g.Map{
		"account_id":  accountId,
		"tenant_id":   tenantId,
		"token_hash":  accessTokenHash,
		"token_type":  "access",
		"expires_at":  time.Now().Add(15 * time.Minute).Unix(),
		"is_revoked":  0,
		"user_agent":  userAgent,
		"ip_address":  ipAddress,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}).Insert()
	if err != nil {
		return "", "", err
	}

	// Store refresh token
	_, err = g.DB().Model("jwt_token").Data(g.Map{
		"account_id":  accountId,
		"tenant_id":   tenantId,
		"token_hash":  refreshTokenHash,
		"token_type":  "refresh",
		"expires_at":  time.Now().Add(7 * 24 * time.Hour).Unix(),
		"is_revoked":  0,
		"user_agent":  userAgent,
		"ip_address":  ipAddress,
		"create_time": time.Now().Unix(),
		"update_time": time.Now().Unix(),
	}).Insert()
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *sTenantRBAC) ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Check if token is revoked
		tokenHash := s.hashToken(tokenString)
		count, err := g.DB().Model("jwt_token").Where("token_hash = ? AND is_revoked = 0", tokenHash).Count()
		if err != nil || count == 0 {
			return nil, fmt.Errorf("token is revoked or invalid")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshToken generates new tokens using a refresh token
func (s *sTenantRBAC) RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error) {
	claims, err := s.ValidateToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if claims.TokenType != "refresh" {
		return "", "", fmt.Errorf("invalid token type")
	}

	// Revoke old refresh token
	oldTokenHash := s.hashToken(refreshToken)
	_, err = g.DB().Model("jwt_token").Where("token_hash = ?", oldTokenHash).Data(g.Map{
		"is_revoked":  1,
		"update_time": time.Now().Unix(),
	}).Update()
	if err != nil {
		return "", "", err
	}

	// Generate new tokens
	return s.GenerateTokens(ctx, claims.AccountId, claims.TenantId, "", "")
}

// RevokeToken revokes a specific token
func (s *sTenantRBAC) RevokeToken(ctx context.Context, tokenHash string) error {
	_, err := g.DB().Model("jwt_token").Where("token_hash = ?", tokenHash).Data(g.Map{
		"is_revoked":  1,
		"update_time": time.Now().Unix(),
	}).Update()
	return err
}

// RevokeAllUserTokens revokes all tokens for a user
func (s *sTenantRBAC) RevokeAllUserTokens(ctx context.Context, accountId int64) error {
	_, err := g.DB().Model("jwt_token").Where("account_id = ?", accountId).Data(g.Map{
		"is_revoked":  1,
		"update_time": time.Now().Unix(),
	}).Update()
	return err
}

// Helper function to hash tokens for storage
func (s *sTenantRBAC) hashToken(token string) string {
	hash := make([]byte, 16)
	rand.Read(hash)
	return hex.EncodeToString(hash)
}