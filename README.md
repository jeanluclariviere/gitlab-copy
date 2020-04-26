# Summary
Gitlab-migrate is a tool to help facilitate the migration of projects from one gitlab environment to another. It decided early on in the project to forgo the use of the `xanzy/go-gitlab` library as their import service is non-functioning.

# Installation
```
go get github.com/jeanluclariviere/gitlab-migrate
cd $GOPATH/src/github.com/jeanluclariviere/gitlab-migrate
go install
```

# Usage
**WARNING! Gitlab-migrate leverages gitlab api tokens stored unencrypted in the users `~/.gitlab-migrate/config.json` directory. It is highly recommended that the api tokens created for use with gitlab-migrate be set to expire shortly have use and should not be kept for long term purposes. Use at your own discretion!**

## Source
The input source should be a valid project ID

## Destination
If the destination is ommited, gitlab-migrate will migrate the project to the token's owner's projects:

```
$ gitlab-migrate 1
```

Imports the project to: `administrator/example`


If the destination is supplied, it will create all parent subgroups as necessary and place the project in the last group:
```
$ gitlab-migrate 1 hello/world
```

Imports the project to: `hello/world/example`

## Setup URIs and tokens: 

### Configure credentials

```
$gitlab-migrate setup
Export URI: https://source.gitlab.com
Export Token: **********
Import URI: https://destination.gitlab.com
Import Token: **********

Login to https://source.gitlab.com successful.
Login to https://destination.gitlab.com successful.
```

### Validate existing credentials (prompt to create if ~/.gitlab-migrate is missing)

```
$gitlab-migrate login
Login to https://source.gitlab.com successful.
Login to https://destination.gitlab.com successful.
```

## Migration

### Migrate to a group
```
$gitlab-migrate 100 group
2020/04/26 14:17:59 Scheduling export...
2020/04/26 14:17:59 Export status: finished
2020/04/26 14:17:59 Downloading ./04-26-2020-sample.tar.gz
2020/04/26 14:17:59 Creating groups
2020/04/26 14:17:59 Importing project
2020/04/26 14:18:00 Import complete
```

### Migrate to a subgroup
```
$gitlab-migrate 100 group/subgroup
2020/04/26 14:17:59 Scheduling export...
2020/04/26 14:17:59 Export status: finished
2020/04/26 14:17:59 Downloading ./04-26-2020-sample.tar.gz
2020/04/26 14:17:59 Creating groups
2020/04/26 14:17:59 Importing project
2020/04/26 14:18:00 Import complete
```

### Migrate to token's owner's projects (ommit the group)
```
$gitlab-migrate 100 
2020/04/26 14:17:59 Scheduling export...
2020/04/26 14:17:59 Export status: finished
2020/04/26 14:17:59 Downloading ./04-26-2020-sample.tar.gz
2020/04/26 14:17:59 Creating groups
2020/04/26 14:17:59 Importing project
2020/04/26 14:18:00 Import complete
```
