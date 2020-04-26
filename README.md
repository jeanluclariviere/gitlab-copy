# Summary
Gitlab-migrate is a tool to help facilitate the migration of projects from one gitlab environment to another. It decided early on in the project to forgo the use of the `xanzy/go-gitlab` library as their import service is non-functioning.

# Installation
go install github.com/jeanluclariviere/gitlab-migrate

# Usage
**WARNING! Gitlab-migrate leverages gitlab api tokens stored unencrypted in the users `~/.gitlab-migrate/config.json` directory. It is highly recommended that the api tokens created for use with gitlab-migrate be set to expire shortly have use and should not be kept for long term purposes. Use at your own discretion!**

##Setup URIs and tokens: 

Configure credentials

```
$gitlab-migrate setup
Export URI: https://source.gitlab.com
Export Token: **********
Import URI: https://destination.gitlab.com
Import Token: **********

Login to https://source.gitlab.com successful.
Login to https://destination.gitlab.com successful.
```

Migrate
```
$gitlab-migrate 100 group/path
```
