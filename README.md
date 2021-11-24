# jira-releaser
Parse a Github release and update all related Jira tickets.

## Why?

Often, a Github release will be created which contains links to tickets in Jira.

This project will find all those links, and set the fixVersion in Jira to the name of the released tag in Github.

## Installation

This project is distributed as a single binary. [Download the latest release](https://github.com/kieranajp/jira-releaser/releases/latest) for your OS and architecture, un-tar it, and add it to your path - making sure it's executable.

In unix-speak, that's `curl | tar xf | chmod +x | mv`.

## Usage

You'll need to provide credentials for both Jira and Github for this to work properly. These can be provided either as command line flags, or environment variables:

```
   --github-user, -g value          Github username [$GITHUB_USER]
   --github-token value, -t value   Github token [$GITHUB_TOKEN]
   --jira-url value, -j value       Jira URL (default: "https://jira.example.com") [$JIRA_URL]
   --jira-user value, -u value      Jira Username [$JIRA_USER]
   --jira-password value, -p value  Jira Password [$JIRA_PASS]
```

A Jira personal access token can also be used instead of a password. 

You will also need to provide the URL to the appropriate Github repository, and the release name (as `--release`). Assuming the above are all set as environment variables, using the tool therefore looks like:

```
jira-releaser -r v1.2.3 https://github.com/octocat/hello-world
```

The body of release v1.2.3 of that repository will be parsed, and `hello-world v1.2.3` will be set as the fix version of any tickets mentioned.

