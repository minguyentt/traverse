# project: Nurse travel contract reviews

## API

### User making a URL request:
* HTTP GET Request => /v1/contracts

// Home page request
// consider pagination
// errors

#### Models

// what the client receives as JSON object
> APIContractsResponse: model
    Contracts `json: "contracts"`

> Contracts: model {
    Name => `json: "name"`
    Location => `json: "location"`
    Agency => `json: "agency"`
    rating => `json: "rating"`
    ReviewCount => `json: "review_count"`
}

> Location: model {
    City => `json: "city"`
    RegionCode => `json: "region_code"`
}

## API flow requests
Handlers => Services => Storage => DB

Things to consider:
- authentication
- authorization
- middlewares
- rate limiting

[NOTES]
- any data being inserted into the db should already be validated from the business logic before executing
- app layer => business/domain/internal/... layer => database storage layer
- id SERIAL PRIMARY KEY, => for auto incrementals when generating user ids

* JWT AUTHENTICATOR
    - we have the admin auth config token settings
    - then we have tokens for sub claims from users

- implement and design seeds for "mock" testing
- remove the default fallbacks from the cfg envs #NOTE potential full scale project
- adapter pattern for business logic and database interactions

[FIX]

- implement a build stage
- running issues with implementing migrations during runtime lol

