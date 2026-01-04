// CUE schema to validate the .geol.yml file
// This schema is intended for validating the YAML file used by the stack command in Geol.
// The fields app_name and app_id are optional, all other fields are required

// geolVersion: the major geol version this file is compatible with
geolVersion: string

// app_name: the name of the report you want as output
app_name?: string

// app_id: an optional identifier for the application
app_id?: string

stack: [...{
    // name: the name of the product as you want
    // it to appear in the report. For example,
    // use 'Red Hat' if you want a report on RHEL.
    // IMPORTANT: The name field must be UNIQUE across all stack items.
    // If you need to track the same product in different environments,
    // use suffixes like 'traefik-prod' and 'traefik-qual'.
    name: string

    // version: the current version of the product.
    // For example, use '3.25' for Quarkus.
    version: string

    // id_eol: the id of the product on
    // endoflife.date (see https://endoflife.date)
    // https://endoflife.date/{id_eol} needs to exist
    id_eol: string

    // skip: whether to skip this product in EOL checks.
    // Optional field. Defaults to false if not specified.
    // Set to true to exclude this product from checks
    // (e.g., when not yet available in endoflife.date API).
    // Products with skip: true remain in the YAML for tracking
    // but are not checked against the endoflife.date API.
    skip?: bool | *false
}]
