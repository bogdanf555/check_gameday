# Check Game-day

This program notifies you if your favorite football team (currently hardcoded to Liverpool) plays today.

## Prerequisites

1.  **Tokens**: You need to have an account and API tokens for:
    *   [football-data.org](https://www.football-data.org/)
    *   [Pushover](https://pushover.net/)
2.  **Environment Configuration**: You **MUST** complete this step before running the installation script.
    *   Copy the example environment file:
        ```bash
        cp .env-example .env
        ```
    *   Open `.env` and fill in your tokens:
        - `MY_FOOTBALL_DATA_ORG_AUTH_TOKEN`
        - `MY_PUSHOVER_API_TOKEN`
        - `MY_PUSHOVER_RECIPIENT`

## Installation

Once you have configured the `.env` file, run the installation script:

```bash
chmod +x install.sh
./install.sh
```

The script will:
- Check if Go is installed (and install it if missing).
- Register a cron job to check for matches every day at **9:00 AM**.
- Log any output to `cron.log` in this directory.