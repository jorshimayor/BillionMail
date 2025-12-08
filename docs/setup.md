# Postfix

**Postfix** is an open-source Mail Transfer Agent (MTA) designed to send and relay emails. It manages delivery from your application to recipients’ mail servers via SMTP.

### Core Functions

* **Sends emails:** Accepts emails from your app (e.g., campaign emails) and routes them to destination servers.
* **Queue management:** Stores emails if delivery is delayed (e.g., recipient server is down).
* **Relaying:** Forwards emails to other MTAs or directly to recipients’ servers.
* **Security:** Supports TLS encryption and authentication (e.g., SASL).
* **SPAM/DNS handling:** Works with DNS systems (SPF, DKIM, DMARC) to improve deliverability and reduce spam flags.

### In Email Marketing

* Handles **bulk email sending** (e.g., 10,000+ emails).
* Retries failed deliveries for reliability.
* Your backend (Go/Node.js) sends emails through Postfix’s SMTP interface.

### Docker Role

Run Postfix in a container configured with your domain and SMTP credentials to send emails from your app.

---

# Dovecot

**Dovecot** is an open-source IMAP and POP3 server for retrieving and managing incoming emails. It’s useful for handling bounces, replies, and unsubscribe messages.

### Core Functions

* **IMAP/POP3 access:** Allows your app to fetch emails (bounces, replies).
* **Mailbox management:** Stores and organizes incoming emails.
* **Security:** Provides secure, encrypted connections.
* **High performance:** Efficient with large mailboxes and multiple connections.

### In Email Marketing

* **Bounce handling:** Fetch bounce notifications (hard/soft).
* **Replies/unsubscribes:** Reads user replies or unsubscribe emails.
* Integrates with your backend to update subscriber records.

### Docker Role

Run Dovecot in a container to manage mailboxes and allow backend access via IMAP (e.g., using Go `imap` or Node `node-imap`).

---

# How They Fit Together in Your Email Marketing Software

* **Postfix:** Your app sends campaign emails via SMTP → Postfix delivers them and handles retries.
* **Dovecot:** Receives incoming mail (bounces/replies) → your backend fetches them via IMAP to update your database.

### Docker Setup

* Postfix container (SMTP)
* Dovecot container (IMAP/POP3)
* Backend container (Go/TypeScript)
* Frontend container (React/TS)
* Database container (PostgreSQL)

### Example Workflow

1. Your app sends a campaign via Postfix (port 25/587).
2. Postfix delivers emails and logs results.
3. Bounces/replies arrive in a Dovecot-managed mailbox.
4. Your backend fetches them via IMAP, processes them, and updates subscriber records.
5. Your frontend displays analytics.

---

# Why Use Them?

* **Postfix:** Robust, scalable, widely used — excellent for sending thousands of emails reliably.
* **Dovecot:** Efficient for receiving and processing incoming mail, critical for bounce handling.
* Together: A complete, Docker-friendly email server setup that integrates cleanly with Go/TypeScript backends.

---

# Definitions of IMAP, POP3, and SMTP

These protocols handle different aspects of email communication.

---

## SMTP (Simple Mail Transfer Protocol)

### What It Does

Protocol for **sending** emails across the internet — used by mail servers and applications.

### Role in Your Project

Postfix uses SMTP to send bulk campaign emails.
Your backend (Go/TS) connects to Postfix’s SMTP port (25 or 587).

### Key Features

* Sends emails to remote servers
* Supports authentication + TLS
* Manages queues and retries

**Example:** When your app sends a campaign email, SMTP ensures it reaches Gmail, Outlook, etc.

---

## IMAP (Internet Message Access Protocol)

### What It Does

Protocol for **retrieving and managing** emails stored on a server. Keeps emails synced without downloading them permanently.

### Role in Your Project

Dovecot exposes IMAP so your backend can fetch bounce emails, replies, unsubscribes.

### Key Features

* Server-based email syncing
* Folder and selective content retrieval
* Secure via TLS

**Example:** Your backend scans IMAP folders for bounce reports and updates your subscriber list.

---

## POP3 (Post Office Protocol, v3)

### What It Does

Downloads emails from the server — often removes them afterwards. Simpler but less flexible than IMAP.

### Role in Your Project

Supported by Dovecot but rarely used for email marketing. IMAP is preferred.

### Key Features

* Downloads messages
* Simpler, no multi-device sync
* Supports TLS

**Example:** Could retrieve replies if you don’t need server-side storage (not common today).

---

# How They Fit in Your Email Marketing Software

* **SMTP (via Postfix):** Backend sends bulk emails.
* **IMAP (via Dovecot):** Backend fetches bounces/replies.
* **POP3 (via Dovecot):** Optional, rarely used.

---

# In Your Docker Setup

* **Postfix container:** Handles SMTP sending
* **Dovecot container:** Handles IMAP/POP3 access
* **Backend:** Connects via SMTP/IMAP
* **Frontend:** Displays analytics
* **Database:** Stores subscribers, logs, campaigns

---

# Why This Matters

* **SMTP** → reliable campaign delivery
* **IMAP** → essential for bounce/reply processing
* **POP3** → optional fallback
  Together they make your email marketing platform maintainable, scalable, and standards-compliant.

---

