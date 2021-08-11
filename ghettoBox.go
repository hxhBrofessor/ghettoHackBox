/*
Title: ghettoBox.go
Author: Bryan Angeles
Notes: Beta
Version: 0.1
Usage: go run ghettoBox.go
--------------------------------------------------------------
Requirements:
#install golang
1- sudo apt install golang -y
2- update .bashrc with the following:
	export GOPATH=$HOME/go
	export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
3- sudo go run ghettoBox.go

Functions:
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
---------------------------------------------------------------
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var colorReset = "\033[0m"
var GREEN = "\033[32m"
var RED = "\033[31m"
var YELLOW = "\033[33m"
var NOCOLOR = "\033[m"

func motd() {
	asciiArt :=
		`
________  ________   ________  ________________________     __  _____   ________ __    ____  ____ _  __
/_  __/ / / / ____/  / ____/ / / / ____/_  __/_  __/ __ \   / / / /   | / ____/ //_/   / __ )/ __ \ |/ /
 / / / /_/ / __/    / / __/ /_/ / __/   / /   / / / / / /  / /_/ / /| |/ /   / ,<     / __  / / / /   / 
/ / / __  / /___   / /_/ / __  / /___  / /   / / / /_/ /  / __  / ___ / /___/ /| |   / /_/ / /_/ /   |  
/_/ /_/ /_/_____/   \____/_/ /_/_____/ /_/   /_/  \____/  /_/ /_/_/  |_\____/_/ |_|  /_____/\____/_/|_|  
																										
                                                                                                                                             
`
	fmt.Println(asciiArt)
}

func checkIfRoot() string {
	//Running exec.command to verify if the user is root
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}
	/*
	   output has trailing \n
	   need to remove the \n
	   otherwise it will cause error for strconv.Atoi
	   log.Println(output[:len(output)-1])
	   0 = root, 501 = non-root user
	*/
	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		log.Fatal(err)
	}
	if i == 0 {
		log.Println((GREEN), "Awesome! You are now running this program with root permissions!", (colorReset))
	} else {
		log.Fatal(string(RED), "[!] This program must be run as root! (sudo)", (colorReset))
	}
	return cmd.String()
}

func checkForInternet() {
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println((GREEN), "[+] Internet connection looks good!", (colorReset))

	} else {
		log.Fatal(string(RED), "[-] Internet connection looks down. You will need internet for this to run (most likely). Fix and try again.", (colorReset))
	}

}

func installKaliRolling() {
	//Adding Kali Rolling into repo
	cmd := exec.Command("sh", "-c", "echo 'deb https://http.kali.org/kali kali-rolling main non-free contrib' > /etc/apt/sources.list.d/kali.list")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf((YELLOW), "Waiting for command to finish adding Kali Linux repositories...", (colorReset))
	err = cmd.Wait()
	log.Printf((GREEN), "Kali Linux repositories added with error: %v", (colorReset), err)
	cmd.ProcessState.ExitCode()

	// Download Kali public key for distro use
	cmd2 := exec.Command("wget", "https://archive.kali.org/archive-key.asc", "-P", "/tmp/")
	err2 := cmd2.Start()
	if err2 != nil {
		log.Fatal(err2)
	}
	log.Printf((YELLOW), "Waiting for Public Key to be installed...", (colorReset))
	err2 = cmd2.Wait()
	log.Printf((GREEN), "Command finished with error: %v", (colorReset), err2)
	cmd2.ProcessState.ExitCode()

	// Installing key
	cmd3 := exec.Command("apt-key", "add", "/tmp/archive-key.asc")
	err3 := cmd3.Start()
	if err3 != nil {
		log.Fatal(err3)
	}
	err3 = cmd3.Wait()
	log.Printf((GREEN), "Adding Key failed: %v", (colorReset), err3)
	cmd3.ProcessState.ExitCode()

	//Setting Prefernce Lower to not break current repo
	cmd4 := exec.Command("sh", "-c", "echo 'Package: *\nPin: Release a=kali-rolling\nPin-Priority: 50' >>/etc/apt/preferences.d/kali.pref")
	err4 := cmd4.Start()
	if err4 != nil {
		log.Fatal(err4)
	}
	log.Printf((YELLOW), "Waiting for kali.list to be added...", (colorReset))
	err4 = cmd4.Wait()
	log.Printf((GREEN), "Command finished with error: %v", (colorReset), err4)
	cmd3.ProcessState.ExitCode()
}

