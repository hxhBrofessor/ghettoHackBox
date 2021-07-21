# ghettoHackBox

## Requirements:

1- `sudo apt install golang -y`

2- update .bashrc with the following:

	export GOPATH=$HOME/go
	export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
  
3- sudo go run ghettoBox.go
​
### Functions:
	+motd()
	+checkIfRoot()
	+checkForInternet()
	+installKaliRolling()
	+updateOS()
	+installStarterPackages()
	+installAptPackages()
	+installMSF()
	+installExploitDb()
	+installOtherKaliTools() //If theres a tool you'd like to add go to the function and add it to the list before running the program
	+installWordlists()
  
#### Note:
***To install a desktop enviroment you must use the following repo `-t kali-rolling` to successfully install kde for example.***
