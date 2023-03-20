# **GoHub Action**


## **Known issues**

It has been observed that a release can break due to rate limits on API calls to Golang packages (action setup).
To "fix" this issue you have to:
1) manually delete the release and the tag in GitHub. 
2) `git pull` down the last changes, remove the release from the changelog
3) Alter `release-please-config.json` by adding `"release-as": "x.y.z"` to you manifest and committing it with a informational `ci: forced x.y.z tag release` commit message
4) remove the `"release-as": "x.y.z"` from the `release-please-config.json` manifest and committing it with a informational `ci: removing forced tag version release` commit message

