# R3tain v1 Architecture Documentation

**Version 1.0** | **Date: October 21, 2025** | 
**Target: Scalable Multi-Tenant Email Marketing Platform (BillionMail Fork)**

---

## 1. Overview

R3tain v1 is a forked and extended version of **BillionMail**, an open-source email marketing and 
mail server solution, transformed into a scalable, self-hosted SaaS platform. 
It supports **multi-tenancy** for agencies and companies, allowing multiple clients to register, 
    manage isolated data, domains, and campaigns while sharing infrastructure securely. 
Key enhancements include **CIDR networking** for private inter-server communication, hybrid tech stack 
    (Go for core logic, TypeScript/Node.js for peripheral services), and tiered plans with IP sharing 
    for cost efficiency.

### Core Objectives
- **Scalability**: Handle 100,000+ tenants (clients/agencies) with millions of emails/day.
- **Multi-Tenancy**: Logical isolation via workspaces, schemas, and RBAC for team collaboration 
    (e.g., adding marketers).
- **Self-Hosted**: Bare metal servers (Hetzner/Cherry) with no cloud dependencies; full control 
    over deliverability.
- **Fork Basis**: Retains BillionMail's containerized mail services (Postfix, Dovecot, Rspamd) 
    and Go core, extended for SaaS features.
- **Compliance**: GDPR-ready with data isolation; CAN-SPAM via unsubscribe tracking.

