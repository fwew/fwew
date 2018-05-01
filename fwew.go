//	Fwew is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU General Public License as published by
// 	the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
//	Fwew is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU General Public License for more details.
//
//	You should have received a copy of the GNU General Public License
//	along with Fwew.  If not, see http://gnu.org/licenses/

// Package main obviously contains all the stuff for the main program
package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tirea/fwew/affixes"
	"github.com/tirea/fwew/config"
	"github.com/tirea/fwew/numbers"
	"github.com/tirea/fwew/util"

	_ "github.com/mattn/go-sqlite3"
)

// Global
const (
	idField  int = 0 // dictionary.tsv line Field 0 is Database ID
	lcField  int = 1 // dictionary.tsv line field 1 is Language Code
	navField int = 2 // dictionary.tsv line field 2 is Na'vi word
	ipaField int = 3 // dictionary.tsv line field 3 is IPA data
	infField int = 4 // dictionary.tsv line field 4 is Infix location data
	posField int = 5 // dictionary.tsv line field 5 is Part of Speech data
	defField int = 6 // dictionary.tsv line field 6 is Local definition
)

var debug *bool

func fwewExperimental(word, lc, posFilter string, reverse, useAffixes bool) []affixes.Word {
	badChars := strings.Split("` ~ @ # $ % ^ & * ( ) [ ] { } < > _ / . , ; : ! ? | + \\", " ")
	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		word = strings.Replace(word, c, "", -1)
	}

	var results []affixes.Word
	var navi, ipa, infixes, pos, def string
	var err error

	database, _ := sql.Open("sqlite3", util.Text("database"))
	query := "SELECT navi, ipa, infixes, fancyPartOfSpeech, definition " +
		"FROM entries " +
		"WHERE navi = \"" + word + "\" " +
		"OR definition LIKE \"" + word + "\""
	rows, _ := database.Query(query)

	for rows.Next() {
		var result affixes.Word
		err = rows.Scan(&navi, &ipa, &infixes, &pos, &def)
		if err != nil {
			break
		}
		result.Navi = navi
		result.IPA = ipa
		result.PartOfSpeech = pos
		result.Definition = def
		results = append(results, result)
	}

	return results
}

func fwew(word, lc, posFilter string, reverse, useAffixes bool) []affixes.Word {
	badChars := strings.Split("` ~ @ # $ % ^ & * ( ) [ ] { } < > _ / . , ; : ! ? | + \\", " ")
	word = strings.ToLower(word)
	// remove all the sketchy chars from arguments
	for _, c := range badChars {
		word = strings.Replace(word, c, "", -1)
	}

	var result affixes.Word
	var results []affixes.Word
	var fields []string
	var defString string
	var added bool

	// Prepare file for searching
	dictData, err := os.Open(util.Text("dictionary"))
	defer dictData.Close()
	if err != nil {
		fmt.Println(errors.New(util.Text("noDataError")))
		os.Exit(1)
	}
	scanner := bufio.NewScanner(dictData)

	// Go through each line and see what we can find
	for scanner.Scan() {
		line := scanner.Text()
		// Store the fields of the line into fields array in lowercase
		fields = strings.Split(line, "\t")

		// Looking for Local word in Definition field
		if reverse {
			// Must hard code "all" here now that config voids guarantee of default filter "all"
			if posFilter == "all" {
				if fields[lcField] == lc {
					// whole-word matching
					defString = util.StripChars(fields[defField], ",;")
					for _, w := range strings.Split(defString, " ") {
						if strings.ToLower(w) == strings.ToLower(word) && !added {
							// Put the stuff from fields into the Word struct
							result = affixes.InitWordStruct(result, fields)
							results = append(results, result)
							added = true
						}
					}
				}
				// filter part of speech
			} else {
				if fields[lcField] == lc && fields[posField] == posFilter {
					// whole-word matching
					defString = util.StripChars(fields[defField], ",;")
					for _, w := range strings.Split(defString, " ") {
						if strings.ToLower(w) == strings.ToLower(word) && !added {
							// Put the stuff from fields into the Word struct
							result = affixes.InitWordStruct(result, fields)
							results = append(results, result)
							added = true
						}
					}
				}
			}
			added = false
			// Looking for Na'vi word in Na'vi field
		} else {
			if fields[lcField] == lc {
				if strings.ToLower(fields[navField]) == strings.ToLower(word) {
					// Put the stuff from fields into the Word struct
					result = affixes.InitWordStruct(result, fields)
					results = append(results, result)
					break
				} else {
					result = affixes.InitWordStruct(result, fields)
					result.Target = word
					result = affixes.Reconstruct(result)

					if result.ID != "-1" {
						results = append(results, result)
					}
					// reset these to not catch the next word
					result.Target = ""
					result.Attempt = ""

				}
			}
		}
	}

	return results
}

