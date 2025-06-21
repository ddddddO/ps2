package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ddddddO/ps2"
)

const (
	base   = `ps2 <<< '< Data serialized by PHP serialize function >'`
	sample = `ps2 <<< 'a:9:{s:10:"string_val";s:27:"こんにちは、世界！";s:7:"int_val";i:123;s:9:"bool_true";b:1;s:10:"bool_false";b:0;s:8:"null_val";N;s:9:"float_val";d:3.14159;s:18:"nested_assoc_array";a:3:{s:4:"name";s:12:"Go Developer";s:7:"details";a:2:{s:3:"age";i:30;s:4:"city";s:8:"Kawasaki";}s:7:"hobbies";a:3:{i:0;s:6:"coding";i:1;s:7:"reading";i:2;s:6:"hiking";}}s:13:"indexed_array";a:5:{i:0;s:9:"りんご";i:1;s:9:"バナナ";i:2;s:12:"チェリー";i:3;i:100;i:4;b:1;}s:15:"object_instance";O:8:"MyObject":3:{s:10:"publicProp";s:15:"パブリック";s:16:"*protectedProp";i:456;s:19:"MyObjectprivateProp";a:1:{s:3:"key";s:5:"value";}}}'`
)

func main() {
	var toJSON, toYAML, toTOML bool
	flag.BoolVar(&toJSON, "json", true, "convert to JSON")
	flag.BoolVar(&toYAML, "yaml", false, "convert to YAML")
	flag.BoolVar(&toTOML, "toml", false, "convert to TOML")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s is: %s\n", os.Args[0], base)
		fmt.Fprintf(os.Stderr, "Here's a quick example you can try:\n\n")
		fmt.Fprintf(os.Stderr, "%s\n\n", sample)
		flag.PrintDefaults()
	}
	flag.Parse()

	optionOfOutputType := ps2.WithOutputTypeJSON()
	if toYAML {
		optionOfOutputType = ps2.WithOutputTypeYAML()
	}
	if toTOML {
		optionOfOutputType = ps2.WithOutputTypeTOML()
	}

	output, err := ps2.Run(os.Stdin, optionOfOutputType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(output)
}
