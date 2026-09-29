# Regenerates app/secrets.properties (gitignored) from the bound Doppler config.
# Doppler names are UPPER_SNAKE, the two signing keys Gradle reads are not, so
# they are mapped back; everything else passes through unchanged.
$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $PSScriptRoot
$out = Join-Path $root 'app/secrets.properties'

$lines = doppler secrets download --no-file --format env-no-quotes
if ($LASTEXITCODE -ne 0) { throw 'doppler secrets download failed' }

$lines = $lines |
    Where-Object { $_ -and $_ -notmatch '^DOPPLER_' } |
    ForEach-Object {
        $_ -replace '^SIGNING_KEY_ALIAS=', 'signing_key_alias=' `
           -replace '^SIGNING_KEYSTORE_PASSWORD=', 'signing_keystore_password='
    }
if (-not $lines) { throw "doppler returned no secrets; $out left untouched" }

# UTF-8 without BOM: Windows PowerShell 5.1's default would prepend one.
[System.IO.File]::WriteAllText($out, (($lines -join "`n") + "`n"), (New-Object System.Text.UTF8Encoding($false)))
Write-Host "wrote $out ($($lines.Count) entries)"
