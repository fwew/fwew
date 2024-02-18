//	Fwew is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU General Public License as published by
// 	the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
//	Fwew is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without gen implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
//	along with Fwew.  If not, see http://gnu.org/licenses/

// Package main contains all the things
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	fwew "github.com/fwew/fwew-lib/v5"
)

// global constants
const (
	space string = " "
)

// global flags & options
var (
	configuration           Config
	configure, filename     *string
	language, posFilter     *string
	showInfixes, showIPA    *bool
	skipFixes, showInfDots  *bool
	showDashed, showVersion *bool
	showSource, useAffixes  *bool
	numConvert, markdown    *bool
	debug, reverse          *bool
)

func setFlags(arg string, argsMode bool) {
	var (
		flagList []string
		err      error
		langs    = strings.Split(Text("languages"), ", ")
	)
	if argsMode {
		start := IndexStr(arg, '[') + 1
		flagList = strings.Split(arg[start:len(arg)-1], ",")
	} else {
		flagList = strings.Split(arg, space)
	}
	for _, f := range flagList {
		switch {
		case f == "":
		case f == "r":
			*reverse = !*reverse
		case f == "i":
			*showInfixes = !*showInfixes
		case f == "ipa":
			*showIPA = !*showIPA
		case f == "id":
			*showInfDots = !*showInfDots
		case f == "s":
			*showDashed = !*showDashed
		case f == "skip":
			*skipFixes = !*skipFixes
		case f == "src":
			*showSource = !*showSource
		case f == "a":
			*useAffixes = !*useAffixes
		case f == "n":
			*numConvert = !*numConvert
		case f == "m":
			*markdown = !*markdown
		case f == "d":
			*debug = !*debug
		case strings.HasPrefix(f, "l="):
			if ContainsStr(langs, f[2:]) {
				*language = f[2:]
			} else {
				err = fmt.Errorf("%s: %s (%s: %s)", Text("invalidLanguageError"), f[2:], Text("options"), Text("languages"))
				fmt.Println(err)
				fmt.Println()
			}
		case strings.HasPrefix(f, "p="):
			*posFilter = f[2:]
		default:
			err = fmt.Errorf("%s (%s)", Text("noOptionError"), f)
			fmt.Println(err)
			fmt.Println()
		}
	}
	if err == nil {
		var out string
		out += Text("set") + space
		out += "[ "
		if *reverse {
			out += "r "
		}
		if *showInfDots {
			out += "id "
		}
		if *showDashed {
			out += "s "
		}
		if *showInfixes {
			out += "i "
		}
		if *showIPA {
			out += "ipa "
		}
		if *showSource {
			out += "src "
		}
		if *useAffixes {
			out += "a "
		}
		if *numConvert {
			out += "n "
		}
		if *markdown {
			out += "m "
		}
		if *debug {
			out += "d "
		}
		out += fmt.Sprintf("l=%s p=%s", *language, *posFilter)
		out += " ]\n"
		if len(*filename) == 0 {
			fmt.Println(out)
		}
	}
}

func printHelp() {
	flag.Usage = func() {
		fmt.Printf("%s: ", Text("usage"))
		fmt.Printf("%s [%s] [%s]\n", Text("bin"), Text("options"), Text("w_words"))
		fmt.Printf("%s:\n", Text("options"))
		flag.PrintDefaults()
	}
	flag.Usage()
}

func output(words [][]fwew.Word) {
	fileMode := len(*filename) > 0
	for _, wordbundle := range words {
		for j, word := range wordbundle {
			if word.ID == "" && !fileMode {
				continue
			}
			if word.ID == "" && fileMode {
				fmt.Printf("cmd %s\n", word.Navi)
			} else {
				entry, err := word.ToOutputLine(fmt.Sprint(j), *markdown, *showIPA, *showInfixes, *showDashed, *showInfDots, *showSource, *language)
				if err != nil {
					panic(err)
				}
				fmt.Println(entry)
			}
		}
	}
	if len(words) == 0 {
		fmt.Println(Text("none"))
	}
}

