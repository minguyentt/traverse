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
- implement and design seeds for "mock" testing
- remove the default fallbacks from the cfg envs #NOTE potential full scale project

 * Implement some real-time update socket to monitor database interactions
    - query executions
    - logs/errors
    - tracers
    - listen/notify

[FIX]
