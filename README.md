# GoTH Stack Example

This is a fairly simple template for making Go web apps with session auth.

# Make it your own

This project is set up to provide a ready to use baseline for a simple web server, and still allow changing things to meet your needs.
The first step would be cloning this repo to your local machine, navigating to the root of the new repo, and following along with the next section.

## Customize

> Customize expects to be run at the root of the cloned repository, with a clean working tree.
> So, if you want to use customize, make sure it's run ***BEFORE*** making any changes, since your changes can be overridden in the process.

This process has some dependencies. Make sure these are met before proceeding further:
- The `git` executable should be available on the PATH of the user running customize.
  - Run `git --version` to confirm this.
- The `go` executable should be available on the PATH of the user running customize.
  - The `go` toolchain should be a version compatible with that shown in the `go.mod` of this repo.
  - Run `go version` to get the version and compare with what is shown in [go.mod](go.mod).

## Running customize

This part is pretty simple.
```shell
go run ./cmd/customize
```

This command will ask for a new Go module name, and a new application name.
The Go module name is not validated, but `go vet` is run before making any permanent changes.
The given app name will replace all instances of "yourapp" in this repo.

These are the operations performed, in order, and any errors will result in rolling back the state of the repo to what you initially checked out.
1. Generates template code using Modmake (see "Technologies" below).
2. Changes the module name to what you specified, and retarget import paths to use the new name.
3. Changes paths referencing the old app name to the new name in files.
   - This includes image names in docker-compose.yaml.
4. Changes directories with the old app name to use the new name.
5. Vets and tests code to ensure that it can still be compiled.

At this point customizations are considered verified, and a rollback will not occur.
1. Removes customize since it won't be needed again.
2. Removes the existing `.git` directory and initialize a new repository, creating a commit with the new state as the base commit.

If at this point you're not happy with the result, then delete the repo and clone down the template again, manually changing things as you see fit.

# Structure

Different top level directories have different semantics.
In general, higher directories in the tree (when sorted alphabetically) depend on lower directories.

## `cmd`
This is the "command" directory, and it has a subdirectory for each application with a `main.go`.
Command directories can have internal packages on which they depend, and depend on those in `feature`.

Nothing outside of `cmd` should depend on anything in `cmd`.

### `cmd/<yourapp>`
This is where the entry point of your web server will be (main.go).
Within this directory is:
- `static`
  - This is where all static web assets will be staged and embedded in the built web server.
  - Additional Modmake steps may be added to build front end resources that output into this directory. Just make sure that page templates are updated accordingly.
- `internal/routes`
  - This is where endpoint routes are defined and their dependencies are made available.
  - Route dependencies can be referenced through the `Router` type for simplicity.
- `internal/templates`
  - This is where Templ templates (see "Technologies" below) are created to define whole pages, reusable components, or HTMX responses.
  - There is a `util.go` file in here that can be used to provide helper functionality to these templates. All other generated Go files are ignored by Git.

## `feature`
This directory holds features that directly support the applications in `cmd`.
What goes in here will largely depend on what your app is intended to do.
Packages in this folder should be specific to the domain in which the applications exist, but are more general purpose.
`feature` relies on `foundation`.

This is where auth, audit, and model code is at, each in their own sub-package within `feature`.

## `foundation`
Foundation packages are those that support `feature` or `cmd` packages.
They should be very general purpose, to the point that they could be harvested into libraries if needed in multiple domains.

## `infra`
The files here are used for system component provisioning (with Docker in this case).
There should be a directory for each component that needs this support.

Files here *may* support higher directories, but it's not a requirement.
The main use-case for that would be embedding SQL scripts for live migration.

## `modmake`
Modmake build system files should go here.
This is a flexible build system that uses plain Go to configure tool acquisition, generate code/data, run tests and benchmarks, and build and package your code.

For more information check out the [Modmake documentation](https://saylorsolutions.github.io/modmake/).

# Technologies

- Go (golang)
  - Currently pinned to the latest Go v1.23.
- PostgreSQL
  - The application uses the [pgx/v5](https://github.com/jackc/pgx) driver to interact with the DB.
  - Currently, the DB is only used for auth information, but it can be extended to hold whatever data your app requires.
    - Passwords are hashed with an intentionally slow hashing function, ([crypt](https://www.postgresql.org/docs/current/pgcrypto.html#PGCRYPTO-PASSWORD-HASHING-FUNCS-CRYPT) with [gen_hash](https://www.postgresql.org/docs/current/pgcrypto.html#PGCRYPTO-PASSWORD-HASHING-FUNCS-GEN-SALT)).
      - On my machine (Linux Mint, Intel i7 2.60GHz 6-core CPU with 32GB RAM) it takes about 1-2 seconds to generate/check a password hash.
      - This cost factor can be changed if desired, by modifying the integer value in the procedures in [01_auth.sql](infra/pg/sql/01_auth.sql).
      - The procedure that checks passwords will still run the hashing function if the username doesn't exist, to prevent leaking that information.
    - There are tables for users, sessions, and authorizations, as well as a user flag for setting a user as an "admin".
    - All of this is moldable to your requirements, but the intent is to provide a reasonably secure starting point so you can focus on your business logic.
  - The data model can be extended by:
    - Adding more SQL scripts in `infra/pg/sql/` to represent new data entities. These will be executed when a fresh PG container starts up.
    - Providing model code to interact with the new entities in `feature/model/`.
- [Docker](https://www.docker.com/) with [Compose](https://docs.docker.com/compose/) configuration to build and run your images.
- The [Modmake](https://github.com/saylorsolutions/modmake) build system.
  - This provides enough structure to make build logic easily extensible, while still providing plenty of flexibility.
  - Built binaries are usually output to `build`, which is ignored in the top level `.gitignore` file.
  - Distributable packages of any kind are usually output to `dist`, which is also ignored in the top level `.gitignore` file.
- [Templ](https://templ.guide/) for generating type-safe HTML templates.
  - These templates can be easily extended, and can accept parameters to change what content is generated.
  - To run the application locally, run `go run ./modmake local:run` in the terminal.
  - Note that the `DBURL` set in [local.go](modmake/local.go) should be changed if the DB is not being run on the local machine.
  - The Modmake CLI's watch mode with the step shown above may be useful for regenerating and re-launching the application when changes are detected.
- [HTMX](https://htmx.org/) for client-side functionality.
  - From the HTMX site:
> htmx gives you access to AJAX, CSS Transitions, WebSockets and Server Sent Events directly in HTML, using attributes, so you can build modern user interfaces with the simplicity and power of hypertext
- [Gorilla SecureCookie](https://github.com/gorilla/securecookie)
  - A really mature library for validating (and optionally encrypting) cookie values.
  - This is used in the [auth.Service](feature/auth/auth.go) for handling [CSRF](https://owasp.org/www-community/attacks/csrf) tokens and authenticated [session](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html) cookies.
    - Note that all cookies transmitted have the `HttpOnly` field set to true.
    - This means that the cookie values will not be available to JavaScript, and you'll need to add fields in form templates to transmit the CSRF token value.
    - An example of this can be seen in [login.templ](cmd/yourapp/internal/templates/login.templ) and the [associated router](cmd/yourapp/internal/routes.go).
  - If you want multiple applications (or multiple instances) to recognize the cookie values, then you'll need to set the same hash key in all instances with the `SESSION_HASHKEY` environment variable.
