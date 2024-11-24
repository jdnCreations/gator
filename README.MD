# To run this program you will need to have [Postgres](https://www.postgresql.org/download/) and [Go](https://go.dev/doc/install) installed

## Installation

Make sure you have Go 1.23.3 or later installed, then run:
`go install github.com/jdnCreations/gator`
Make sure your Go bin directory is in your PATH to run the `gator` command from anywhere.

## Configuration

Create a config file at `~/.gatorconfig.json` with the following structure:
{
"db_url":<yourdburl>,
"current_user_name":""
}

## Available Commands

- `gator register <username>` - Register a user
- `gator login <username>` - Change user
- `gator addfeed <feedname> <feedurl>` - Add a feed
- `gator follow <feedurl>` - Follow a feed
- `gator unfollow <feedurl>` - Unfollow a feed
- `gator following` - Displays followed feeds
- `gator reset` - Reset users
- `gator feeds` - Lists all feeds
- `gator browse <limit>` - Shows you posts from feeds you follow, with an optional limit (default is 2)
