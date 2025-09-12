# Persona 3: Patrick, the DevOps Engineer

- **Role:** At the crossroads of development and operations, he is the guardian of the software factory.
- **Keywords:** CI/CD, automation, Infrastructure as Code, efficiency, speed.
- **Key Quote:** *"If it's manual, it's broken. Automate everything."*

---

### CLI Expectations
- **JSON Output:** A `--output json` flag is non-negotiable for integration with other tools.
- **Strict Exit Codes:** Must return non-zero exit codes on failure or when an EOL is detected (e.g., via a `--fail-on-eol` flag).
- **Blazing Fast:** Must be lightweight and fast to avoid slowing down the CI/CD pipeline.
- **Dependency-Free:** Must be a single, portable binary.

### Goals
- **Integrate and deploy continuously:** Automate the application lifecycle from build to production.
- **Manage infrastructure as code:** Ensure the reproducibility and reliability of environments.
- **Streamline developer workflow:** Provide tools that accelerate the feedback loop.

### Frustrations
- Tools that are not "API-first" or do not offer structured output (JSON, YAML).
- Manual steps that slow down the deployment pipeline.
- Lack of consistency between dev, staging, and prod environments.

### User Stories
- **As Patrick,** I want to fail a CI/CD pipeline if the project's dependencies are EOL to prevent security risks from reaching production.
- **As Patrick,** I want to parse JSON output from the tool to feed EOL data into our central dashboarding system.
- **As Patrick,** I want to install the CLI in a Docker container with a single command without managing system dependencies.