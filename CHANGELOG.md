# v1.3.0
- Optimalization - sending mails without authentication to multiple recipients - initialize the connection only once.
- Fix return only the error when sending with authentication
- Upgraded dependencies
- Run staticcheck against go version 1.21.x and 1.22.x

# v1.2.0
- Fixed mail duplicates
- Added missing error messages for `SendMail` and `m.SendMail` functions
  - These functions currently returns multiple errors, which is dependent on recipient count

# v1.1.0
- Added return-path header

# v1.0.0
- Initial implementation