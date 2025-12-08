package database_initialization

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

// Initialize tenant and tenant RBAC related table structures
func init() {
	registerHandler(func() {
		sqlList := []string{
			`-- Tenant table
			CREATE TABLE IF NOT EXISTS tenant (
				tenant_id SERIAL PRIMARY KEY,
				tenant_name VARCHAR(255) NOT NULL,
				tenant_slug VARCHAR(255) NOT NULL UNIQUE,
				tenant_type VARCHAR(50) NOT NULL DEFAULT 'free',
				status INT NOT NULL DEFAULT 1,
				max_users INT NOT NULL DEFAULT 5,
				max_emails BIGINT NOT NULL DEFAULT 100000,
				ip_pool_type VARCHAR(50) NOT NULL DEFAULT 'shared',
				dedicated_ips TEXT NOT NULL DEFAULT '',
				database_name VARCHAR(255) NOT NULL DEFAULT '',
				settings TEXT NOT NULL DEFAULT '{}' ,
				create_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				update_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())
			)`,

			`-- Tenant users table
			CREATE TABLE IF NOT EXISTS tenant_user (
				id SERIAL PRIMARY KEY,
				tenant_id BIGINT NOT NULL,
				account_id BIGINT NOT NULL,
				role VARCHAR(64) NOT NULL DEFAULT 'viewer',
				status INT NOT NULL DEFAULT 1,
				create_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				update_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				UNIQUE(tenant_id, account_id)
			)`,

			`-- Tenant workspaces table
			CREATE TABLE IF NOT EXISTS tenant_workspace (
				workspace_id SERIAL PRIMARY KEY,
				tenant_id BIGINT NOT NULL,
				workspace_name VARCHAR(255) NOT NULL,
				workspace_slug VARCHAR(255) NOT NULL,
				client_name VARCHAR(255) NOT NULL DEFAULT '',
				settings TEXT NOT NULL DEFAULT '{}',
				status INT NOT NULL DEFAULT 1,
				create_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				update_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				UNIQUE(tenant_id, workspace_slug)
			)`,

			`-- Tenant roles table
			CREATE TABLE IF NOT EXISTS tenant_role (
				tenant_role_id SERIAL PRIMARY KEY,
				tenant_id BIGINT NOT NULL,
				role_name VARCHAR(64) NOT NULL,
				display_name VARCHAR(255) NOT NULL DEFAULT '',
				description TEXT NOT NULL DEFAULT '',
				permissions TEXT NOT NULL DEFAULT '[]',
				is_default INT NOT NULL DEFAULT 0,
				status INT NOT NULL DEFAULT 1,
				create_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				update_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				UNIQUE(tenant_id, role_name)
			)`,

			`-- Tenant user-role mapping table
			CREATE TABLE IF NOT EXISTS tenant_user_role (
				id SERIAL PRIMARY KEY,
				tenant_id BIGINT NOT NULL,
				account_id BIGINT NOT NULL,
				tenant_role_id BIGINT NOT NULL,
				create_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				update_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				UNIQUE(tenant_id, account_id)
			)`,

			`-- JWT token storage for tenant-aware authentication
			CREATE TABLE IF NOT EXISTS jwt_token (
				token_id SERIAL PRIMARY KEY,
				account_id BIGINT NOT NULL,
				tenant_id BIGINT NOT NULL,
				token_hash TEXT NOT NULL,
				token_type VARCHAR(20) NOT NULL,
				expires_at BIGINT NOT NULL DEFAULT 0,
				is_revoked SMALLINT NOT NULL DEFAULT 0,
				user_agent TEXT NOT NULL DEFAULT '',
				ip_address TEXT NOT NULL DEFAULT '',
				create_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
				update_time INT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())
			)`,

			`CREATE INDEX IF NOT EXISTS idx_jwt_token_hash ON jwt_token(token_hash)`,
			`CREATE INDEX IF NOT EXISTS idx_jwt_token_account ON jwt_token(account_id, tenant_id)`,
			`CREATE INDEX IF NOT EXISTS idx_jwt_token_type ON jwt_token(token_type)`,
		}

		for _, sql := range sqlList {
			_, err := g.DB().Exec(context.Background(), sql)
			if err != nil {
				g.Log().Error(context.Background(), "Failed to execute tenant RBAC SQL:", err, sql)
				return
			}
		}
	})
}