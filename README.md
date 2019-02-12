# badboys
**ba**ckup **d**ata**b**ases **o**n **y**our **s**erver

> Sysadmin: \*slaps roof of database\* this bad boy can fit so much fucking data in it!

## Features
- Support for MySQL, PostgreSQL, SQLite3 and MongoDB
- Uses Docker image with the latest database client tools
- Supports databases within Docker containers
- Run command after successful backup

## Installation
```
sudo wget https://get.momar.io/badboys -O /usr/local/bin/
sudo chmod +x /usr/local/bin/badboys
echo "0 * * * * | sudo tee /etc/cron.d/badboys
```

## Usage
1. Create an empty folder somewhere in your project (e.g. `/data/myproject/badboys`)
2. Add that folder to `repositories` in `/etc/badboys.yaml` (defaults to `/var/backup/badboys`, also support globs, e.g. `/data/*/badboys`)
3. Create a file called `databases.yaml` inside that folder
4. Add a line in the `databases` section for each database to back up, e.g. `mydatabase: sqlite3://../data.db`  
   You can use the following schemes:
   - `mysql://[[user][:password]@]server[:port][/database[/table1[/table2[...]]]]`
   - `postgres://[[user][:password]@]server[:port][/database[/table1[/table2[...]]]]``
   - `sqlite3://[path-relative-to-repository]` (currently without Docker support)
   - `mongodb://[connection-uri]`
   For `server`<!-- (or as a path prefix, e.g. `sqlite3://docker=example/data/test.db`)-->, you can also use `docker=containername` to connect directly to a database inside a docker container
5. Run `badboys --dry /data/myproject/badboys` to verify that badboys has access to the database and can create a complete backup of it.

## Using badboys with a file-based backup solution

- Set `filename: backup`. That way, there will always be a single backup copy of the current state of the database.
- Set `oncomplete: your-backup-command` to run the file-based backup solution after every run of badboys