func slashCommand(s string, argsMode bool) {
	var (
		sc      []string
		command string
		args    []string
		numArgs int
		setArg  string
		confArg string
		k       int
		err     error
		words   []fwew.Word
	)
	sc = strings.Split(s, space)
	sc = DeleteEmpty(sc)
	command = sc[0]
	if len(sc) > 1 {
		args = sc[1:]
		numArgs = len(args)
	}
	switch command {
	case "/help":
		printHelp()
	case "/commands":
		fmt.Println(Text("slashCommandHelp"))
	case "/set", "/unset":
		setArg = strings.Join(args, space)
		setFlags(setArg, argsMode)
	// aliases for /set
	case "/a", "/id", "/s", "/i", "/ipa", "/l", "/n", "/p", "/r", "/src":
		for _, c := range command {
			if c != '/' {
				setArg += string(c)
			}
		}
		if numArgs > 0 {
			setArg += space
		}
		setArg += strings.Join(args, space)
		setFlags(setArg, argsMode)
	case "/list":
		words, err := fwew.List(args, 1)
		if err != nil {
			panic(err)
		}
		output([][]fwew.Word{words})
	case "/random":
		if numArgs > 0 {
			// get number of random words requested
			k, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(Text("invalidNumericError") + "\n")
			} else {
				// get filter arguments
				if numArgs >= 5 && args[1] == "where" {
					args = args[2:]
				} else {
					args = []string{}
				}
				words, err = fwew.Random(k, args, 1)
				if err != nil {
					panic(err)
				}
				output([][]fwew.Word{words})
			}
		} else {
			fmt.Println()
		}
	case "/lenition", "/len":
		fmt.Println(Text("lenTable"))
	case "/update":
		err := fwew.UpdateDict()
		if err != nil {
			panic(err)
		}
		// Version.DictBuild = SHA1Hash(Text("dictionary"))
	case "/version":
		// fmt.Println(Version)
		fmt.Println(fwew.Version.String())
	case "/config":
		confArg = strings.Join(args, space)
		configuration = WriteConfig(confArg)
	case "/quit", "/exit", "/q", "/wc":
		os.Exit(0)
	default:
		fmt.Println()
	}
}

func main() {
	// Assure Dictionary is downloaded
	err := fwew.AssureDict()
	if err != nil {
		panic(err)
	}
	var (
		argsMode bool
		fileMode bool
	)
	fwew.StartEverything()

	configuration = ReadConfig()
	// Version flag, for displaying version data
	showVersion = flag.Bool("v", false, Text("usageV"))
	// Reverse direction flag, for local_lang -> Na'vi lookups
	reverse = flag.Bool("r", configuration.Reverse, Text("usageR"))
	// Language specifier flag
	language = flag.String("l", configuration.Language, Text("usageL"))
	// Infixes flag, opt to show infix location data
	showInfixes = flag.Bool("i", configuration.ShowInfixes, Text("usageI"))
	// Infix locations in dot notation
	showInfDots = flag.Bool("id", configuration.ShowInfDots, Text("usageID"))
	// IPA flag, opt to show IPA data
	showIPA = flag.Bool("ipa", configuration.ShowIPA, Text("usageIPA"))
	// Show syllable breakdown / stress
	skipFixes = flag.Bool("skip", configuration.ShowIPA, Text("usageSkip"))
	// Show syllable breakdown / stress
	showDashed = flag.Bool("s", configuration.ShowDashed, Text("usageS"))
	// Source flag, opt to show source data
	showSource = flag.Bool("src", configuration.ShowSource, Text("usageSrc"))
	// Filter part of speech flag, opt to filter by part of speech
	posFilter = flag.String("p", configuration.PosFilter, Text("usageP"))
	// Attempt to find all matches using affixes
	useAffixes = flag.Bool("a", configuration.UseAffixes, Text("usageA"))
	// Convert numbers
	numConvert = flag.Bool("n", configuration.NumConvert, Text("usageN"))
	// Markdown formatting
	markdown = flag.Bool("m", configuration.Markdown, Text("usageM"))
	// Input file / Fwewscript
	filename = flag.String("f", "", Text("usageF"))
	// Configuration editing
	configure = flag.String("c", "", Text("usageC"))
	// Debug mode
	debug = flag.Bool("d", configuration.DebugMode, Text("usageD"))
	flag.Parse()
	argsMode = flag.NArg() > 0
	fileMode = len(*filename) > 0

	if *showVersion {
		fmt.Println(fwew.Version.String())
		if flag.NArg() == 0 {
			os.Exit(0)
		}
	}

	if *configure != "" {
		configuration = WriteConfig(*configure)
		os.Exit(0)
	}

	if fileMode { // FILE MODE
		if *markdown {
			// restrict Discord users to cwd
			*filename = "./" + *filename
		}
		inFile, err := os.Open(*filename)
		if err != nil {
			fmt.Printf("%s (%s)\n", Text("noFileError"), *filename)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(inFile)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "#") && line != "" {
				executor(line)
			}
		}
		err = inFile.Close()
		if err != nil {
			fmt.Println(Text("fileCloseError"))
			os.Exit(1)
		}
	} else if argsMode { // ARGS MODE
		for _, arg := range flag.Args() {
			arg = strings.Replace(arg, "’", "'", -1)
			if strings.HasPrefix(arg, "set[") && strings.HasSuffix(arg, "]") {
				setFlags(arg, argsMode)
			} else if strings.HasPrefix(arg, "unset[") && strings.HasSuffix(arg, "]") {
				setFlags(arg, argsMode)
			} else {
				executor(arg)
			}
		}
	} else { // INTERACTIVE MODE
		fmt.Println(Text("header"))

		p := prompt.New(executor, completer,
			prompt.OptionTitle(Text("name")),
			prompt.OptionPrefix(Text("prompt")),
			prompt.OptionSelectedDescriptionTextColor(prompt.DarkGray),
		)
		p.Run()
	}
}