### Tech Stack Summary
- **Core**: Go (BillionMail's management logic, extended for multi-tenancy).
- **Services**: TypeScript/Node.js (API gateway, auth, queues via Express/BullMQ).
- **Data**: PostgreSQL (sharded), Redis (clustered).
- **Networking**: CIDR-based private subnets.
- **Deployment**: Docker Compose/Swarm on bare metal.

---

## 2. High-Level Architecture

R3tain v1 uses a microservices model with containerized components, mimicking BillionMail's layers 
    but adding multi-tenant isolation and CIDR for scalability.

### Text-Based Diagram
```
Bare Metal Servers (Hetzner/Cherry: 3-12 Servers for 10k-100k+ Tenants)
├── Server Cluster 1 (API + Core) - Public IPs for External Access, Private CIDR: 10.0.1.0/24
│   ├── Docker Network (CIDR: 172.20.0.0/16)
│   │   ├── API Gateway (Node.js/TS) - REST Entry Point, Tenant Routing
│   │   ├── Core Management Service (Go) - Campaign, Template, Analytics Logic
│   │   ├── Auth Service (Node.js/TS) - JWT/RBAC for Teams/Marketers
│   │   └── PostgreSQL Master (Sharded by tenant_id)
├── Server Cluster 2 (Workers + Mail) - Private CIDR: 10.0.2.0/24
│   ├── Docker Network (CIDR: 172.21.0.0/16)
│   │   ├── Queue Workers (Node.js/TS, BullMQ) - Async Email Batching (5-30 Replicas)
│   │   ├── Postfix (SMTP) - Outbound with IP Pools (Shared/Dedicated)
│   │   ├── Dovecot (IMAP) - Inbound Storage
│   │   └── Rspamd (Anti-Spam) - Filtering with Per-Tenant Rules
└── Server Cluster 3 (Storage + Replicas) - Private CIDR: 10.0.3.0/24
    ├── Docker Network (CIDR: 172.22.0.0/16)
    │   ├── Redis Cluster (3 Masters) - Caching, Quotas, Pub/Sub
    │   ├── MinIO (Object Storage) - Tenant-Prefixed Files (Templates, Attachments)
    │   ├── PostgreSQL Replicas - Read Scaling
    │   └── Backup/ELK Stack - Logs, Metrics (Grafana/Prometheus)
```

- **External Access**: API via HTTPS (port 443); SMTP/IMAP for emails.
- **Internal Flow**: API Gateway routes tenant requests to Core; events published to Redis for queues.

---

## 3. Core Components

### 3.1 API Gateway (Node.js/TS)
- Handles REST requests, tenant authentication, and routing to internal services.
- Multi-Tenancy: Extracts tenant_id from JWT; enforces quotas via Redis.

### 3.2 Core Management Service (Go)
- Forked from BillionMail's core; manages campaigns, templates, contacts, and analytics.
- Extensions: Tenant-aware logic (e.g., schema switching); publishes events to Redis.

### 3.3 Auth Service (Node.js/TS)
- Manages JWT, RBAC for roles (admin, marketer, viewer).
- Supports team invites and workspace hierarchies for agencies.

### 3.4 Queue Workers (Node.js/TS)
- Processes batch emails asynchronously; scales horizontally across servers.

### 3.5 Mail Services
- Retained from BillionMail: Postfix (outbound), Dovecot (inbound), Rspamd (filtering).
- Enhancements: Tiered IP pools (shared for free, dedicated for pro).

### 3.6 Data Layer
- **PostgreSQL**: Sharded by tenant_id (Citus extension); schemas for isolation.
- **Redis**: Clustered for caching, quotas, and inter-service communication.
- **MinIO**: Self-hosted S3-like storage for files, with tenant prefixes.

---

## 4. Multi-Tenancy System

### 4.1 Tenant Structure
- **Workspaces**: Agencies/companies as master tenants; sub-tenants for clients/brands.
- **Isolation**: PostgreSQL schemas (e.g., `tenant_123`), prefixed Redis keys, per-tenant DKIM/SPF.
- **Onboarding**: Automated schema creation, domain verification, IP assignment.
- **RBAC**: Admins add marketers with limited roles; enforced via JWT claims.

### 4.2 Tiered Plans
- **Pro/Company**: Dedicated IPs, opt-in multi-tenancy, unlimited team members (10k+ emails/month).
- **Enterprise**: Custom IP pools, full workspaces, priority support (unlimited).

---

## 5. CIDR Networking

### 5.1 Private CIDR Subnets
- **Inter-Server**: 10.0.0.0/16 (e.g., 10.0.1.0/24 for API cluster) via Hetzner vSwitch or Cherry VLAN.
- **Docker Networks**: 172.20.0.0/16 per cluster for container isolation.
- **Benefits**: Secure internal traffic (e.g., DB replication over private links); firewall rules 
    limit exposure.

### 5.2 IP Management
- **Public IPs**: For API (443), SMTP (25/587), IMAP (143/993).
- **Pools**: 10-50 dedicated IPv4 (~$50-250/month); shared for free tiers.

---

## 6. Scaling Strategies

### 6.1 Horizontal Scaling  
- **Services**: Replicate API/queues/core via Docker Swarm (e.g., 2-4 API instances).
- **Data**: PostgreSQL sharding (Citus); Redis clustering (3 masters + replicas).
- **Triggers**: CPU >70% or queue depth >1k; automated via scripts/Prometheus alerts.

### 6.2 Deliverability Scaling
- IP warm-up/rotation; hybrid SES fallback for high-volume tenants.
- Batch processing: Up to 100 emails/API call.

### 6.3 Capacity Projections
- **10k Tenants**: 3 servers (~$250/month).
- **100k Tenants**: 12 servers (~$900/month); 10M+ emails/day.

---

## 7. Deployment Model

### 7.1 Bare Metal Setup
- **Servers**: 3-12 (e.g., Hetzner AX52: 64GB RAM, NVMe storage).
- **Orchestration**: Docker Compose for dev; Swarm/K3s for prod clustering.
- **Deployment Steps**: Provision servers, configure CIDR, clone fork, compose up.

### 7.2 Monitoring
- Prometheus/Grafana for metrics; ELK for logs; per-tenant dashboards.

---

## 8. Security and Compliance

- **Isolation**: Schema RLS, tenant-prefixed storage.
- **Email Security**: DKIM/SPF/DMARC per tenant; Fail2ban for abuse.
- **Compliance**: GDPR data export; audit logs; encrypted backups.

---

## 9. Roadmap for R3tain v1

- **Phase 1 (MVP)**: Fork BillionMail, add basic multi-tenancy/CIDR (4 weeks).
- **Phase 2**: Tiered IPs, RBAC for marketers (4 weeks).
- **Phase 3**: Full scaling, agency optimizations (4 weeks).

This architecture positions R3tain v1 as a robust, self-hosted alternative to managed services, 
    scalable for global agencies while retaining BillionMail's core strengths. For updates, 
    reference your fork repo.