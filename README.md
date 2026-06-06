# gator

a minimal rss feed aggregator cli.

## prereqs

- go
- postgresql

## setup

1. clone the repo
2. install the cli:

```sh
go install
```

this makes `gator` available in your go bin directory.

3. config:
   create `~/.gatorconfig.json` in your home directory:

```json
{
    "db_url": "postgres://user:pass@localhost:5432/gator?sslmode=disable",
    "current_user_name": ""
}
```

4. database:
   run goose migrations from `sql/schema` to initialize the database.

## usage

run `gator <command> [args...]`

**commands:**

- `register <name>` - create a new user and log in
- `login <name>` - switch active user
- `users` - list all users
- `reset` - wipe the database
- `addfeed <name> <url>` - add a new rss feed and follow it
- `feeds` - list all available feeds
- `follow <url>` - follow an existing feed
- `unfollow <url>` - unfollow a feed
- `following` - list your followed feeds
- `browse [limit]` - view recent posts from followed feeds
- `agg <interval>` - start the background scraper (e.g. `1m`, `1h`)
