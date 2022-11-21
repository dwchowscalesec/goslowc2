# goslowc2
A demonstration of using GoLang and GCP Storage for a quick and easy evasion payload.

## Requirements
**DISCLAIMER**: The author is NOT responsible for misuse or abuse of this educational source
// goslowc2 is a demo payload execution for using GCP Storage bucket and objects as the C2 mechanism the time of development Q4 2022 windows defender did not trigger the OS.exec method unless functions were stripped at build time

// For stealthier payload consider using my complement consider using direct WinAPI syscalls converting to unsigned integer pointers

// Author: Dennis Chow dchow[AT]xtecsystems.com 11/20/2022

// Requirements: GCP storage bucket and service account credential foo.json in the same runtime directory

// Confirmed working on GoLang 1.18.4 and GCP Storage SDK 1.28.0
