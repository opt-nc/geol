# Persona 2: Grace, the CISO (Chief Information Security Officer)

- **Role:** Responsible for risk management, compliance, and data protection.
- **Keywords:** Risk, compliance, audit, traceability, control.
- **Key Quote:** *"I need proof, not opinions. Show me the data that proves we are compliant."*

---

### CLI Expectations
- **Verifiable Data Source:** The origin of the information must be clear and reputable.
- **Exportable Data:** Needs output in standard, machine-readable formats (CSV, JSON) for auditing and reporting.
- **Timestamped Information:** Must be able to prove the freshness of the data (e.g., when the cache was last updated).

### Goals
- **Ensure compliance:** Make sure all software configurations and dependencies comply with the company's security policy.
- **Have complete traceability:** Know who deployed what and when, and with which components.
- **Generate reports:** Provide proof of compliance to auditors and management.

### Frustrations
- "Black boxes": systems and tools without clear audit logs.
- Tools that do not allow for granular security policy definitions.
- The difficulty of getting a comprehensive view of risk exposure.

### User Stories
- **As Grace,** I want to generate a monthly report of all obsolete products to present to the steering committee.
- **As Grace,** I want to be notified immediately if a production deployment contains a critical end-of-life dependency.
- **As Grace,** I want to see when the EOL data was last updated to ensure our reports are based on current information.
