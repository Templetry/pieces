# Secrets

This project uses [Doppler](https://docs.doppler.com) as the source of truth for secrets. Values never live in git, chats or wikis; only their names appear in the repository.

## How it fits the project

Gradle reads `app/secrets.properties`, which is gitignored and ships with `changeme` placeholders so a fresh checkout builds. It does not read the process environment, so `doppler run` alone is not enough here. Instead, a script regenerates that file from Doppler:

| Platform | Command |
|---|---|
| macOS / Linux / Git Bash | `./scripts/doppler-secrets.sh` |
| Windows PowerShell | `powershell -ExecutionPolicy Bypass -File scripts/doppler-secrets.ps1` (or `./scripts/doppler-secrets.ps1` if your policy allows scripts) |

Naming: Doppler secrets are `UPPER_SNAKE`. `SIGNING_KEY_ALIAS` and `SIGNING_KEYSTORE_PASSWORD` become `signing_key_alias` and `signing_keystore_password`; `API_BASE_URL_DEVELOPMENT`, `API_BASE_URL_STAGING` and `API_BASE_URL_PRODUCTION` pass through as they are. `DOPPLER_*` metadata is dropped. Any other secret you add is copied verbatim.

The script overwrites `app/secrets.properties`. If Doppler returns nothing it leaves the file untouched and fails.

## One-time setup

1. Install the [Doppler CLI](https://docs.doppler.com/docs/install-cli) and run `doppler login`.
2. Create the project named in `doppler.yaml` (its `dev`, `stg` and `prd` configs are created for you).
3. Bind this checkout with `doppler setup` — it reads `doppler.yaml`, no prompts.
4. Store the secrets without touching shell history; the CLI prompts for the value:

   ```sh
   doppler secrets set SIGNING_KEY_ALIAS
   doppler secrets set SIGNING_KEYSTORE_PASSWORD
   doppler secrets set API_BASE_URL_DEVELOPMENT
   ```

5. Run the script for your platform, then build as usual. Re-run it after rotating a secret.

## CI

The workflows create placeholder secrets so pull requests build without any account. For a signed release, give the job a service token in a `DOPPLER_TOKEN` repository secret and replace the placeholder step:

```yaml
- uses: dopplerhq/cli-action@v4
- run: ./scripts/doppler-secrets.sh
  env:
    DOPPLER_TOKEN: ${{ secrets.DOPPLER_TOKEN }}
```

Keep the placeholder step on pull requests from forks, where secrets are not available.
