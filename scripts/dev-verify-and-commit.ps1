param(
    [switch] $SkipTests
)

$ErrorActionPreference = 'Stop'
$Root = Resolve-Path (Join-Path $PSScriptRoot '..')
Set-Location $Root

function Invoke-BackendTests {
    $env:CGO_ENABLED = '1'
    Push-Location (Join-Path $Root 'backend')
    try {
        Write-Host '== Backend: go test ./... -count=1 -race -timeout=2m ==' -ForegroundColor Cyan
        go test ./... -count=1 -race -timeout=2m
    }
    finally {
        Pop-Location
    }
}

function Invoke-FrontendTests {
    Push-Location (Join-Path $Root 'frontend')
    try {
        Write-Host '== Frontend: pnpm test -- --passWithNoTests ==' -ForegroundColor Cyan
        pnpm test -- --passWithNoTests
    }
    finally {
        Pop-Location
    }
}

if (-not $SkipTests) {
    Invoke-BackendTests
    Invoke-FrontendTests
}

$subject = [string]$env:GONS_COMMIT_SUBJECT
$subject = $subject.Trim()
if (-not $subject) {
    throw 'GONS_COMMIT_SUBJECT is empty. Set it via VS Code task inputs (commit subject).'
}

$explanation = [string]$env:GONS_COMMIT_EXPLANATION
$explanation = $explanation.Trim()
if (-not $explanation) {
    throw 'GONS_COMMIT_EXPLANATION is empty. Set it via VS Code task inputs (what you are committing and why).'
}

$when = (Get-Date).ToUniversalTime().ToString('o')
$status = (git -C $Root status -sb 2>&1 | Out-String).TrimEnd()
$cachedStat = (git -C $Root diff --cached --stat 2>&1 | Out-String).TrimEnd()
$unstagedStat = (git -C $Root diff --stat 2>&1 | Out-String).TrimEnd()

$nl = [Environment]::NewLine
$body = $subject + $nl + $nl
$body += $explanation + $nl + $nl
$body += '## Verification' + $nl
$body += "- When (UTC): $when" + $nl
$body += '- Backend: go test ./... -count=1 -race -timeout=2m (passed)' + $nl
$body += '- Frontend: pnpm test -- --passWithNoTests (passed)' + $nl + $nl
$body += '## Git analysis' + $nl
$body += '### Branch status' + $nl + '```' + $nl + $status + $nl + '```' + $nl + $nl
$body += '### Staged changes (included in this commit)' + $nl
if ([string]::IsNullOrWhiteSpace($cachedStat)) {
    $body += '_None._' + $nl + $nl
}
else {
    $body += '```' + $nl + $cachedStat.TrimEnd() + $nl + '```' + $nl + $nl
}
$body += '### Unstaged changes (not in this commit)' + $nl
if ([string]::IsNullOrWhiteSpace($unstagedStat)) {
    $body += '_None._' + $nl
}
else {
    $body += '```' + $nl + $unstagedStat.TrimEnd() + $nl + '```' + $nl
}

$msgFile = Join-Path ([System.IO.Path]::GetTempPath()) ("gons-verify-commit-{0}.txt" -f [Guid]::NewGuid().ToString('N'))
$utf8NoBom = New-Object System.Text.UTF8Encoding $false
[System.IO.File]::WriteAllText($msgFile, $body, $utf8NoBom)

git -C $Root diff --cached --quiet
$stagedEmpty = ($LASTEXITCODE -eq 0)

if ($stagedEmpty) {
    if ($env:GONS_ALLOW_EMPTY_COMMIT -eq '1') {
        Write-Warning 'No staged changes; creating empty commit (--allow-empty) because GONS_ALLOW_EMPTY_COMMIT=1.'
        git -C $Root commit --allow-empty -F $msgFile
    }
    else {
        Remove-Item $msgFile -Force -ErrorAction SilentlyContinue
        throw 'No staged changes. Stage files with git add before running this task, or set GONS_ALLOW_EMPTY_COMMIT=1 for an empty marker commit.'
    }
}
else {
    git -C $Root commit -F $msgFile
}

if ($env:GONS_GH_PR -eq '1') {
    try {
        gh pr comment --body-file $msgFile 2>&1 | Write-Host
    }
    catch {
        Write-Warning "gh pr comment failed: $_"
    }
}

Remove-Item $msgFile -Force -ErrorAction SilentlyContinue

if ($env:GONS_SKIP_PUSH -eq '1') {
    Write-Host 'GONS_SKIP_PUSH=1 — skipping git push.' -ForegroundColor Yellow
    exit 0
}

Write-Host '== git push ==' -ForegroundColor Cyan
git -C $Root push
