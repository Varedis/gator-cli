# Gator CLI

A command-line RSS feed aggregator and reader built in Go. Gator CLI allows you to manage RSS feeds, follow/unfollow feeds, and browse posts from your followed feeds.

## Description

Gator CLI is a personal RSS feed aggregator that helps you:
- Register and manage user accounts
- Add RSS feeds to the system
- Follow and unfollow feeds
- Browse posts from your followed feeds
- Automatically fetch and aggregate RSS content

The tool uses PostgreSQL as its database backend and stores configuration in a local `.gatorconfig.json` file.

## Prerequisites

- **Go 1.24.5 or later** - Required to build and run the application
- **PostgreSQL database** - Required for data storage
- **Internet connection** - Required to fetch RSS feeds

## Installation

Install Gator CLI using Go's install command:

```bash
go install github.com/varedis/gator-cli@latest
```

After installation, you'll need to set up your configuration file at `~/.gatorconfig.json` with your database connection string:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator_db?sslmode=disable",
  "current_user_name": ""
}
```

## Commands

### User Management

#### `register <name>`
Register a new user account and set it as the current user.

```bash
gator-cli register "John Doe"
```

#### `login <name>`
Switch to an existing user account.

```bash
gator-cli login "John Doe"
```

#### `users`
List all registered users. The current user is marked with "(current)".

```bash
gator-cli users
```

### Feed Management

#### `addfeed <name> <url>`
Add a new RSS feed to the system and automatically follow it. Requires authentication.

```bash
gator-cli addfeed "Tech News" "https://example.com/feed.xml"
```

#### `feeds`
List all RSS feeds in the system.

```bash
gator-cli feeds
```

#### `follow <url>`
Follow an existing RSS feed. Requires authentication.

```bash
gator-cli follow "https://example.com/feed.xml"
```

#### `unfollow <url>`
Unfollow an RSS feed. Requires authentication.

```bash
gator-cli unfollow "https://example.com/feed.xml"
```

#### `following`
List all feeds that the current user is following. Requires authentication.

```bash
gator-cli following
```

### Content Browsing

#### `browse [limit]`
Browse posts from your followed feeds. The optional `limit` parameter defaults to 2 posts. Requires authentication.

```bash
gator-cli browse
gator-cli browse 10
```

### Feed Aggregation

#### `agg <time_between_reqs>`
Start the RSS feed aggregator that continuously fetches feeds at the specified interval. The time should be specified in Go duration format (e.g., "30s", "5m", "1h").

Aggregation is expected to happen in the background, so run this command in a different terminal and you can still browse and do all other actions whilst data is being fetched.

```bash
gator-cli agg 5m
```

### Database Management

#### `reset`
Reset all data in the database. This is a destructive action and will remove users, feeds, follows, posts and all associated content

```bash
gator-cli reset
```

## Examples

### Complete Workflow

1. **Register a new user:**
   ```bash
   gator-cli register alice
   ```

2. **Add and follow a feed:**
   ```bash
   gator-cli addfeed "Tech Crunch" "https://techcrunch.com/feed/"
   ```

3. **Start feed aggregation (in a different terminal):**
   ```bash
   gator-cli agg 10m
   ```

4. **Browse recent posts you are following:**
   ```bash
   gator-cli browse 5
   ```


## Configuration

The application stores configuration in `~/.gatorconfig.json`:

- `db_url`: PostgreSQL connection string
- `current_user_name`: Currently logged-in user (set automatically via `register` and `login` commands)

## Database Schema

The application uses the following main tables:
- `users`: User accounts
- `feeds`: RSS feed definitions
- `feed_follows`: User-feed relationships
- `posts`: Aggregated RSS posts

## Development

To run from source:

```bash
git clone "https://github.com/Varedis/gator-cli.git"
cd gator-cli
go run . <command> [args...]
```
