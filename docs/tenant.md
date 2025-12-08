# Enhancing R3tain Platform for Companies and Tiered Users

## 1. Overview

This document extends your self-hosted BillionMail-based email marketing platform to accommodate **companies** (as tenants) that may want to **add marketers** (team members with limited access) and introduces a **tiered system** for IP sharing. Free/basic users can share IPs for cost efficiency, while pro/enterprise users (including companies) opt into the multi-tenant model for better isolation and deliverability. This draws from Resend-like patterns but keeps everything on-premises, using bare metal servers for control and scalability.

**Key Additions:**
- **Role-Based Access Control (RBAC)** for adding marketers within company tenants.
- **Shared IP Pools** for free users to reduce costs, with dedicated pools for pro users.
- **Optional Multi-Tenancy**: Basic users get simple accounts; pro users enable full isolation if needed.

This ensures your platform supports agencies/companies with team collaboration while maintaining high deliverability for all tiers.

---

## 2. Handling Companies Adding Marketers

### 2.1 RBAC for Team Management
Companies can register as tenants and add marketers (or other roles like admins, viewers) using **Role-Based Access Control (RBAC)** integrated into your tenant model. This allows a company to collaborate without needing full multi-tenancy for every user—marketers operate within the company's tenant boundaries.

**Implementation Strategy:**
- **Registration Flow**: A company admin registers, creating a tenant (e.g., via `/signup/company`). They receive master admin credentials.
- **Adding Marketers**: Use an API endpoint (`/tenants/{tenant_id}/users/add`) where the admin invites marketers by email. The marketer receives a link to set up their account, linked to the tenant.
- **Role Definitions**: Define roles in your database (e.g., PostgreSQL table `roles`):
  - **Admin**: Full access (campaigns, contacts, domains, billing).
  - **Marketer**: Limited to creating/sending campaigns and viewing analytics; no access to billing or user management.
  - **Viewer**: Read-only for reports.
- **Enforcement**: On API calls, validate JWT tokens with role claims (e.g., `{"tenant_id": 123, "role": "marketer"}`). Use RLS in PostgreSQL to restrict queries (e.g., marketers see only their assigned campaigns).
- **Agency Optimization**: Companies can assign marketers to specific sub-clients (if using workspaces); limit team size per tier (e.g., free: 1 user, pro: 10 users).

**Scaling Considerations**: For 100k clients, use Redis to cache role permissions for fast checks; shard user tables by tenant_id to handle large teams.

---

## 3. Tiered IP Sharing System

### 3.1 Shared IPs for Free/Basic Users
Introduce a **shared IP pool** for free or basic users to minimize costs and simplify onboarding, while reserving dedicated resources for pro users. This mirrors common email platforms where free tiers share IPs for low-volume sending, accepting slightly lower deliverability in exchange for accessibility.

**Implementation Strategy:**
- **Pool Setup**: Allocate 2-5 shared IPs (~$10-25/month) in Postfix configs (e.g., `main.cf` with a shared transport map routing low-tier sends to these IPs).
- **Assignment**: On registration, free users are auto-assigned to the shared pool. Monitor collective reputation (e.g., bounces from all free users affect the pool).
- **Warm-Up and Limits**: Start free users at 100 emails/day; use Redis counters for rate limiting (e.g., `tenant:{id}:daily_sends`).
- **Mitigation**: If reputation drops (e.g., >2% bounces), rotate IPs or temporarily pause high-abuse users. Integrate Rspamd for aggressive filtering on shared sends.

**Benefits for Free Users**: No setup needed; ideal for small marketers testing the platform.

### 3.2 Dedicated/Multi-Tenant for Pro/Company Users
Pro users (including companies) can upgrade to **dedicated IP pools** and **opt-in multi-tenancy** for better isolation, higher limits, and agency workflows.

**Implementation Strategy:**
- **Upgrade Flow**: Via `/billing/upgrade` API; on pro activation, assign from a dedicated IP pool (10-20 IPs total, $50-100/month). Update Postfix virtual maps to route their domains to the assigned IP.
- **Opt-In Multi-Tenancy**: Pro users enable workspaces/sub-tenants (e.g., via toggle in dashboard). This creates isolated schemas/pools only when needed, avoiding overhead for simple companies.
- **Pro Perks**: Higher limits (e.g., 100k emails/month); dedicated DKIM/SPF; priority queues in BullMQ for faster processing.
- **Company-Specific**: Companies on pro can add unlimited marketers via RBAC; multi-tenancy allows sub-clients (e.g., brands) with isolated sending.

**Scaling Considerations**: For 100k pro clients, use automated IP rotation scripts (e.g., cron jobs checking reputation via MX Toolbox API integration); shard dedicated pools across servers.

---

## 4. Tiered Plans Overview

```
Tier          | IP Type    | Multi-Tenancy | Team Features     | Monthly Cost | Limits
--------------|------------|---------------|-------------------|--------------|-------
Free/Basic   | Shared     | No            | 1 User (No Add)   | $0           | 500 emails/mo, 1 domain
Pro/Company  | Dedicated  | Opt-In        | Unlimited + RBAC  | $10-50       | 100k+ emails/mo, 10+ domains
Enterprise   | Custom Pool| Yes           | Full Workspaces   | $100+        | Unlimited, Priority Support
```

**Migration Path**: Free users can upgrade seamlessly; on pro, auto-migrate data to dedicated resources.

---

## 5. Scaling for Agencies/Companies

- **Agency Workflow**: Agencies start on pro, add sub-clients via workspaces; assign marketers per sub-client with RBAC.
- **Load Balancing**: Use HAProxy across servers for API/queue traffic; monitor per-tier usage to allocate resources (e.g., more workers for pro queues).
- **Deliverability Safeguards**: Shared pools for free reduce risk to pro IPs; regular audits (e.g., weekly reputation checks) ensure overall platform health.

This tiered approach makes your platform accessible for small users while scalable for companies/agencies, all self-hosted for full control.