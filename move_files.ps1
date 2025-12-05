# Перемещаем все файлы из solid-edu-platform в корень
$source = "solid-edu-platform"
$dest = "."

Get-ChildItem $source -Exclude node_modules | ForEach-Object {
    $targetPath = Join-Path $dest $_.Name
    if (Test-Path $targetPath) {
        Write-Host "Skipping $($_.Name) - already exists"
    } else {
        Move-Item $_.FullName -Destination $targetPath -Force
        Write-Host "Moved: $($_.Name)"
    }
}

# Удаляем пустую папку solid-edu-platform (если осталась только node_modules)
if (Test-Path $source) {
    $remaining = Get-ChildItem $source -Exclude node_modules
    if ($remaining.Count -eq 0) {
        Remove-Item $source -Recurse -Force
        Write-Host "Removed empty folder: $source"
    }
}

Write-Host "`nFiles moved successfully!"

