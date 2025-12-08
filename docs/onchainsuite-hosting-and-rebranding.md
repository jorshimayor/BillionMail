# External Hosting Plan and Rebranding to "onchainsuite"

This document provides a practical plan to host all core APIs on external third‑party servers (initially VPS, later bare‑metal), and to rebrand server and service names to "onchainsuite". It also lists the folders you should study to update API naming and email infrastructure (including DNS records).

## Goals

- Host the core API on external infrastructure with secure, reproducible deployments.
- Phase migration from a single VPS to multi‑node bare‑metal without downtime.
- Rebrand services, server names, and DNS to the `onchainsuite` naming convention.
- Identify exact folders/files to update for API names and email infrastructure.

## Architecture Overview

- Core API: GoFrame-based service in `core/` with controllers, services, and generated API handlers.
- Email stack: Postfix (MTA), Dovecot (IMAP/POP), Rspamd (filtering), Redis; optionally database.
- Deployment options: Docker Compose for single-node (VPS) or multiple hosts; can evolve to container orchestration later.
- Frontend SPA: `core/frontend` serving the UI with environment-specific configs.

### Frontend on Vercel (Next.js)

If you are using a separate Next.js frontend hosted on Vercel, you can ignore the `core/frontend` directory and configure cross-origin access and authentication as follows:

- API base URL: set `NEXT_PUBLIC_API_BASE` in your Vercel project to `https://api.onchainsuite.<your-domain>`.
- CORS: allow origin `https://<your-vercel-domain>` on the core API. You can implement CORS either via GoFrame middleware or a reverse proxy (e.g., Nginx) in front of the API.
- Cookies vs Bearer:
  - If using cookies for auth, set `SameSite=None`, `Secure=true`, and `Domain=.onchainsuite.<your-domain>` so Vercel can send cookies cross-site.
  - Alternatively, use `Authorization: Bearer <token>` headers from Next.js; ensure HTTPS and short-lived access tokens.
- Next.js rewrites (optional to avoid CORS): in `next.config.js`, proxy API routes to your backend:
  ```js
  // next.config.js
  module.exports = {
    async rewrites() {
      return [
        {
          source: '/api/:path*',
          destination: 'https://api.onchainsuite.<your-domain>/:path*',
        },
      ];
    },
  };
  ```
- CSRF: if you use cookie-based sessions, include CSRF tokens in requests; otherwise prefer stateless JWT with Bearer tokens.
- Rate limiting: add per-IP and per-account rate limits at the API; consider an API key header for automation clients.

## Naming Convention: "onchainsuite"

Use a consistent naming scheme across hosts, services, and DNS:

- Hostnames: `onchainsuite-core`, `onchainsuite-mta`, `onchainsuite-imap`, `onchainsuite-filter`, `onchainsuite-redis`, `onchainsuite-db`.
- DNS subdomains: `api.onchainsuite.<your-domain>`, `smtp.onchainsuite.<your-domain>`, `imap.onchainsuite.<your-domain>`, `filter.onchainsuite.<your-domain>`, `webmail.onchainsuite.<your-domain>`.
- Docker Compose service names: `onchainsuite-core`, `onchainsuite-postfix`, `onchainsuite-dovecot`, `onchainsuite-rspamd`.

Replace `<your-domain>` with the domain you control. If you own `onchainsuite.com`, you can drop the subdomain prefix and use direct records (e.g., `api.onchainsuite.com`).

---

## Phase 1: VPS Deployment (Single Node)

- Provision a VPS with at least:
  - CPU: 2–4 vCPU, RAM: 8–16 GB, Disk: 100+ GB SSD
  - OS: Debian/Ubuntu LTS
  - Static public IP and reverse DNS control
- Security baseline:
  - Firewall: allow `22/tcp` (SSH), `80/tcp` and `443/tcp` (HTTP/HTTPS), `25/tcp` (SMTP), `587/tcp` (Submission), `993/tcp` (IMAPS), restrict internal ports
  - TLS: valid certificates for all public endpoints
  - User accounts, sudo, audit; fail2ban enabled (see `Dockerfiles/core/fail2ban.conf`)
- Deployment:
  - Use `docker-compose.yml` for services; rebrand service names in Compose as needed
  - Configure environment with `.env.*` in `core/frontend` and service configs in `conf/`
  - Bind `api.onchainsuite.<your-domain>` to the core API
- Data:
  - Persist volumes for mail queues, user data, database, and logs
  - Backups: daily snapshots of volumes and DB exports

