# `fbackup`, it's never too late to make backups

`fbackup` is a simple tool that takes a directory in your filesystem, a Google Drive destination, and proceeds to upload every file from the local folder into Google Drive. If the local file's modification date is newer than the one reported by Google Drive, **the file will be overwritten prefering the local file** over the remote one. Extraneous files (files that exist in Google Drive but they don't exist in the local folder) **will be deleted from Google Drive** and the preference is always to persist local files over remote files anytime.

## Usage

The project is meant to be used with Docker in mind, so settings are provided via environment variables:

* `FBACKUP_QUIET` will make the program not output a single thing to `stdout`.
* `FBACKUP_ACCOUNT_FILE` a path to the local `service-account.json` file. You can [create service accounts here](https://console.cloud.google.com/iam-admin/serviceaccounts). Don't specify any permission in particular in the "Grant this service account access to the project" tab from the walkthrough, and simply create it. Continue reading for more details about the account.
* `FBACKUP_FOLDER` the target folder you want to sync to Google Drive (say, `~/Projects`).
* `FBACKUP_DESTINATION_ID` the ID of the Google Drive folder where you want to store the folder you're uploading / synchronizing. To get an ID, simply go to `drive.google.com`, navigate to the folder you want to use by double clicking on it, then in the URL you'll see something like this: `https://drive.google.com/drive/u/1/folders/1z-xSDBGxR6WzktgIca4HdiipQP9sJ9lG`, note, at the end, that there's a random ID, `1z-xSDBGxR6WzktgIca4HdiipQP9sJ9lG`. That's the destination ID you'll be using to synchronize. Keep in mind that your service account **won't see this folder** because it belongs to your personal account, **unless you explicitly share it** with the account. To do this, open the account file from the second bullet point, find the e-mail address associated with it in the `client_email` field from the JSON file (it will look like `my-service-account@my-project.iam.gserviceaccount.com`). In Google Drive, right click the folder and click on "Share", then **paste the e-mail there and make it so the account has "Editor" privileges**, then use the folder ID from above to configure `fbackup`.
* `FBACKUP_EVERY` sets the ticking clock that syncs every n amount of time. By default, it will sync every 30 minutes. Provide a human-readable form of seconds (`20s`), minutes (`10m`), and hours (`7h`).

## Usage with Docker

`fbackup` is automatically published [under `patrickdappollonio/fbackup` in Docker Hub](https://hub.docker.com/r/patrickdappollonio/fbackup). You can use this image to synchronize any folder in your computer using Docker. In this case, the current folder contains the `account.json` file used to log in to Google Drive, and a folder called `my-photos`.

The process to back up this using Docker is:
* We mount the account file inside the container (the location doesn't matter) as read-only and we provide its location to the `FBACKUP_ACCOUNT_FILE`
* We mount the folder we want to synchronize / back up into any location in the container as read-only, and we provide its location to `FBACKUP_FOLDER`
* We provide the folder ID where we want to save the contents of `my-photos` in Google Drive using the ID (see above for how to get the Google Drive ID of a folder)
* We assume the default ticker to synchronize every 30 minutes, although you can provide `FBACKUP_EVERY` to change that

In the end, the result looks more or less like this:

```bash
$ ls
account.json  my-photos

$ docker run -d --restart=unless-stopped \
  -v "$(pwd)/account.json:/account.json:ro" \
  -e FBACKUP_ACCOUNT_FILE="/account.json" \
  -v "$(pwd)/my-photos:/backup:ro" \
  -e FBACKUP_FOLDER="/backup" \
  -e FBACKUP_DESTINATION_ID="1z-xSDBGxR6WzktgIca4HdiipQP9sJ9lG" \
  patrickdappollonio/fbackup
9ec2a6410ead26d8f948168d1a1543c6c2d9f44a26a3908a23df0f18b06a0a67
```

This will mount the contents of the folder inside the container, as well as the account file as read-only (since the application doesn't really need to change anything in it) and then start synchronizing. You can retrieve the logs of the container by looking at the ID provided at the end, in our case, `9ec2a6410ead26d8f948168d1a1543c6c2d9f44a26a3908a23df0f18b06a0a67`. You can use `docker logs -f 9ec2a641` to get the logs (pro tip: I put the first 8 or so characters of the ID, but Docker will know it belogs to that container):

```bash
$ docker logs 9ec2a6
[fbackup] 2020/01/01 00:00:01 Setting account file to: /account.json
[fbackup] 2020/01/01 00:00:01 Setting target folder to: /backup
[fbackup] 2020/01/01 00:00:01 Setting ticker to: 30m0s
[fbackup] 2020/01/01 00:00:01 Configured service account with account file: /account.json
[fbackup] 2020/01/01 00:00:01 Google Drive client configured, retrieving folder information for ID: 1z-xSDBGxR6WzktgIca4HdiipQP9sJ9lG
# remaining log statements here
```

The container will continue to run unless it's stopped (as stated by `--restart=unless-stopped` above) and it will run on detached mode (provided with `-d` above).