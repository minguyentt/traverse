# project: Nurse travel contract reviews

#### Models

[TODO]

* System Designs for scaling backend api
    - Ring hashing for redis cache => userid
    - add database indexing for faster lookup performance => faster searches
    - task/message queues for emails? unsure.

[REMINDER]
* For new features:
    1. Model creation
    2. Database queries
    3. service layer
    4. handler layer


job_details
- id
- Profession => "registered nurse"
- pay => ex. "$3,888 to $4000 weekly"
- assignment_length => 13 weeks
- certifications ? => ALCS, BLS, NRP
- experience => "1 year"

