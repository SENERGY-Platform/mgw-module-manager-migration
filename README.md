# mgw-module-manager-migration

One-time migration between two major versions of the [mgw-module-manager](https://github.com/SENERGY-Platform/mgw-module-manager) service.

## Overview

The new major version of the module-manager introduces breaking changes to its persistence layer: the database schema was
redesigned, modules now have a repository source and channel, deployment configuration values are stored in
a typed value model and the manager ID is kept in a JSON file instead of a plain text file.

This program bridges that gap. It embeds a read-only copy of the **old** implementation and reads the existing state
through the old code, transforms it and writes it back into the same database using the new
schema. The old tables are kept as `bk_`-prefixed backups, so the source data remains available for inspection or a
manual rollback.

The migration runs as a Docker container that is provided by the
[mgw-core-installer](https://github.com/SENERGY-Platform/mgw-core-installer) during the core update, before the new
module-manager service starts.

### Requirements

The container needs access to:

- the module-manager MySQL database,
- the container-engine-wrapper (used to resolve container names and volume names of deployed modules),
- the old module and deployment work directories and the old manager ID file,
- the new service directory, to write the manager ID file.

### What is migrated

| Old | New |
| --- | --- |
| Modules | Modules, extended with the configured repository source and channel |
| Deployments | Deployments, with the module version taken from the old deployment record |
| Host resources, secrets, volumes | Deployment host resources, secrets (flattened into secret items), volumes |
| Configs | Deployment user configs, converted into the new typed value model |
| Deployment containers | Deployment containers, identified by container **name** instead of container ID |
| Deployment advertisements | Deployment advertisements |
| Manager ID (plain text file) | Manager ID (JSON file) |

### What is not migrated

- auxiliary deployments
- deployment and module dependency relations (deprecated)
- deployment names (deprecated)

### Constraints and error handling

- Modules whose modfile cannot be read are logged to stderr and skipped.
- Containers that the container-engine-wrapper does not know are logged to stderr; the old container ID is used as the
  container name in that case.
- Config values that do not match the data type declared in the modfile are logged to stderr and skipped.

### Repeated runs

The migration is safe to re-run:

- renaming the old tables is skipped once all `bk_` backup tables exist,
- creating the new tables uses `CREATE TABLE IF NOT EXISTS`,
- modules that already exist in the new schema are skipped,
- all writes of a run happen in a single transaction and are rolled back on error.

## Migration steps

1. **Load configuration** from the optional JSON file given via `-config` and from environment variables, then connect to
   the database.
2. **Initialize the old implementation.** Determine the manager ID and prepare the module and deployment work directories.
3. **Back up the old tables.** All known old tables are renamed to `bk_<table>` and Skipped if they already exist.
4. **Read the old state** from the backup tables: modules (including their modfile), deployments with their assets, 
5. and deployment advertisements. Container names and volume names are resolved through the container-engine-wrapper 
6. filtered by the manager ID label.
5. **Transform** the old state into the new model: modules get a repository source and channel, config values
   are converted into the new typed value model, and map-keyed assets become rows keyed by deployment ID and reference.
6. **Write the new manager ID file** to the new location in JSON form. Skipped if `MANAGER_ID` is set explicitly.
7. **Create the new tables** (`global_configs.sql`, `modules.sql`, `deployments.sql`, `aux_deployments.sql`, 
8. `deployment_advertisements.sql`).
8. **Write the transformed modules** and their deployments and advertisements in one transaction, skipping modules that
   already exist. On error the transaction is rolled back.

The process reacts to `SIGINT`, `SIGTERM` and `SIGQUIT` by cancelling the migration context, which aborts the run and
rolls back the open transaction.

## Configuration

Configuration is read from an optional JSON file and from environment variables; environment variables take precedence
over the file, which takes precedence over the defaults.

```
./bin -config /path/to/config.json
```

Duration values use Go duration syntax (e.g. `30s`, `5m`) as strings in JSON and in environment variables.

| Environment variable               | Config json path                    | Default                                             | Description                                                                                                               |
|------------------------------------|-------------------------------------|-----------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------|
| `MANAGER_ID`                       | `manager_id`                        |                                                     | Manager ID. If set, it overrides the ID read from the old manager ID file and suppresses writing the new manager ID file. |
| `DATABASE_ADDRESS`                 | `database.address`                  |                                                     | Address of the MySQL server, e.g. `mysql:3306`.                                                                           |
| `DATABASE_NAME`                    | `database.database`                 | `module_manager`                                    | Database name.                                                                                                            |
| `DATABASE_USER`                    | `database.user`                     |                                                     | Database user.                                                                                                            |
| `DATABASE_PASSWORD`                | `database.password`                 |                                                     | Database password.                                                                                                        |
| `DATABASE_TIMEOUT`                 | `database.timeout`                  | `30s`                                               | Database connection timeout.                                                                                              |
| `DATABASE_MAX_OPEN_CONNECTIONS`    | `database.max_open_connections`     | `25`                                                | Maximum number of open connections.                                                                                       |
| `DATABASE_MAX_IDLE_CONNECTIONS`    | `database.max_idle_connections`     | `25`                                                | Maximum number of idle connections.                                                                                       |
| `DATABASE_CONNECTION_MAX_LIFETIME` | `database.connection_max_lifetime`  | `5m`                                                | Maximum lifetime of a connection.                                                                                         |
| `OLD_MOD_WORKDIR_PATH`             | `old_impl.mod_handler_workdir_path` | `/opt/module-manager/modules`                       | Module work directory of the old service; modfiles are read from here. Must be absolute.                                  |
| `OLD_DEP_WORKDIR_PATH`             | `old_impl.dep_handler_workdir_path` | `/opt/module-manager/deployments`                   | Deployment work directory of the old service. Must be absolute.                                                           |
| `OLD_MANAGER_ID_PATH`              | `old_impl.manager_id_path`          | `/opt/module-manager/data/mid`                      | Plain text manager ID file of the old service. Only read if `MANAGER_ID` is unset.                                        |
| `CEW_BASE_URL`                     | `old_impl.cew_base_url`             | `http://core-api/ce-wrapper`                        | Base URL of the container-engine-wrapper API, used to look up containers and volumes.                                     |
| `HTTP_TIMEOUT`                     | `old_impl.http_timeout`             | `30s`                                               | Timeout for container-engine-wrapper requests.                                                                            |
| `NEW_MANAGER_ID_PATH`              | `new_impl.manager_id_path`          | `/opt/module-manager/service/mid`                   | JSON manager ID file of the new service; written by the migration.                                                        |
| `REPOSITORY_SOURCE`                | `new_impl.repository_source`        | `github.com/SENERGY-Platform/mgw-module-repository` | Module repository source recorded for every migrated module and deployment.                                               |
| `REPOSITORY_CHANNEL`               | `new_impl.repository_channel`       | `main`                                              | Module repository channel recorded for every migrated module and deployment.                                              |

