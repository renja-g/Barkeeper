# Barkeeper
> Barkeeper is a Discord bot that help to organize custom matches, by generating fair teams and also provides statistical data of past matches.

## Commands
- `/rate @user 9`
  > Rates a user.

- `/info @user`
  > Shows the information of a user.

- `/list filter=Online, Offline`
  > Shows a list of all users and there stats. If a filter is set the list for example only contains online uses. Giving a quick overview if enough would be there for a custom match.

- `/teams`
  > generates two teams with the people in the voice channel the command user is in. The teams will be balanced using the rating. The message will have two buttons:
  > 
  > [Start Match]
  > [Reshuffle Teams]
  > 
  > If "Start Match" is pressed, the user will be automatically moved into their teams voice channels and the buttons change to:
  > [Team 1 wins]
  > [Team 2 wins]
  > [Cancel match]
  >
  > When a match ends, the stats of all participants will be updated.
  <img width="318" alt="image" src="https://github.com/renja-g/Barkeeper/assets/76645494/b191ba3a-ba07-4897-9562-430c844e64db">
  <img width="427" alt="image" src="https://github.com/renja-g/Barkeeper/assets/76645494/7dd78292-0aea-4713-8d87-edc14a268e9a">
  <img width="249" alt="image" src="https://github.com/renja-g/Barkeeper/assets/76645494/43fe4833-11e1-48e1-8943-5bd0d6c149a3">




- `/history user=@user`
  > Shows all past matches. Or only the matches of a certain user.

- `/help`
  > Shows the help message.

## Setup
0. Make sure that Docker is installed on your system.
1. Clone the repository `git clone http://github.com/renja-g/Barkeeper`
2. Rename the `example.config.json` file to `config.json` and fill in the required information.
3. Run `docker-compose up -d` in the root directory of the project.
4. The bot should now be running and you can invite it to your server.
