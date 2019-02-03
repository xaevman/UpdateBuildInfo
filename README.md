# UpdateBuildInfo
UpdateBuildInfo is a command line utility which rewrites source files to embed branch, build number, and config information into applications during automated builds.

## Example source files
The examples directory contains some example source files for Golang, Python, JSON, C, and C#.

## Variable naming
Build and branch variables must be named branch_id, build_id, branch_name and build_config. Case doesn't matter, and the underscore is optional.

## Example command line
```bash
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BuildInfo.go
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BuildInfo.py
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BuildInfo.json
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BuildInfo.cs
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BuildInfo.h
```
