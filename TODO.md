# GoToSocial TODOs

things noted down by the maintainers that could do with being done!

## database chores
- tidy up db.GetAccountStatuses() (separate functions perhaps?)
- update remaining database queries / API endpoints to use paging.Page{}
- add migrations to convert edits, accounts (, etc) flags to bitsets
- drop unnecessary 'updated_at' columns
- replace 'created_at' columns with time parsed from ULIDs (where possible)
- replace ULID columns with binary representation
- replace "upsert" queries with more performant alternatives (these do a lot of runtime logic which could be determined ahead of time)

## miscellaneous chores
- update to Go 1.25 (we're now being held back from dependency updates)
- finish code commenting where missing (search for '// (\w+\b)?...')
- move away from using Gin, they're all-in on "AI", blegh

## performance
- kim: update ffmpreg to use ncruces/wasm2go (?)
- kim: write alternative bundb serializing backend (?)
- kim: fork bundb to support not always going through "database/sql" (?)
