# Melkor
Melkor is a sample third party application written on top of the AeroFS API and
SDK. It enumerates list of files, folders and the total number of users on an
AeroFS deployment to showcase the SDK

### Use
1. Register a third party application on your AeroFS Appliance
2. Download the "appconfig.json" for the registered application.
3. Run the following:
```sh
$ make
$ ./melkor <hostname> <port> <appConfigFile>
```
