# BillionMail Infrastructure & Backend Overview

This breakdown explains the full infrastructure, backend architecture, and how the API layer works so you can confidently **self-host**, **extend the backend**, or **build your own custom frontend**.

---

# 1. Core Architecture

* **Backend Language:** Go (Golang)
* **Architecture Style:** Microservices, containerized via Docker
* **Backend Pattern:** MVC (Model–View–Controller), visible in the `/core` directory structure
* **Deployment:** Fully orchestrated using `docker-compose`

---

# 2. Infrastructure Components (From `docker-compose.yml`)

## PostgreSQL (v17.4)

* Primary data storage layer
* Runs on Alpine Linux
* Exposed on a configurable port (default: `25432`)
* Persistent storage using Docker volumes

## Redis (v7.4.2)

* Used for:

  * Caching
  * Message queueing
  * Performance optimization for email processing

* Configured with:

  * Elevated `somaxconn` for high connection throughput
  * Password protection
  * Alpine-based Docker image

## Rspamd

* Spam filtering and email processing engine
* Integrates deeply with Redis for optimized performance
* Uses a custom Docker image `billionmail/rspamd:1.2`

---

# 3. Mail Server Components

## Postfix (MTA)

* Handles:

  * SMTP sending
  * SMTP receiving

* Config directory: `postfix/`
* Core responsibility: route outgoing email + accept incoming messages

## Dovecot (MDA)

* Handles:

  * IMAP
  * POP3

* Used for email retrieval (bounces, replies)
* Config directory: `dovecot/`

These two form the full mail server stack that powers BillionMail.

---

# 4. Backend Structure (`/core` Directory)

The backend consists of two major layers: `api` and `internal`.

## API Layer (`/core/api`)

Provides RESTful endpoints for:

* Campaign management
* Contact management
* Email template management
* Batch email operations
* Domain management
* Role-Based Access Control (RBAC)
* Other app-level CRUD operations

This is the part you will use when building your custom frontend.

## Internal Layer (`/core/internal`)

Includes core backend components:

* **Controllers:** Business logic entry points
* **DAO:** Data access layer (PostgreSQL queries)
* **Models:** Database and application data structures
* **Services:** Core business operations (campaign execution, email sending, etc.)

---

# 5. Security Features

* Integrated **Fail2ban** for intrusion prevention
* **SSL/TLS** support (self-signed certificates for dev environments)
* **RBAC** (Role-Based Access Control)
* Password-protected Redis
* Database-level access restrictions
* Email authentication features supported through Postfix & DNS config (SPF/DKIM/DMARC — if implemented)

---

# 6. Email Marketing Features

BillionMail includes a full suite of marketing automation tools:

* Campaign creation & scheduling
* Subscriber list management
* Custom templating system
* Batch email sending (high volume)
* Analytics + event tracking
* Custom embeddable subscription forms
* Bounce processing (via Dovecot IMAP integrations)

---

# 7. Development Tools

* Development helper scripts

  * `run_dev.sh`
  * `go-build.sh`
* Makefile for automated builds
* Fully Dockerized dev environment
* Built-in **multi-language (i18n)** support

---

# 8. Frontend Integration

* Frontend lives in the `/frontend` folder
* Written with **TypeScript**
* Uses **rsbuild** for bundling + development
* Managed via **pnpm**
* Communicates with the Go backend via the API layer

You can replace the default frontend with your own custom UI (React, Vue, Svelte, etc.) by calling the same API endpoints.

---

# 9. Monitoring & Logging

* **rsyslog** handles logging across services
* **Supervisor** manages long-running processes
* Built-in operational logging for audits and debugging

