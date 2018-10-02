# UpdateBuildInfo
UpdateBuildInfo is a command line utility which rewrites source files to embed branch, build number, and config information into applications during automated builds.

## Example source files
The examples directory contains some example source files for Golang, Python, JSON, and C#.

## Variable naming
Build and branch variables must be named BranchId, BuildId, BranchName and BuildConfig. Casing can be all lowercase, pascal case, or camel case. Underscores separating the tokens (ie build_id) are also acceptable.

## Example command line
```bash
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BuildInfo.go
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BuildInfo.py
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BuildInfo.json
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BuildInfo.cs
```
