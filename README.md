# socialhunter
Crawls the given URL and finds broken social media links that can be hijacked. Broken social links may allow an attacker to conduct phishing attacks. It also can cost a loss of the company's reputation. Broken social media hijack issues are usually accepted on the bug bounty programs.

Currently, it supports Twitter, Facebook, Instagram and Tiktok without any API keys.

[![asciicast](https://asciinema.org/a/wYMVXIHCxxOB3QPWq4Fe8Advn.svg)](https://asciinema.org/a/wYMVXIHCxxOB3QPWq4Fe8Advn)

# Installation

## From Binary

You can download the pre-built binaries from the [releases](https://github.com/utkusen/socialhunter/releases) page and run. For example:

`wget https://github.com/utkusen/socialhunter/releases/download/v0.1.1/socialhunter_0.1.1_Linux_amd64.tar.gz`

`tar xzvf socialhunter_0.1.1_Linux_amd64.tar.gz`

`./socialhunter --help`

## From Source

1. Install Go on your system
2. Run: `go get -u github.com/utkusen/socialhunter`

# Usage

socialhunter requires 2 parameters to run: 

`-f` : Path of the text file that contains URLs line by line. The crawl function is path-aware. For example, if the URL is `https://utkusen.com/blog`, it only crawls the pages under `/blog` path

`-w` : The number of workers to run (e.g `-w 10`). The default value is 5. You can increase or decrease this by testing out the capability of your system.

# Donation

Loved the project? You can buy me a coffee

<a href="https://www.buymeacoffee.com/utkusen" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/default-orange.png" alt="Buy Me A Coffee" height="41" width="174"></a>
