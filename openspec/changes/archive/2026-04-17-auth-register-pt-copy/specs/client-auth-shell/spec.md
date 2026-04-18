# Delta: client-auth-shell (auth-register-pt-copy)

## ADDED Requirements

### Requirement: Registration visible copy in European Portuguese

All **user-visible** strings on `/auth/register` that are authored in the registration page module (field labels, placeholders, and client-constructed general error text shown in the auth shell banner) SHALL be in **European Portuguese** and SHALL use the **same terminology** as `/auth/login` for equivalent fields (e.g. e-mail field: label and placeholder aligned with login).

#### Scenario: Email field matches login wording

- GIVEN the user opens `/auth/register`
- WHEN they inspect the email field label and placeholder
- THEN the label and placeholder SHALL be in Portuguese and SHALL match the login form’s e-mail label and placeholder semantics (not English marketing phrasing such as “Email Address”).

#### Scenario: Confirm password label in Portuguese

- GIVEN the user opens `/auth/register`
- WHEN they inspect the confirm-password field label
- THEN the label SHALL be in Portuguese describing confirmation of the password (not English “Confirm Password”).

#### Scenario: No English “Error:” prefix from client catch path

- GIVEN the registration submit path catches an unexpected thrown error
- WHEN the general error banner is shown on the registration page
- THEN the banner text SHALL NOT begin with the English literal `Error:` as a client-added prefix.
