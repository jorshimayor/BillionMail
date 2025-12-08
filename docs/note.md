For the Postfix folder:
Feature,Why you need it,Tool
IP warm-up scheduler,Avoid instant blocks,Custom script or Postal/MailWizz/Mautic feature
FBL (Feedback Loop) registration,Auto-remove complainers,Parse Abuse Feedback Reports from Gmail/Yahoo
List-Unsubscribe + One-Click,Required by Gmail/Yahoo from 2024,Add headers automatically
DKIM signing per domain,Mandatory,OpenDKIM + PostgreSQL lookup or Billion Mail built-in
SPF/DMARC alignment,Mandatory,Your app must set correct From
Bounce & complaint parsing,Keep lists clean,Custom parser → webhook → your DB
Reputation monitoring,Blacklists,Avoid silent blocking