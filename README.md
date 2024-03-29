# utils
Random useful stuff for various projects

# getjwt

- Backend dev utility for generating a valid firebase JWT for testing protected APIs.
- Requires an active firebase project, IAM API enables and GCP service account with sufficient permission

Usage:

For windows:

- Download the executable in `/bin`

For Linux/Mac

- [Download Go](https://go.dev/doc/install)
- Clone repo
- `cd /getjwt`
- `go build -o getjwt.exe` to build on your architecture
* I will add compiled binaries for more architectures later

To run
- Option 1: `getjwt -f <path/to.config.yml`
- Option 2: `getjwt <FIREBASE_WEB_API_KEY> <FIREBASE_USER_ID> <PATH_TO_GOOGLE_JSON>`

For option 1, a `yaml` file containing the following items is required:
```yml
apiKey: AI...
userId: ...
googleCreds: path/to/google.json
```