### Rollout Steps

1. Prepare DNS (see checklist below) and obtain TLS certs.
2. Build and run core API (`core/`), then mail services (`Dockerfiles/postfix`, `dovecot`, `rspamd`).
3. Verify API health checks and RBAC login flow.
4. Send test emails to multiple ISPs; validate SPF/DKIM/DMARC.
5. Enable fail2ban and monitor SSH/auth logs.

---

## Phase 2: Bare‑Metal Migration (Multi‑Node)

- Separate concerns across hosts:
  - `onchainsuite-core`: API + frontend
  - `onchainsuite-mta`: Postfix (outbound), Rspamd (filter)
  - `onchainsuite-imap`: Dovecot
  - `onchainsuite-db`: Database (e.g., MySQL/PostgreSQL), Redis on dedicated node
- Networking:
  - Private network/VLAN between nodes; public ingress via reverse proxy/LB
  - Health checks and per‑service firewalls
- Zero‑downtime cutover:
  - Warm up the new MTA IPs; set lower TTLs on DNS records
  - Migrate DB with snapshot and replication; switch traffic via DNS/LB
- Observability:
  - Centralized logging; metrics (CPU, memory, queue depth)
  - Alerting on mail queue backlog, 4xx/5xx API rates, authentication errors

### Rollout Steps

1. Stand up new nodes and configure private networking.
2. Migrate databases and caches; verify replication.
3. Point `api.onchainsuite.<your-domain>` to the new API node; verify.
4. Gradually shift MX and outbound SMTP to new MTA; warm IPs.
5. Disable services on VPS once stable.

---

## Rebranding Checklist: "onchainsuite"

Update naming across configs, code, and deployment manifests:

- Hostnames in OS and Docker Compose service names
- DNS records for API and mail
- Application environment variables in `.env.*` and config files
- Any hardcoded banners or `hostname` fields in Postfix/Dovecot/Rspamd
- Documentation, README, and user-facing messages

Search/replace patterns to audit:

- `hostname`, `myhostname`, `mydomain`, `smtpd_banner`, `smtp_helo_name`
- Rspamd `bind_hostname`, `hostname`, worker names
- Dovecot `auth_mechanisms`, service names, syslog identifiers
- Frontend env vars and API base URLs

---

## DNS and Email Infrastructure Checklist

Configure DNS for proper mail delivery and API access under the `onchainsuite` naming scheme.

- A/AAAA records
  - `api.onchainsuite.<your-domain>` → core API public IP
  - `smtp.onchainsuite.<your-domain>` → Postfix public IP
  - `imap.onchainsuite.<your-domain>` → Dovecot public IP
  - `filter.onchainsuite.<your-domain>` → Rspamd public IP
  - `webmail.onchainsuite.<your-domain>` → webmail app IP (if used)
- MX records
  - `MX onchainsuite.<your-domain>` → `smtp.onchainsuite.<your-domain>` with appropriate priority
- SPF (TXT)
  - `v=spf1 a mx include:spf.<your-sending-provider> ~all` (adjust for your sending pattern)
- DKIM (TXT/CNAME)
  - Publish selector TXT at `selector._domainkey.onchainsuite.<your-domain>`; ensure Postfix signs outbound mail
- DMARC (TXT)
  - `_dmarc.onchainsuite.<your-domain>` with policy: `v=DMARC1; p=none|quarantine|reject; rua=mailto:dmarc@onchainsuite.<your-domain>`
- PTR (Reverse DNS)
  - Set rDNS of MTA IP to `smtp.onchainsuite.<your-domain>`
- TLS
  - Certificates for all public endpoints (`api`, `smtp`, `imap`, `filter`, `webmail`)
- MTA ports
  - `25/tcp` (SMTP), `587/tcp` (submission), `465/tcp` (if used), ensure STARTTLS
- IMAP/POP ports
  - `993/tcp` (IMAPS), `995/tcp` (POP3S if used)

---

## Folders to Study and Update (API)

- `core/api/`
  - Generated API definitions and request/response models. Do not hand-edit generated files; update controllers/services instead.
- `core/internal/controller/`
  - HTTP handlers. Update any user-facing messages or URLs referencing old names.
- `core/internal/service/`
  - Business logic, RBAC, JWT, tenant features. Review for any hostname or domain assumptions.
- `core/internal/cmd/`
  - Application initialization. `InitDatabase()` is wired here; ensure boot-time configs align with the new environment.
