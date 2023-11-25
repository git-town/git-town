# define CLI arguments
param (
  # enables CI mode
  [switch]$ci
)

# stop the script on errors
$ErrorActionPreference = "Stop"

# the Git Town version to release
Set-Variable -Name "GitTownVersion" -Value "v10.0.3" -Option Constant

# dependencies
Set-Variable -Name "GoMsiVersion" -Value "1.0.2" -Option Constant
Set-Variable -Name "GoReleaserVersion" -Value "1.22.1" -Option Constant

Set-Variable -Name "MsiFileName" -Value "git-town_windows_intel_64.msi" -Option Constant

function Main() {
  if ($ci) {
    Install-Tools
  }
  Add-MSI
  .\run-that-app goreleaser@$GoReleaserVersion --clean
}

# generates the .msi file
function Add-MSI() {
  # build the executable that will be inside the .msi file
  go build
  # copy the files needed to build the .msi file on the C: drive to bypass this bug: https://github.com/mh-cbon/go-msi/issues/51
  $tempDir = Join-Path ([System.IO.Path]::GetTempPath()) "git-town"
  if (Test-Path $tempDir) {
    Remove-Item -Path $tempDir -Recurse -Force
  }
  New-Item -Path $tempDir -ItemType Directory
  Copy-Item -Path ".\installer" -Destination $tempDir -Recurse
  Copy-Item -Path ".\LICENSE" -Destination $tempDir
  Copy-Item -Path ".\git-town.exe" -Destination $tempDir
  # change into the temp dir
  $currentDir = Get-Location
  Set-Location -Path $tempDir
  # build the .msi file in the temp dir
  go-msi make --msi $MsiFileName --version $GitTownVersion --src 'installer/templates/' --path 'installer/wix.json'
  # go back to the Git workspace
  Set-Location $currentDir
  # copy the .msi file into the Git workspace
  $msiPath = Join-Path $tempDir $MsiFileName
  Copy-Item -Path $msiPath -Destination $currentDir
  # delete the temp dir
  Remove-Item -Path $tempDir -Recurse -Force
}

# installs the third-party tools needed for the release
function Install-Tools() {
  # install go-msi
  choco install go-msi --version=$GoMsiVersion --no-progress
  # refresh the PATH in this shell instance
  Import-Module $env:ChocolateyInstall\helpers\chocolateyProfile.psm1
  refreshenv
  # add the WiX installation that already exists on CI to the PATH
  $env:PATH = $env:PATH + ";C:\Program Files (x86)\WiX Toolset v3.11\bin"
  # install run-that-app
  Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/kevgo/run-that-app/main/download.ps1" -UseBasicParsing).Content
}

Main