func printResults(results []affixes.Word, reverse, showInfixes, showIPA, useAffixes, markdown bool) {
	if len(results) != 0 {
		var out string

		for i, w := range results {
			num := fmt.Sprintf("[%d] ", i+1)
			nav := fmt.Sprintf("%s", w.Navi)
			ipa := fmt.Sprintf("[%s] ", w.IPA)
			pos := fmt.Sprintf("%s", w.PartOfSpeech)
			inf := fmt.Sprintf("%s ", w.InfixLocations)
			def := fmt.Sprintf("%s\n", w.Definition)

			if markdown {
				nav = "**" + nav + "** "
				pos = "*" + pos + "* "
			} else {
				nav += " "
				pos += " "
			}

			out += num
			out += nav
			if showIPA {
				out += ipa
			}
			if showInfixes && w.InfixLocations != "\\n" {
				out += inf
			}
			out += pos
			out += def
			if useAffixes && len(w.Affixes) > 0 {
				for key, value := range w.Affixes {
					out += fmt.Sprintf("    %s: %s\n", key, value)
				}
			}
		}
		out += fmt.Sprintf("\n")

		fmt.Printf(out)

	} else {
		fmt.Println(util.Text("none"))
	}
}

func setFlags(arg string, debug, r, i, ipa, a, n *bool, l, p *string) {
	const start int = 4 // s,e,t,[ = 0,1,2,3
	flagList := strings.Split(arg[start:len(arg)-1], ",")
	setList := []string{}

	for _, f := range flagList {
		switch {
		case f == "":
			fmt.Printf("<! %s: debug=%t, r=%t, i=%t, ipa=%t, a=%t, l=%s, p=%s >\n\n", util.Text("cset"), *debug, *r, *i, *ipa, *a, *l, *p)
		case f == "debug":
			*debug = true
			setList = append(setList, f)
		case f == "r":
			*r = true
			setList = append(setList, f)
		case f == "i":
			*i = true
			setList = append(setList, f)
		case f == "ipa":
			*ipa = true
			setList = append(setList, f)
		case f == "a":
			*a = true
			setList = append(setList, f)
		case f == "n":
			*n = true
			setList = append(setList, f)
		case strings.HasPrefix(f, "l="):
			*l = f[2:]
			setList = append(setList, f)
		case strings.HasPrefix(f, "p="):
			*p = f[2:]
			setList = append(setList, f)
		default:
			fmt.Printf("<! %s: %s >\n", util.Text("noOptionError"), f)
		}
	}

	if len(setList) != 0 {
		fmt.Printf("<! %v %s >\n\n", setList, util.Text("set"))
	}
}

func unsetFlags(arg string, debug, r, i, ipa, a, n *bool) {
	const start int = 6 // u,n,s,e,t,[ = 0,1,2,3,4,5
	flagList := strings.Split(arg[6:len(arg)-1], ",")
	unsetList := []string{}
	for _, f := range flagList {
		switch f {
		case "":
			fmt.Println()
		case "debug":
			*debug = false
			unsetList = append(unsetList, f)
		case "r":
			*r = false
			unsetList = append(unsetList, f)
		case "i":
			*i = false
			unsetList = append(unsetList, f)
		case "ipa":
			*ipa = false
			unsetList = append(unsetList, f)
		case "a":
			*a = false
			unsetList = append(unsetList, f)
		case "n":
			*n = false
			unsetList = append(unsetList, f)
		default:
			fmt.Printf("<! %s: %s >\n", util.Text("noOptionError"), f)
		}
	}
	if len(unsetList) != 0 {
		fmt.Printf("<! %v %s >\n\n", unsetList, util.Text("unset"))
	}
}

