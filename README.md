# UpdateBuildInfo
UpdateBuildInfo is a command line utility which rewrites source files to embed branch, build number, and config information into applications during automated builds.

## Example source files
The examples directory contains some example source files for Golang, Python, and C#.

## Example command line
```bash
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BulidInfo.go
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BulidInfo.py
UpdateBuildInfo -branchId 1 -buildId 1234 -branchName Main -buildConfig Release -path ./BulidInfo.cs
```