func updateOS() {
	//Updating OS
	log.Printf((YELLOW), "Running Updates")
	cmd := exec.Command("apt", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	cmd2 := exec.Command("apt", "upgrade", "-y")
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	cmd2.Start()

	cmd3 := exec.Command("apt", "dist-upgrade", "-y")
	cmd3.Stdout = os.Stdout
	cmd3.Stderr = os.Stderr
	cmd3.Start()

	log.Printf((GREEN), "[+] Updates Complete", colorReset)
}

func installStarterPackages() {
	starterPack := []string{
		"install",
		"-y",
		"net-tools",
		"curl",
		"git",
	}

	slice := (starterPack[2:])
	for _, str := range slice {
		log.Printf("[*] Attempting installation of the following APT package: %s\n", str)
	}
	cmd := exec.Command("apt", starterPack...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	err := cmd.Wait()
	if err != nil {
		log.Fatal(cmd.Stderr, (RED), "[-] APT Updating failed. Fix and try again. Error:", err, (colorReset))
		os.Exit(1)

	}
	log.Println((GREEN), "[+] Starter Packages installed.", (colorReset))
}
func installAptPackages() {
	aptPackages := []string{
		"install",
		"-y",
		"terminator",
		"flameshot",
		"tmux",
		"torbrowser-launcher",
		"nmap",
		"smbclient",
		"locate",
		"radare2-cutter",
		//"snort", //snort seems freeze the install completely, install must be run on it's own for now.
		"dirb",
		"gobuster",
		"medusa",
		"masscan",
		"whois",
		"autopsy",
		"hashcat",
		"kismet",
		"kismet-plugins",
		"airgraph-ng",
		"wifite",
		"dnsenum",
		"dnsmap",
		"ettercap-common",
		"ettercap-graphical",
		"netdiscover",
		"chromium-browser",
		"python3-pandas",
	}
	slice := (aptPackages[2:])
	for _, str := range slice {
		log.Printf("[*] Attempting installation of the following APT package: %s\n", str)
	}
	cmd := exec.Command("apt", aptPackages...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	err := cmd.Wait()
	if err != nil {
		log.Fatal(cmd.Stderr, (RED), "APT Packages installation failed:", err, colorReset)
		os.Exit(1)

	}
	log.Println((GREEN), "[+] APT Packages installed..", (colorReset))
}

func installMSF() {
	log.Print("[+] Installing Metasploit Framework.")
	cmd := exec.Command("apt", "install", "-t", "kali-rolling", "-y", "metasploit-framework")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	err := cmd.Wait()
	if err != nil {
		log.Fatal(cmd.Stderr, (RED), "-] MSF installation failed. Error:", err, colorReset)
		os.Exit(1)

	}
	log.Println((GREEN), "[+] MSF Installed Successfully.", (colorReset))
}

func installExploitDb() {
	log.Print("[+] Installing ExploitDb.")
	cmd := exec.Command("apt", "install", "-t", "kali-rolling", "-y", "exploitdb")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	err := cmd.Wait()
	if err != nil {
		log.Fatal(cmd.Stderr, (RED), "[-] ExploitDb installation failed. Error:", err, colorReset)
		os.Exit(1)

	}
	log.Println((GREEN), "[+] ExploitDb Installed Successfully.", (colorReset))
}

func installOtherKaliTools() {

	//Add addiotional tools that you want to use into the list below
	kaliRepoTools := []string{
		"install",
		"-t",
		"kali-rolling",
		"-y",
		"bloodhound",
		"enum4linux",
		"feroxbuster",
		"ffuf",
		"nbtscan",
		"nikto",
		"onesixtyone",
		"oscanner",
		"smbclient",
		"smbmap",
		"smtp-user-enum",
		"snmp",
		"sslscan",
		"whatweb",
		"wkhtmltopdf",
		"wpscan",
		"webshells",
		"python3-impacket",
		"kali-tools-windows-resources",
		"john",
		"sqlmap",
		"zaproxy",
		"burpsuite",
		"legion",
		"sparta-scripts",
		"spiderfoot",
		"theharvester",
		"sherlock",
		"maltego",
		"python3-shodan",
		"theharvester",
		"webhttrack",
		"outguess",
		"stegosuite",
		"metagoofil",
		"eyewitness",
		"exifprobe",
		//"ruby-bundler",
		"recon-ng",
		"instaloader",
		"photon",
		"sublist3r",
		"osrframework",
		"drawing",
		"finalrecon",
	}
	slice := (kaliRepoTools[4:])
	for _, str := range slice {
		log.Printf("[*] Attempting installation of the following package from the Kali-Rolling Repo: %s\n", str)
	}
	cmd := exec.Command("apt", kaliRepoTools...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	err := cmd.Wait()
	if err != nil {
		log.Fatal(cmd.Stderr, (RED), "[-] Installation failed. Error:", err, colorReset)
		os.Exit(1)
	}
	log.Println((GREEN), "[+] Kali repo tools installed Successfully.", (colorReset))
}

func installWordlists() {
	log.Print("[+] Installing Wordlists")
	cmd := exec.Command("apt", "install", "-t", "kali-rolling", "-y", "wordlists", "seclists")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	err := cmd.Wait()
	if err != nil {
		log.Fatal(cmd.Stderr, (RED), "[-] Installation failed. Error:", err, colorReset)
		os.Exit(1)

	}
	log.Println((GREEN), "[+] Wordlists installed Successfully.", (colorReset))
}

func main() {
	motd()
	checkIfRoot()
	checkForInternet()
	installKaliRolling()
	updateOS()
	installStarterPackages()
	installAptPackages()
	installMSF()
	installExploitDb()
	installOtherKaliTools()
	installWordlists()
	println("All done!")
}
