# **gohub-action**

A demo pipeline of how to build, version and consume github actions written in Golang!

## **General idea**

The purpose of this projects is to build a Golang binary, embed it to a tag/release as an asset and then finally invocate the composite action in order to run that binary from a linux runtime.\
The invocation step consists of downloading the Golang asset, decompress & chmodding it before running it by passing parameters as environment variables to propagate inputs.

## **Release**

The [release action](./.github/workflows/release.yaml) depends on [release-please-action](https://github.com/google-github-actions/release-please-action) in lockstep with [go-release-action](https://github.com/wangyoucao577/go-release-action).\
The former is responsible for creating a release whilst the latter builds and appends the compressed (tar.gz) binary to the released tag.

Whilst the `release-please-action` is configured through the usage of a [manifest](https://github.com/googleapis/release-please/blob/main/docs/manifest-releaser.md), the `go-release-action` requires some tinkering inside the release pipeline when/if porting it to another go-binary reusable actions workflow. In particular the following variables will have downstream effects for the usability of the reusable composite action:
- binary_name: **xxx**
- asset_name: **yyy**-linux-amd64

Of the two variables, the former (`binary_name`) will be the name of the binary inside the compressed asset.\
The latter, `asset_name`, will be the name of the compressed artifact that is appended to the released tag.

It's advisable to keep **tag** and **GOOS/GOARCH** outside of the binary name (as this can be gathered from the asset_name and it's omission will greatly improve reusability).\
In terms of reusability it's also advisable to keep the **tag** outside of the asset_name, but include the `GOOS` and `GOARCH` variables (for multi OS support).

## **Reusable Action**

The action itself is configured in the [root level action.yml](action.yml) file.

This composite action relies on the [release-downloader](https://github.com/robinraju/release-downloader) action in order to download the artifact:

```yaml
- uses: robinraju/release-downloader@v1.7
    with: 
      repository: "tsanton/gohub-action"
      tag: ${{ inputs.tag }}
      fileName: "greeter-linux-amd64.tar.gz"
      # Relative path under $GITHUB_WORKSPACE to place the downloaded file(s)
      # It will create the target directory automatically if not present
      # eg: out-file-path: "my-downloads" => It will create directory $GITHUB_WORKSPACE/my-downloads
      out-file-path: "artifact-download"
```
As we can see from above, the `fileName` can be hardcoded if we ommit the tag from that name and simply controll which version of the code we're running targeting a tag as configured by the input variable `tag`.\
This can of course be made more dynamic by adding `go_os` and `go_arch` as input variables if needed.

Lastly there's the complexity added by the `binary_name`:

```yaml
 - name: run
   shell: bash
   run: |
       tar -xf ./artifact-download/greeter-linux-amd64.tar.gz --directory $GITHUB_WORKSPACE/artifact-download
       chmod +x ./artifact-download/greeter
       ./artifact-download/greeter
```

As we can see we first have to untar the downloaded artifact. This is a temporary step as there is a PR on the way for automatically unzipping (see **TODO**).\
We will still have the challenge that we must chmod the binary prior to invocating it. Therefore it's advisable to give it a simple name that will persist across all releases and simply hardcode the same value in both the release pipeline & reusable composite action yaml.


## **Invocation**

To incocate this action, simply run the following workflow from our `.github/workflows` folder:

```yaml
name: test-invocate-reusable-action

on:
  workflow_dispatch:

jobs:
  go-echo:
    name: release run
    runs-on: ubuntu-latest
    steps:
    
    - uses: tsanton/gohub-action@0.1.11
      with:
        tag: "0.1.11"
        message: "Eat pudding and fall asleep"
        name: "You fool!"
```

Which should give you the nice output of: ```Eat pudding and fall asleep, You fool!```

## **Known issues**

It has been observed that a release can break due to rate limits on API calls to Golang packages (action setup).
To "fix" this issue you have to:
1) manually delete the release and the tag in GitHub. 
2) `git pull` down the last changes, remove the release from the changelog
3) Alter `release-please-config.json` by adding `"release-as": "x.y.z"` to you manifest and committing it with a informational `ci: forced x.y.z tag release` commit message
4) remove the `"release-as": "x.y.z"` from the `release-please-config.json` manifest and committing it with a informational `ci: removing forced tag version release` commit message

## **Sauce**

Inspired by [this](https://full-stack.blend.com/how-we-write-github-actions-in-go.html)

## **TODO**

When [this PR](https://github.com/robinraju/release-downloader/pull/613) is merged, refactor the run step in the reusable action to remove the ```tar-xf ./<src> ./<dest>```.