func main() {
	var configuration = config.ReadConfig()
	var results []affixes.Word
	var language, posFilter *string
	var showVersion, showInfixes, showIPA, reverse, useAffixes, numConvert, markdown, experimental *bool

	// Debug flag, for verbose probing output
	debug = flag.Bool("debug", false, util.Text("usageDebug"))
	// Version flag, for displaying version data
	showVersion = flag.Bool("v", false, util.Text("usageV"))
	// Reverse direction flag, for local_lang -> Na'vi lookups
	reverse = flag.Bool("r", false, util.Text("usageR"))
	// Language specifier flag
	language = flag.String("l", configuration.Language, util.Text("usageL"))
	// Infixes flag, opt to show infix location data
	showInfixes = flag.Bool("i", false, util.Text("usageI"))
	// IPA flag, opt to show IPA data
	showIPA = flag.Bool("ipa", false, util.Text("usageIPA"))
	// Show part of speech flag
	posFilter = flag.String("p", configuration.PosFilter, util.Text("usageP"))
	// Attempt to find all matches using affixes
	useAffixes = flag.Bool("a", configuration.UseAffixes, util.Text("usageA"))
	// Convert numbers
	numConvert = flag.Bool("n", false, util.Text("usageN"))
	// Markdown formatting
	markdown = flag.Bool("m", false, util.Text("usageM"))
	// Experimental Database algo switch
	experimental = flag.Bool("e", false, util.Text("usageE"))
	flag.Parse()

	if *showVersion {
		fmt.Println(util.Version)
		if flag.NArg() == 0 {
			os.Exit(0)
		}
	}

	// ARGS MODE
	if flag.NArg() > 0 {
		for _, arg := range flag.Args() {
			arg = strings.Replace(arg, "’", "'", -1)
			if strings.HasPrefix(arg, "set[") && strings.HasSuffix(arg, "]") {
				setFlags(arg, debug, reverse, showInfixes, showIPA, useAffixes, numConvert, language, posFilter)
			} else if strings.HasPrefix(arg, "unset[") && strings.HasSuffix(arg, "]") {
				unsetFlags(arg, debug, reverse, showInfixes, showIPA, useAffixes, numConvert)
			} else {
				if *numConvert {
					fmt.Println(numbers.Convert(arg, *reverse))
				} else {
					if *experimental {
						results = fwewExperimental(arg, *language, *posFilter, *reverse, *useAffixes)
					} else {
						results = fwew(arg, *language, *posFilter, *reverse, *useAffixes)
					}
					printResults(results, *reverse, *showInfixes, *showIPA, *useAffixes, *markdown)
				}
			}
		}

		// INTERACTIVE MODE
	} else {
		fmt.Println(util.Text("header"))

		for {
			fmt.Print(util.Text("prompt"))

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			input = strings.Trim(input, "\n")
			input = strings.Replace(input, "’", "'", -1)

			// catch EOF error
			if err != nil {
				fmt.Println()
				os.Exit(0)
			}

			if input != "" {
				if strings.HasPrefix(input, "set[") && strings.HasSuffix(input, "]") {
					setFlags(input, debug, reverse, showInfixes, showIPA, useAffixes, numConvert, language, posFilter)
				} else if strings.HasPrefix(input, "unset[") && strings.HasSuffix(input, "]") {
					unsetFlags(input, debug, reverse, showInfixes, showIPA, useAffixes, numConvert)
				} else {
					if *numConvert {
						fmt.Println(numbers.Convert(input, *reverse))
					} else {
						if *experimental {
							results = fwewExperimental(input, *language, *posFilter, *reverse, *useAffixes)
						} else {
							results = fwew(input, *language, *posFilter, *reverse, *useAffixes)
						}
						printResults(results, *reverse, *showInfixes, *showIPA, *useAffixes, *markdown)
					}
				}
			} else {
				fmt.Println()
			}
		}
	}
}
