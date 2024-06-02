# Golang Ebiten project with local server

# how to run
## option 1 - run main file
 `go run ./ebiten/main.go`
 
 this will simply run the main file
## option 2 - build and run windows exe
 `cd ebiten`
 
 `.\updateBuilds.bat`

 now you have a windows exe file in ebiten/builds/windows-amd64.exe

you are ready to run the file, however note your system antivirus will need to scan your exe file

## opiton 3 - run in browser
for this, you need to have run the updatebuilds.bat file from step 2

assuming you have done step 2, and are now on the root again

`cd server`

`.\buildForAllPlatforms.bat`

`cd builds`

you are ready to run the file, however note your system antivirus will need to scan your exe file

`.\windows-amd64.exe`

open the browser and got to [localhost:900](http://localhost:9090/static/)

# controls
keys that are listened too are

spacebar - to select a piece

the four standard arrow keys - to move around the board

esc - to get up the menue

enter - to end turn (will automaticaly end turn when you have 0 actions remaining)