$source = "solid-edu-platform"
$dest = "."

Write-Host "Moving files from $source to $dest..."

Get-ChildItem $source -Exclude node_modules | ForEach-Object {
    $targetPath = Join-Path $dest $_.Name
    if (Test-Path $targetPath) {
        Write-Host "Skipping $($_.Name) - already exists in root"
    } else {
        Move-Item $_.FullName -Destination $targetPath -Force
        Write-Host "Moved: $($_.Name)"
    }
}

Write-Host "`nFiles moved successfully!"
Write-Host "Now you can run: git add ."

