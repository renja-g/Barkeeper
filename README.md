# Barkeeper
> Barkeeper is a Discord bot that help to organize custom matches, by generating fair teams and hand also provides statistical data of past matches.

## Commands
- `/rate @user 9`
  > Rates a user.

- `/info @user`
  > Shows the information of the user.

- `/list`
  > Shows a list of all users and there stats, sorted by a modified Wilson score algorithm.

- `/teams`
  > generates two teams with the people in the voice channel with the command user. The teams will be balanced using the rating. The message will have two buttons:
  > 
  > [Start Match]
  > [Reshuffle Teams]
  > 
  > If "Start Match" is pressed, the user will be automatically moved into their teams voice channels and the buttons change to:
  > [Team 1 Wins]
  > [Team 2 Wins]
  > [Cancel Match]
  >
  > When a match ends, the stats of all participants will be updated.

- `/history`
  > Shows all past matches.

- `/help`
  > Shows the help message.
