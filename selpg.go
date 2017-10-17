package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

/*===========types==========*/
type selpg_args struct {
	start_page  int
	end_page    int
	in_filename string
	page_len    int  //defalut value "-1"
	page_type   rune //l for lines-delimited,f for form-feed-delimited
	print_dest  string
}

var progname string //program name , for error messages

func main() {
	argc := len(os.Args)
	progname = os.Args[0]
	//for i, a := range os.Args[0:] {
	//	fmt.Printf("Argument %d is %s\n", i+1, a)
	//}
	//fmt.Printf("count: %d", argc)

	var sa selpg_args
	sa.start_page = -1
	sa.end_page = -1
	sa.in_filename = ""
	sa.page_type = 'l'
	sa.print_dest = ""
	sa.page_len = 72

	process_args(argc, os.Args, &sa)
	process_input(sa)
	/*  i, err := strconv.Atoi("42")
	if err != nil {
	    fmt.Printf("couldn't convert number: %v\n", err)
	    return
	}
	fmt.Println("Converted integer:", i)*/
}
func process_args(ac int, args []string, psa *selpg_args) {
	//check the arguments foe validity
	if ac < 3 {
		err := errors.New("no enough arguments")
		panic(err)
		//fmt.Println(progname, ":", err)
		usage()
		os.Exit(1)
	}
	//handel the first argument
	s1 := args[1]
	s1_c := []rune(s1)
	if s1_c[0] != '-' || s1_c[1] != 's' {
		err := errors.New("the first arg should be -sstart_page")
		panic(err) //fmt.Println(progname, ":", err)
		usage()
		os.Exit(2)
	}
	s := s1[2:]
	s_p, err1 := strconv.Atoi(s)
	if err1 != nil {
		panic(err1) //fmt.Printf("couldn't convert number: %v\n", err1)
		return
	}
	if s_p < 1 {
		err1 := errors.New("invalid start page")
		panic(err1) //fmt.Println(progname, ":", err1)
		usage()
		os.Exit(3)
	}
	psa.start_page = s_p

	//handel the second arg
	s2 := args[2]
	s2_c := []rune(s2)
	if s2_c[0] != '-' || s2_c[1] != 'e' {
		err := errors.New("the second arg should be -eend_page")
		panic(err)
		//fmt.Println(progname, ":", err)
		usage()
		os.Exit(4)
	}
	s = s2[2:]
	e_p, err2 := strconv.Atoi(s)
	if err2 != nil {
		panic(err2)
		return
	}
	if e_p < 1 || e_p < s_p {
		err := errors.New("invalid end page")
		panic(err)
		usage()
		os.Exit(5)
	}
	psa.end_page = e_p

	// now handle oprional Args
	argno := 3
	//while there more args and they start with a '-'
	for argno <= ac-1 && args[argno][:1] == "-" {
		s := args[argno]
		s_c := []rune(s)
		switch s_c[1] {
		case 'l':
			l, err := strconv.Atoi(s[2:])
			if err != nil {
				panic(err) //fmt.Printf("couldn't convert number: %v\n", err)
				return
			}
			if l < 1 {
				err = errors.New("invalid page length")
				panic(err) //fmt.Println(progname, ":", err)
				usage()
				os.Exit(6)
			}
			psa.page_len = l
			argno++
		case 'f':
			//check if just "-f" or something more
			if s != "-f" {
				err := errors.New("option should be \"-f\"")
				panic(err) //fmt.Println(progname, ":", err)
				usage()
				os.Exit(7)
			}
			psa.page_type = 'f'
			argno++
		case 'd':
			if s == "-d" {
				err := errors.New("-d option requires a printer destination")
				panic(err) //fmt.Println(progname, ":", err)
				usage()
				os.Exit(8)
			}
			des := s[2:]
			psa.print_dest = des
			argno++
		default:
			fmt.Println("stop")
			err := errors.New("unknown option")
			panic(err) //fmt.Println(progname, ":", err, s)
			usage()
			os.Exit(9)

		}
	}

	//one more arg
	if argno <= (ac - 1) {
		psa.in_filename = args[argno]
		_, err := os.Stat(args[argno])
		if os.IsNotExist(err) {
			err := errors.New("inputfile dose not exist")
			panic(err) //fmt.Println(progname, ":", err, args[argno])
			usage()
			os.Exit(10)
		}
	}
}
func process_input(sa selpg_args) {
	//set the input source
	fin, fout := os.Stdin, os.Stdout
	var page_count, line_count int
	if sa.in_filename != "" {
		f, err := os.Open(sa.in_filename)
		if err != nil {
			panic(err)
			os.Exit(12)
		}
		fin = f
	}
	buff := bufio.NewReader(fin)

	//set the output destination
	if sa.print_dest != "" {
		f, err := os.OpenFile(sa.print_dest, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
			os.Exit(13)
		}
		fout = f
	}

	//begin one of two main loops based on page type
	if sa.page_type == 'l' {
		line_count, page_count = 0, 1
		for true {
			line, err := buff.ReadString('\n')
			if err != nil {
				//panic(err)
				break
			}
			line_count++
			if line_count > sa.page_len {
				page_count++
				line_count = 1
			}
			if page_count >= sa.start_page && page_count <= sa.end_page {
				io.WriteString(fout, line)
			}
		}
	} else {
		page_count = 1
		for true {
			str, err := buff.ReadString('\f')
			if err != nil {
				//panic(err)
				break //EOF
			}
			page_count++
			if page_count >= sa.start_page && page_count <= sa.end_page {
				io.WriteString(fout, str)
			}
		}
	}
	//end main loop
	if page_count < sa.start_page {
		err := errors.New("start_page greater than total pages")
		panic(err) //fmt.Println(progname, ":", err)
	} else if page_count < sa.end_page {
		err := errors.New("end_page greater than total pages")
		panic(err) //fmt.Println(progname, ":", err)
	}
}

func usage() {
	fmt.Printf("\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n", progname)
}
