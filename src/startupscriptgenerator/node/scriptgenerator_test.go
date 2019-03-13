// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
// --------------------------------------------------------------------------------------------

package main

import (
	"fmt"
)

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_simpleNodeCommand() {
	gen := &NodeStartupScriptGenerator{}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node a/b/c.js
}

func ExampleNodeStartupScriptGenerator_GenerateEntrypointScript_UserStartupCommandIsUsed() {
	gen := &NodeStartupScriptGenerator{
		SourcePath:         "output",
		UserStartupCommand: "abc",
		BindPort:           "8080",
	}
	command := gen.GenerateEntrypointScript()
	fmt.Println(command)
	// Output:
	// #!/bin/sh
	//
	// # Enter the source directory to make sure the script runs where the user expects
	// cd output
	//
	// export PORT=8080
	//
	// if [ -f node_modules.zip ]; then
	//     echo "Found 'node_modules.zip', will extract its contents as node modules."
	//     echo "Removing existing modules directory..."
	//     rm -fr /node_modules
	//     mkdir -p /node_modules
	//     echo "Extracting modules..."
	//     tar -xzf node_modules.zip -C /
	//     echo "Done."
	// fi
	//
	// abc
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_customServerPassedIn() {
	gen := &NodeStartupScriptGenerator{
		CustomStartCommand: "pm2 start --no-daemon",
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// pm2 start --no-daemon a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingFlagShouldBeIncluded() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:                 true,
		RemoteDebuggingBreakBeforeStart: false,
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --inspect a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingbrkFlagShouldBeIncluded() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:                 true,
		RemoteDebuggingBreakBeforeStart: true,
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --inspect-brk a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingWithHostButNoPort() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:   true,
		RemoteDebuggingIp: "0.1.2.3",
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --inspect=0.1.2.3 a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingWithHostAndPort() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:     true,
		RemoteDebuggingIp:   "0.1.2.3",
		RemoteDebuggingPort: "1234",
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --inspect=0.1.2.3:1234 a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingbrkWithHostButNoPort() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:                 true,
		RemoteDebuggingBreakBeforeStart: true,
		RemoteDebuggingIp:               "0.1.2.3",
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --inspect-brk=0.1.2.3 a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingbrkWithHostAndPort() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:                 true,
		RemoteDebuggingBreakBeforeStart: true,
		RemoteDebuggingIp:               "0.1.2.3",
		RemoteDebuggingPort:             "1234",
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --inspect-brk=0.1.2.3:1234 a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingLegacyNodeVersion() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:                 true,
		RemoteDebuggingBreakBeforeStart: false,
		UseLegacyDebugger:               true,
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --debug a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingbrkLegacyNodeVersion() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:                 true,
		RemoteDebuggingBreakBeforeStart: true,
		UseLegacyDebugger:               true,
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --debug-brk a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingbrkWithHostLegacyNodeVersion() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:                 true,
		RemoteDebuggingBreakBeforeStart: true,
		RemoteDebuggingIp:               "0.1.2.3",
		UseLegacyDebugger:               true,
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --debug-brk=0.1.2.3 a/b/c.js
}

func ExampleNodeStartupScriptGenerator_getStartupCommandFromJsFile_debuggingbrkWithHostAndPortLegacyNodeVersion() {
	gen := &NodeStartupScriptGenerator{
		RemoteDebugging:                 true,
		RemoteDebuggingBreakBeforeStart: true,
		RemoteDebuggingIp:               "0.1.2.3",
		RemoteDebuggingPort:             "1234",
		UseLegacyDebugger:               true,
	}
	command := gen.getStartupCommandFromJsFile("a/b/c.js")
	fmt.Println(command)
	// Output:
	// node --debug-brk=0.1.2.3:1234 a/b/c.js
}
