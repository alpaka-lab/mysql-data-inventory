# MySQL Data Inventory

This tool scans all user databases in a MySQL cluster and exports the schema inventory (schema, table, column, type) to a CSV file.

## Prerequisites

- Go 1.21 or newer
- Access to a MySQL server
- MySQL user with access to `information_schema`

## Setup

1. Clone this repository.
2. Install dependencies:

   ```sh
   go mod tidy
   ```

3. Edit `main.go` and update the `dsn` variable with your MySQL credentials:

   ```
   dsn := "user:password@tcp(localhost:3306)/information_schema"
   ```

## Usage

Run the tool:

```sh
go run main.go
```

This will generate a file named `mysql_inventory.csv` in the current directory, containing the schema inventory for all non-system databases.

## Output

The CSV file will have the following columns:

- Schema
- Table
- Column
- Type

## Example

```
Schema,Table,Column,Type
mydb,users,id,int
mydb,users,name,varchar
...
```
