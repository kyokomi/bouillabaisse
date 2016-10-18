# bouillabaisse

bouillabaisse is Firebase Authentication client command-line tool.

![base_12203_50](https://cloud.githubusercontent.com/assets/1456047/17456119/472cc21a-5c07-11e6-8a59-d7977347295b.jpg)

## Install

```sh
go get github.com/kyokomi/bouillabaisse
```

### create authstore.yaml

```sh
$ touch authstore.yaml # config.yaml from `authstorefilename` key
```

## Usage

### dialogue Mode

```sh
$ bouillabaisse
```

### commandLine mode

```sh
$ bouillabaisse --help
```

## Example config.yaml

```yaml

default:
  server:
    listenaddr              : ":8000"
    firebaseapikey          : "<firebase auth api key>"
  local:
    authstoredirpath        : "./"
    authstorefilename       : "authstore.yaml"
  auth:
    authsecretkey           : "firebaseAuth"
    githubclientid          : "<github clientId>"
    githubsecretkey         : "<github secretKey>"
    googleclientid          : "<google clientId>"
    googlesecretkey         : "<google secretKey>"
    facebookclientid        : "<facebook clientId>"
    facebooksecretkey       : "<facebook secretKey>"
    twitterconsumerid       : "<twitter consumerId>"
    twitterconsumersecretkey: "<twitter consumerSecretKey>"
```

## TODO

- [x] Authenticate with Firebase using Password-Based Accounts.
- [x] Authenticate Using Google Sign-In.
- [x] Authenticate Using Facebook Login.
- [x] Authenticate Using Twitter.
- [x] Authenticate Using GitHub.
- [x] Authenticate local save.
- [x] Authenticate local load.
- [x] Show current Authenticate.
- [x] Authenticate with Firebase Anonymously.
- [x] Link Multiple Auth Providers to an Account.
- [x] Manage Users in Firebase.
- [x] Exchange access token and a new refresh token.
- [x] Manage Authenticates at LocalFile.
- [x] Verify Email.
- [x] Password Reset.
- [ ] New Email Accept.
- [x] CommandLine help.
- [x] Remove Local Account.