- `core/manifest/deploy/` and `core/manifest/docker/`
  - Deployment manifests and scripts. Adjust service names and images to `onchainsuite-*`.
- `core/Makefile` and `core/go-build.sh`
  - Build targets. Align naming conventions and CI/CD tasks.
- `core/main.go` and `core/run_dev.go`
  - Entry points and dev runner. Confirm ports and bind addresses.
- External frontend (Vercel/Next.js)
  - Skip `core/frontend/*`. Instead, set `NEXT_PUBLIC_API_BASE` in Vercel to `https://api.onchainsuite.<your-domain>` and configure CORS as described above.

## Folders to Study and Update (Email Infrastructure)

- `Dockerfiles/postfix/`
  - Postfix image and scripts. Update `postfix.sh`, service names, and any hostname references.
- `conf/postfix/`
  - `main.cf`, `master.cf`, `rsyslog.conf` — set `myhostname = smtp.onchainsuite.<your-domain>`, `mydomain`, `smtpd_banner`, HELO policies.
- `Dockerfiles/dovecot/` and `conf/dovecot/`
  - Dovecot image and configs. Ensure correct `imap` hostnames and TLS.
- `Dockerfiles/rspamd/` and `conf/rspamd/`
  - Rspamd image and configs under `local.d/` and `rspamd.conf`. Update `hostname` and WebUI bind settings.
- `conf/redis/`
  - Redis config script. Check bind addresses and passwords.
- `docker-compose.yml`
  - Rebrand service names, volumes, and environment variables; split per-node or per-stack as needed.
- `env_init` and `init.sql`
  - Initial environment and SQL bootstrap; ensure any default domains or banners match `onchainsuite` branding.
- `conf/webmail/`
  - Webmail settings and customizations if you ship webmail.

---

## Configuration Examples (Placeholders)

Replace placeholders with your actual domain and selectors.

- Postfix (`conf/postfix/main.cf`):
  - `myhostname = smtp.onchainsuite.<your-domain>`
  - `mydomain = onchainsuite.<your-domain>`
  - `smtpd_banner = onchainsuite MTA ready`
- Dovecot (`conf/dovecot/dovecot.conf`):
  - Ensure `ssl_cert`, `ssl_key`, listeners on `imap.onchainsuite.<your-domain>`
- Rspamd (`conf/rspamd/local.d/options.inc` or similar):
  - `hostname = "filter.onchainsuite.<your-domain>"`
- Frontend (`core/frontend/.env.production`):
  - `VITE_API_BASE=https://api.onchainsuite.<your-domain>`

---

## Testing and Verification

- Core API
  - Build: `cd core && go mod tidy && go build ./internal/...`
  - Health: call `/health` or equivalent; run auth/RBAC flows
- Mail Delivery
  - Send to Gmail/Outlook; check headers for SPF/DKIM pass and DMARC alignment
  - Verify rDNS and HELO match `smtp.onchainsuite.<your-domain>`
- Security
  - TLS scans on `api`, `smtp`, `imap` endpoints
  - Fail2ban and firewall rules enforced

---

## Migration Runbook (High-Level)

1. Prepare VPS with Docker, firewall, users, and TLS.
2. Configure DNS records and certificates for `onchainsuite` endpoints.
3. Deploy core API and verify with test tenants/users.
4. Deploy mail stack components; validate outbound and inbound flows.
5. Monitor logs, queues, and metrics; adjust reputational warmup.
6. For bare‑metal, provision nodes, migrate DB/caches, then cut over traffic.

---

## Notes and Recommendations

- Keep DNS TTL low during migration to allow rapid changes.
- Warm up new MTA IPs gradually to avoid deliverability issues.
- Keep `init.sql` and `database_initialization` handlers aligned with tenant RBAC changes.
- Avoid editing generated files (e.g., `core/api/**` headers indicate generated code); make changes in controllers/services.
- Document all renames in `README.md` or this doc to track the rebranding.

---

## Appendix: Quick Search Terms

Use these queries to find places needing updates:

- API: `core/internal/controller`, `core/internal/service`, `core/manifest/**`, `core/frontend/.env.*`
- Email: `conf/postfix/*`, `conf/dovecot/*`, `conf/rspamd/*`, `Dockerfiles/(postfix|dovecot|rspamd)/*`, `docker-compose.yml`
- General: search for `hostname`, `domain`, `smtpd_banner`, `VITE_API_BASE`, `onchain`/old brand names
