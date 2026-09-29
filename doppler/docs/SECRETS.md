# Secrets

This project uses [Doppler](https://docs.doppler.com) as the source of truth for secrets. Values never live in git, chats or wikis; only their names appear in the repository.

## How it fits the project

The service reads its configuration from the environment, and **a real environment variable always beats the `.env.<profile>` files**. `doppler run` injects secrets as environment variables, so no code changes are involved.

`.env.<profile>` files stay for non-secret defaults; `.env.local` (gitignored) stays available for anything machine-specific. Doppler is the third layer, not a replacement.

## One-time setup

1. Install the [Doppler CLI](https://docs.doppler.com/docs/install-cli) and run `doppler login`.
2. Create the project named in `doppler.yaml` (its `dev`, `stg` and `prd` configs are created for you).
3. Bind this checkout with `doppler setup` — it reads `doppler.yaml`, no prompts.
4. Store secrets without touching shell history; the CLI prompts for the value:

   ```sh
   doppler secrets set SOME_API_KEY
   ```

## Daily use

Prefix whatever you normally run:

```sh
doppler run -- <your run command>
doppler run -- <your test command>
```

To use another config: `doppler run --config stg -- <command>`.

## CI

Prefer Doppler's **GitHub Actions sync integration** (Doppler dashboard → project → Integrations): it keeps repository secrets in sync on every rotation and the workflows keep reading `secrets.NAME` — nothing to change here.

If you would rather fetch at run time, give the job a service token in a `DOPPLER_TOKEN` repository secret and wrap the step:

```yaml
- uses: dopplerhq/cli-action@v4
- run: doppler run -- <your command>
  env:
    DOPPLER_TOKEN: ${{ secrets.DOPPLER_TOKEN }}
```

## Containers

Do not bake secrets into the image. Set `DOPPLER_TOKEN` (a service token) in the runtime and start the process through `doppler run`, or inject the variables with your orchestrator.

## Rotating

`doppler secrets set NAME` → syncs propagate it. Doppler keeps per-secret version history and access logs.
