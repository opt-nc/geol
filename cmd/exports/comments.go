package exports

// Table comments
const (
	CommentTableAbout              = "Metadata about the geol version and platform that generated this database"
	CommentTableDetailsTemp        = "Temporary table for storing raw product release details before type conversion"
	CommentTableDetails            = "Product release lifecycle details including release dates, latest versions, and end-of-life (EOL) dates"
	CommentTableProducts           = "Catalog of all products tracked by geol with their labels, categories, and documentation URIs"
	CommentTableAliases            = "Alternative names or aliases for products to facilitate searching and identification"
	CommentTableProductIdentifiers = "Various identifiers for products such as SKUs, model numbers, or codes used by manufacturers"
	CommentTableTags               = "Tags used to categorize and group products by common characteristics or use cases"
	CommentTableCategories         = "Categories used to group products by their primary function or domain"
	CommentTableProductTags        = "Junction table linking products to their associated tags for many-to-many relationships"
)

// Column comments: about table
const (
	CommentColAboutGitVersion  = "Git tag version of geol used to generate this database"
	CommentColAboutGitCommit   = "Git commit hash of geol used to generate this database"
	CommentColAboutGoVersion   = "Go compiler version used to build geol"
	CommentColAboutPlatform    = "Operating system and architecture where geol was executed"
	CommentColAboutGithubURL   = "GitHub repository URL for the geol project"
	CommentColAboutGeneratedAt = "UTC timestamp when this database was generated"
	CommentColAboutGeneratedTz = "Local timestamp with timezone when this database was generated"
)

// Column comments: details table
const (
	CommentColDetailsProductID         = "Product id referencing the products table"
	CommentColDetailsCycle             = "Product release cycle or version number"
	CommentColDetailsIsLTS             = "Whether this release cycle is a Long Term Support (LTS) version"
	CommentColDetailsReleaseDate       = "Initial release date for this cycle"
	CommentColDetailsLatest            = "Latest patch version within this cycle"
	CommentColDetailsLatestReleaseDate = "Release date of the latest patch version"
	CommentColDetailsEOLDate           = "End-of-life date when this cycle stops receiving support"
)

// Column comments: products table
const (
	CommentColProductsID         = "Unique product id (primary key)"
	CommentColProductsLabel      = "Human-readable display name for the product"
	CommentColProductsCategoryID = "Category id grouping related products"
	CommentColProductsURI        = "URI to the product documentation on endoflife.date"
)

// Column comments: aliases table
const (
	CommentColAliasesID        = "Alternative name or alias for the product (primary key)"
	CommentColAliasesProductID = "Product id referencing the products table"
)

// Column comments: product_identifiers table
const (
	CommentColProductIdentifiersProductID = "Product id referencing the products table"
	CommentColProductIdentifiersType      = "Type of identifier (e.g., repology, purl, cpe)"
	CommentColProductIdentifiersValue     = "Value of the identifier. For repology type, stored with 'repology:' prefix. Reconstruct full URL with: https://repology.org/project/ + identifier_value (after removing 'repology:' prefix)"
)

// Column comments: tags table
const (
	CommentColTagsID  = "Unique tag identifier (primary key)"
	CommentColTagsURI = "URI to the tag page on endoflife.date"
	CommentColTagsWWW = "Human-readable web URL to the tag page on endoflife.date"
)

// Column comments: categories table
const (
	CommentColCategoriesID  = "Unique category identifier (primary key)"
	CommentColCategoriesURI = "URI to the category page on endoflife.date"
)

// Column comments: product_tags table
const (
	CommentColProductTagsProductID = "Product id referencing the products table"
	CommentColProductTagsTagID     = "Tag id referencing the tags table"
)
