# `fbackup`, because it's never late to make backups

`fbackup` is a simple tool to take a directory in your filesystem, a Google Drive location and upload every file from the local folder into Google Drive. If the local file modification date is newer than the one reported by Google Drive, the file will be overwritten. Extraneous files (files that exist in Google Drive but not locally) **will be deleted from Google Drive** and the preference is always to persist local files over remote files anytime.

## Usage

The project is meant to be used with Docker in mind, so settings are provided via environment variables:

* `FBACKUP_QUIET` will make the program not output a single thing to `stdout`.
* `FBACKUP_ACCOUNT_FILE` a path to the local `service-account.json` file. You can [create service accounts here](https://console.cloud.google.com/iam-admin/serviceaccounts). Don't specify any permission in particular in the "Grant this service account access to the project" and simply create it. More setup below.
* `FBACKUP_FOLDER` the target folder you want to sync to Google Drive (say, `~/Projects`).
* `FBACKUP_DESTINATION_ID` the ID of the Google Drive folder where you want to store the folder you're uploading / synchronizing. To get an ID, simply go to `drive.google.com`, navigate to the folder you want to use by double clicking on it, then in the URL you'll see something like this: `https://drive.google.com/drive/u/1/folders/1z-xSDBGxR6WzktgIca4HdiipQP9sJ9lG`, note, at the end, that there's a random ID, `1z-xSDBGxR6WzktgIca4HdiipQP9sJ9lG`. That's the destination ID you'll be using to synchronize. Keep in mind that your service account **won't see this folder** unless you explicitly share it with the account. To do this, open the account file, find the e-mail address associated with it in the `client_email` field from the JSON file (it will look like `my-service-account@my-project.iam.gserviceaccount.com`). In Google Drive, right click the folder and click on "Share", then paste the e-mail there and make it so the account has "Editor" privileges, then use the folder ID from above to configure `fbackup`.
* `FBACKUP_EVERY` sets the ticking clock that syncs every n amount of time. By default, it will sync every 30 minutes. Provide a human-readable form of seconds (`20s`), minutes (`10m`), and hours (`7h`).