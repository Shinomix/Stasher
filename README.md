# Stasher
> Add reminders using Slack on top of git stash

## Setup
### Slack side
- Create a new basic app on [Slack admin](https://api.slack.com/apps?new_app=1) and setup the team you be it to integrated into
- In `Features >> OAuth & Permissions`, add new scopes under the `Permissions Scopes` section:
  - chat:write:bot
  - reminders:write
  - users:read
- Once done, `Save Changes` and click on the banner to reinstall the app
- When page has reloaded, copy the OAuth Access Token under the `Token for Your Team` section


### Stasher side
- Clone this repository in your $GOLANG path
- In the `conf/` folder, edit the `conf.json.example` file to add your Slack username and the OAuth Access Token Previously copied and rename it to `conf.json`
- Run `make init; make build ; go install`
- *Optional* For an ease-of-use, run `sudo ln -s $GOPATH/bin/Stasher /usr/bin/Stasher`


## Use
- Simply run `stasher` instead of `git stash` in your git repository
- For detailed use, run `stasher --help`

## To add to the project
- Distribute `Stasher` Slack app on the Slack App Directory
- Auto-install using `go get -u`
