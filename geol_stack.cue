// CUE schema to validate the .geol.yml file
// This schema is intended for validating the YAML file used by the stack command in Geol.
// The fields app_name and critical are optional, all other fields are required

geol: {
    // geolVersion: the major geol version this file is compatible with
    geolVersion: string

    // app_name: the name of the report you want as output
    app_name?: string

    stack: [...{
        // name: the name of the product as you want
        // it to appear in the report. For example,
        // use 'Red Hat' if you want a report on RHEL.
        name: string

        // version: the current version of the product.
        // For example, use '3.25' for Quarkus.
        version: string

        // id_eol: the id of the product on
        // endoflife.date (see https://endoflife.date)
        // https://endoflife.date/{id_eol} needs to exist
        id_eol: string

        // critical: causes the geol execution to exit
        // with an error code if the product is in an
        // end of life status, ex: true/false, default is false is not set
        critical?: bool
    }]
}
