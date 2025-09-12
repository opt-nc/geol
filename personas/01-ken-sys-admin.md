# Persona 1: Ken, the SysAdmin (System Administrator)

- **Role:** Manages the stability, performance, and security of the infrastructure.
- **Keywords:** Pragmatic, efficiency, automation, reliability.
- **Key Quote:** *"I don't care if it's pretty. Is it scriptable and will it not wake me up at 3 AM?"*

---

### CLI Expectations
- **Reliable Exit Codes:** The tool's success or failure must be trustworthy for scripting.
- **Clear, Concise Output:** Prefers non-verbose, fact-based output.
- **Non-Interactive:** Must be fully controllable via flags and arguments for automation.
- **Fast Execution:** Should not add significant overhead to scripts.

### Goals
- **Automate repetitive tasks:** Deploy security patches, update servers, check configurations.
- **Diagnose and resolve quickly:** Identify the cause of an incident and fix it without wasting time.
- **Maintain consistency:** Ensure that all environments (dev, staging, prod) are configured identically.

### Frustrations
- Tools that require a graphical user interface, which are slow and not scriptable.
- The lack of clear and usable logs for troubleshooting.
- Manual processes that are sources of human error.

### User Stories
- **As Ken,** I want to run a script that checks the OS version of all my servers against the EOL database so I can plan upgrades.
- **As Ken,** I want to get a non-zero exit code when a checked product is EOL so I can trigger an alert in my monitoring system.
- **As Ken,** I want to refresh the local cache via a cron job so my checks are always using recent data without hitting the API every time.
