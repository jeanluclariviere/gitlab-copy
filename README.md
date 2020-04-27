# Summary
Gitlab-copy is a tool to help facilitate the copying of projects from one gitlab environment to another. I decided early on to forgo the use of the `xanzy/go-gitlab` library as their import service is non-functioning.

# Installation
```
go get github.com/jeanluclariviere/gitlab-copy
cd $GOPATH/src/github.com/jeanluclariviere/gitlab-copy
go install
```

# Usage
**WARNING! Gitlab-copy leverages gitlab api tokens stored unencrypted in the users `~/.gitlab-copy/config.json` directory. It is highly recommended that the api tokens created for use with gitlab-copy be set to expire shortly after use and that they should not be kept for long term purposes. Use at your own discretion!**

## Source Project
The input should be a valid project ID:
```
$ gitlab-copy 100 example
```

Will copy the project with ID 100

## Destination
If the destination is ommited, gitlab-copy will copy the project to the token's owner's projects:

```
$ gitlab-copy 1
```

Imports the project to: `administrator/example`


If the destination is supplied, it will create all parent subgroups as necessary and place the project in the last group:
```
$ gitlab-copy 1 hello/world
```

Imports the project to: `hello/world/example`

## Setup URIs and tokens: 

### Configure credentials

```
$ gitlab-copy setup
Export URI: https://source.gitlab.com
Export Token: **********
Import URI: https://destination.gitlab.com
Import Token: **********

Login to https://source.gitlab.com successful.
Login to https://destination.gitlab.com successful.
```

### Validate existing credentials (prompt to create if ~/.gitlab-copy is missing)

```
$ gitlab-copy login
Login to https://source.gitlab.com successful.
Login to https://destination.gitlab.com successful.
```

## Migration

### Migrate to a group
```
$ gitlab-copy 100 group
2020/04/26 14:17:59 Scheduling export...
2020/04/26 14:17:59 Export status: finished
2020/04/26 14:17:59 Downloading ./04-26-2020-sample.tar.gz
2020/04/26 14:17:59 Creating groups
2020/04/26 14:17:59 Importing project
2020/04/26 14:18:00 Import complete
```

### Migrate to a subgroup
```
$ gitlab-copy 100 group/subgroup
2020/04/26 14:17:59 Scheduling export...
2020/04/26 14:17:59 Export status: finished
2020/04/26 14:17:59 Downloading ./04-26-2020-sample.tar.gz
2020/04/26 14:17:59 Creating groups
2020/04/26 14:17:59 Importing project
2020/04/26 14:18:00 Import complete
```

### Migrate to token's owner's projects (ommit the group)
```
$ gitlab-copy 100 
2020/04/26 14:17:59 Scheduling export...
2020/04/26 14:17:59 Export status: finished
2020/04/26 14:17:59 Downloading ./04-26-2020-sample.tar.gz
2020/04/26 14:17:59 Creating groups
2020/04/26 14:17:59 Importing project
2020/04/26 14:18:00 Import complete
```
