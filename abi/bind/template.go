// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bind

import "github.com/FISCO-BCOS/go-sdk/abi"

// tmplData is the data structure required to fill the binding template.
type tmplData struct {
	Package   string                   // Name of the package to place the generated file in
	Contracts map[string]*tmplContract // List of contracts to generate into this file
	Libraries map[string]string        // Map the bytecode's link pattern to the library name
	Structs   map[string]*tmplStruct   // Contract struct type definitions
}

// tmplContract contains the data needed to generate an individual contract binding.
type tmplContract struct {
	Type        string                 // Type name of the main contract binding
	InputABI    string                 // JSON ABI used as the input to generate the binding from
	InputBin    string                 // Optional EVM bytecode used to denetare deploy code from
	FuncSigs    map[string]string      // Optional map: string signature -> 4-byte signature
	Constructor abi.Method             // Contract constructor for deploy parametrization
	Calls       map[string]*tmplMethod // Contract calls that only read state data
	Transacts   map[string]*tmplMethod // Contract calls that write state data
	Events      map[string]*tmplEvent  // Contract events accessors
	Libraries   map[string]string      // Same as tmplData, but filtered to only keep what the contract needs
	Library     bool                   // Indicator whether the contract is a library
}

// tmplMethod is a wrapper around an abi.Method that contains a few preprocessed
// and cached data fields.
type tmplMethod struct {
	Original   abi.Method // Original method as parsed by the abi package
	Normalized abi.Method // Normalized version of the parsed method (capitalized names, non-anonymous args/returns)
	Structured bool       // Whether the returns should be accumulated into a struct
}

// tmplEvent is a wrapper around an a
type tmplEvent struct {
	Original   abi.Event // Original event as parsed by the abi package
	Normalized abi.Event // Normalized version of the parsed fields
}

// tmplField is a wrapper around a struct field with binding language
// struct type definition and relative filed name.
type tmplField struct {
	Type    string   // Field type representation depends on target binding language
	Name    string   // Field name converted from the raw user-defined field name
	SolKind abi.Type // Raw abi type information
}

// tmplStruct is a wrapper around an abi.tuple contains a auto-generated
// struct name.
type tmplStruct struct {
	Name   string       // Auto-generated struct name(before solidity v0.5.11) or raw name.
	Fields []*tmplField // Struct fields definition depends on the binding language.
}

// tmplSource is language to template mapping containing all the supported
// programming languages the package can generate to.
var tmplSource = map[Lang]string{
	LangGo:         tmplSourceGo,
	LangJava:       tmplSourceJava,
	LangObjC:       tmplSourceObjc,
	LangObjCHeader: tmplSourceObjcHeader,
}

// tmplSourceGo is the Go source template use to generate the contract binding
// based on.
const tmplSourceGo = `
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"math/big"
	"strings"

	"github.com/FISCO-BCOS/go-sdk/abi"
	"github.com/FISCO-BCOS/go-sdk/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/core/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
)

{{$structs := .Structs}}
{{range $structs}}
	// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
	type {{.Name}} struct {
	{{range $field := .Fields}}
	{{$field.Name}} {{$field.Type}}{{end}}
	}
{{end}}

{{range $contract := .Contracts}}
	// {{.Type}}ABI is the input ABI used to generate the binding from.
	const {{.Type}}ABI = "{{.InputABI}}"

	{{if $contract.FuncSigs}}
		// {{.Type}}FuncSigs maps the 4-byte function signature to its string representation.
		var {{.Type}}FuncSigs = map[string]string{
			{{range $strsig, $binsig := .FuncSigs}}"{{$binsig}}": "{{$strsig}}",
			{{end}}
		}
	{{end}}

	{{if .InputBin}}
		// {{.Type}}Bin is the compiled bytecode used for deploying new contracts.
		var {{.Type}}Bin = "0x{{.InputBin}}"

		// Deploy{{.Type}} deploys a new contract, binding an instance of {{.Type}} to it.
		func Deploy{{.Type}}(auth *bind.TransactOpts, backend bind.ContractBackend {{range .Constructor.Inputs}}, {{.Name}} {{bindtype .Type $structs}}{{end}}) (common.Address,  *types.Receipt, *{{.Type}}, error) {
		  parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
		  if err != nil {
		    return common.Address{}, nil, nil, err
		  }
		  {{range $pattern, $name := .Libraries}}
			{{decapitalise $name}}Addr, _, _, _ := Deploy{{capitalise $name}}(auth, backend)
			{{$contract.Type}}Bin = strings.Replace({{$contract.Type}}Bin, "__${{$pattern}}$__", {{decapitalise $name}}Addr.String()[2:], -1)
		  {{end}}
		  address, receipt, contract, err := bind.DeployContract(auth, parsed, common.FromHex({{.Type}}Bin), backend {{range .Constructor.Inputs}}, {{.Name}}{{end}})
		  if err != nil {
		    return common.Address{}, nil, nil, err
		  }
		  return address, receipt, &{{.Type}}{ {{.Type}}Caller: {{.Type}}Caller{contract: contract}, {{.Type}}Transactor: {{.Type}}Transactor{contract: contract}, {{.Type}}Filterer: {{.Type}}Filterer{contract: contract} }, nil
		}

		func AsyncDeploy{{.Type}}(auth *bind.TransactOpts, handler func(*types.Receipt, error), backend bind.ContractBackend {{range .Constructor.Inputs}}, {{.Name}} {{bindtype .Type $structs}}{{end}}) (*types.Transaction, error) {
		  parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
		  if err != nil {
		    return nil, err
		  }
		  {{range $pattern, $name := .Libraries}}
			{{decapitalise $name}}Addr, _, _, _ := AsyncDeploy{{capitalise $name}}(auth, backend)
			{{$contract.Type}}Bin = strings.Replace({{$contract.Type}}Bin, "__${{$pattern}}$__", {{decapitalise $name}}Addr.String()[2:], -1)
		  {{end}}
		  tx, err := bind.AsyncDeployContract(auth, handler, parsed, common.FromHex({{.Type}}Bin), backend {{range .Constructor.Inputs}}, {{.Name}}{{end}})
		  if err != nil {
		    return nil, err
		  }
		  return tx, nil
		}
	{{end}}

	// {{.Type}} is an auto generated Go binding around a Solidity contract.
	type {{.Type}} struct {
	  {{.Type}}Caller     // Read-only binding to the contract
	  {{.Type}}Transactor // Write-only binding to the contract
	  {{.Type}}Filterer   // Log filterer for contract events
	}

	// {{.Type}}Caller is an auto generated read-only Go binding around a Solidity contract.
	type {{.Type}}Caller struct {
	  contract *bind.BoundContract // Generic contract wrapper for the low level calls
	}

	// {{.Type}}Transactor is an auto generated write-only Go binding around a Solidity contract.
	type {{.Type}}Transactor struct {
	  contract *bind.BoundContract // Generic contract wrapper for the low level calls
	}

	// {{.Type}}Filterer is an auto generated log filtering Go binding around a Solidity contract events.
	type {{.Type}}Filterer struct {
	  contract *bind.BoundContract // Generic contract wrapper for the low level calls
	}

	// {{.Type}}Session is an auto generated Go binding around a Solidity contract,
	// with pre-set call and transact options.
	type {{.Type}}Session struct {
	  Contract     *{{.Type}}        // Generic contract binding to set the session for
	  CallOpts     bind.CallOpts     // Call options to use throughout this session
	  TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
	}

	// {{.Type}}CallerSession is an auto generated read-only Go binding around a Solidity contract,
	// with pre-set call options.
	type {{.Type}}CallerSession struct {
	  Contract *{{.Type}}Caller // Generic contract caller binding to set the session for
	  CallOpts bind.CallOpts    // Call options to use throughout this session
	}

	// {{.Type}}TransactorSession is an auto generated write-only Go binding around a Solidity contract,
	// with pre-set transact options.
	type {{.Type}}TransactorSession struct {
	  Contract     *{{.Type}}Transactor // Generic contract transactor binding to set the session for
	  TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
	}

	// {{.Type}}Raw is an auto generated low-level Go binding around a Solidity contract.
	type {{.Type}}Raw struct {
	  Contract *{{.Type}} // Generic contract binding to access the raw methods on
	}

	// {{.Type}}CallerRaw is an auto generated low-level read-only Go binding around a Solidity contract.
	type {{.Type}}CallerRaw struct {
		Contract *{{.Type}}Caller // Generic read-only contract binding to access the raw methods on
	}

	// {{.Type}}TransactorRaw is an auto generated low-level write-only Go binding around a Solidity contract.
	type {{.Type}}TransactorRaw struct {
		Contract *{{.Type}}Transactor // Generic write-only contract binding to access the raw methods on
	}

	// New{{.Type}} creates a new instance of {{.Type}}, bound to a specific deployed contract.
	func New{{.Type}}(address common.Address, backend bind.ContractBackend) (*{{.Type}}, error) {
	  contract, err := bind{{.Type}}(address, backend, backend, backend)
	  if err != nil {
	    return nil, err
	  }
	  return &{{.Type}}{ {{.Type}}Caller: {{.Type}}Caller{contract: contract}, {{.Type}}Transactor: {{.Type}}Transactor{contract: contract}, {{.Type}}Filterer: {{.Type}}Filterer{contract: contract} }, nil
	}

	// New{{.Type}}Caller creates a new read-only instance of {{.Type}}, bound to a specific deployed contract.
	func New{{.Type}}Caller(address common.Address, caller bind.ContractCaller) (*{{.Type}}Caller, error) {
	  contract, err := bind{{.Type}}(address, caller, nil, nil)
	  if err != nil {
	    return nil, err
	  }
	  return &{{.Type}}Caller{contract: contract}, nil
	}

	// New{{.Type}}Transactor creates a new write-only instance of {{.Type}}, bound to a specific deployed contract.
	func New{{.Type}}Transactor(address common.Address, transactor bind.ContractTransactor) (*{{.Type}}Transactor, error) {
	  contract, err := bind{{.Type}}(address, nil, transactor, nil)
	  if err != nil {
	    return nil, err
	  }
	  return &{{.Type}}Transactor{contract: contract}, nil
	}

	// New{{.Type}}Filterer creates a new log filterer instance of {{.Type}}, bound to a specific deployed contract.
 	func New{{.Type}}Filterer(address common.Address, filterer bind.ContractFilterer) (*{{.Type}}Filterer, error) {
 	  contract, err := bind{{.Type}}(address, nil, nil, filterer)
 	  if err != nil {
 	    return nil, err
 	  }
 	  return &{{.Type}}Filterer{contract: contract}, nil
 	}

	// bind{{.Type}} binds a generic wrapper to an already deployed contract.
	func bind{{.Type}}(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	  parsed, err := abi.JSON(strings.NewReader({{.Type}}ABI))
	  if err != nil {
	    return nil, err
	  }
	  return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
	}

	// Call invokes the (constant) contract method with params as input values and
	// sets the output to result. The result type might be a single field for simple
	// returns, a slice of interfaces for anonymous returns and a struct for named
	// returns.
	func (_{{$contract.Type}} *{{$contract.Type}}Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
		return _{{$contract.Type}}.Contract.{{$contract.Type}}Caller.contract.Call(opts, result, method, params...)
	}

	// Transfer initiates a plain transaction to move funds to the contract, calling
	// its default method if one is available.
	func (_{{$contract.Type}} *{{$contract.Type}}Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
		return _{{$contract.Type}}.Contract.{{$contract.Type}}Transactor.contract.Transfer(opts)
	}

	// Transact invokes the (paid) contract method with params as input values.
	func (_{{$contract.Type}} *{{$contract.Type}}Raw) TransactWithResult(opts *bind.TransactOpts, result interface{}, method string, params ...interface{}) (*types.Transaction, *types.Receipt, error) {
		return _{{$contract.Type}}.Contract.{{$contract.Type}}Transactor.contract.TransactWithResult(opts, result, method, params...)
	}

	// Call invokes the (constant) contract method with params as input values and
	// sets the output to result. The result type might be a single field for simple
	// returns, a slice of interfaces for anonymous returns and a struct for named
	// returns.
	func (_{{$contract.Type}} *{{$contract.Type}}CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
		return _{{$contract.Type}}.Contract.contract.Call(opts, result, method, params...)
	}

	// Transfer initiates a plain transaction to move funds to the contract, calling
	// its default method if one is available.
	func (_{{$contract.Type}} *{{$contract.Type}}TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, *types.Receipt, error) {
		return _{{$contract.Type}}.Contract.contract.Transfer(opts)
	}

	// Transact invokes the (paid) contract method with params as input values.
	func (_{{$contract.Type}} *{{$contract.Type}}TransactorRaw) TransactWithResult(opts *bind.TransactOpts, result interface{}, method string, params ...interface{}) (*types.Transaction, *types.Receipt, error) {
		return _{{$contract.Type}}.Contract.contract.TransactWithResult(opts, result, method, params...)
	}

	{{range .Calls}}
		// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatmethod .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}Caller) {{.Normalized.Name}}(opts *bind.CallOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} },{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}}{{end}} error) {
			{{if .Structured}}ret := new(struct{
				{{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}}
				{{end}}
			}){{else}}var (
				{{range $i, $_ := .Normalized.Outputs}}ret{{$i}} = new({{bindtype .Type $structs}})
				{{end}}
			){{end}}
			out := {{if .Structured}}ret{{else}}{{if eq (len .Normalized.Outputs) 1}}ret0{{else}}&[]interface{}{
				{{range $i, $_ := .Normalized.Outputs}}ret{{$i}},
				{{end}}
			}{{end}}{{end}}
			err := _{{$contract.Type}}.contract.Call(opts, out, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
			return {{if .Structured}}*ret,{{else}}{{range $i, $_ := .Normalized.Outputs}}*ret{{$i}},{{end}}{{end}} err
		}

		// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatmethod .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}Session) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if ne $i 0}},{{end}} {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} }, {{else}} {{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}} {{end}} error) {
		  return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.CallOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
		}

		// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatmethod .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}CallerSession) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if ne $i 0}},{{end}} {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} }, {{else}} {{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}} {{end}} error) {
		  return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.CallOpts {{range .Normalized.Inputs}}, {{.Name}}{{end}})
		}
	{{end}}

	{{range .Transacts}}
		// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatmethod .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}Transactor) {{.Normalized.Name}}(opts *bind.TransactOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} },{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}}{{end}} *types.Transaction, *types.Receipt, error) {
			{{if .Structured}}ret := new(struct{
				{{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}}
				{{end}}
			}){{else}}var (
				{{range $i, $_ := .Normalized.Outputs}}ret{{$i}} = new({{bindtype .Type $structs}})
				{{end}}
			){{end}}
			out := {{if .Structured}}ret{{else}}{{if eq (len .Normalized.Outputs) 1}}ret0{{else}}&[]interface{}{
				{{range $i, $_ := .Normalized.Outputs}}ret{{$i}},
				{{end}}
			}{{end}}{{end}}
			transaction, receipt, err := _{{$contract.Type}}.contract.TransactWithResult(opts, out, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
			return {{if .Structured}}*ret,{{else}}{{range $i, $_ := .Normalized.Outputs}}*ret{{$i}},{{end}}{{end}} transaction, receipt, err
		}

		func (_{{$contract.Type}} *{{$contract.Type}}Transactor) Async{{.Normalized.Name}}(handler func(*types.Receipt, error), opts *bind.TransactOpts {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) (*types.Transaction, error) {
					return _{{$contract.Type}}.contract.AsyncTransact(opts, handler, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
		}

		// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatmethod .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}Session) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if ne $i 0}},{{end}} {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} },{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}}{{end}} *types.Transaction, *types.Receipt, error) {
		  return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.TransactOpts {{range $i, $_ := .Normalized.Inputs}}, {{.Name}}{{end}})
		}

		func (_{{$contract.Type}} *{{$contract.Type}}Session) Async{{.Normalized.Name}}(handler func(*types.Receipt, error), {{range $i, $_ := .Normalized.Inputs}}{{if ne $i 0}},{{end}} {{.Name}} {{bindtype .Type $structs}} {{end}}) (*types.Transaction, error) {
		  return _{{$contract.Type}}.Contract.Async{{.Normalized.Name}}(handler, &_{{$contract.Type}}.TransactOpts {{range $i, $_ := .Normalized.Inputs}}, {{.Name}}{{end}})
		}

		// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatmethod .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}TransactorSession) {{.Normalized.Name}}({{range $i, $_ := .Normalized.Inputs}}{{if ne $i 0}},{{end}} {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} },{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}}{{end}} *types.Transaction, *types.Receipt, error) {
		  return _{{$contract.Type}}.Contract.{{.Normalized.Name}}(&_{{$contract.Type}}.TransactOpts {{range $i, $_ := .Normalized.Inputs}}, {{.Name}}{{end}})
		}

		func (_{{$contract.Type}} *{{$contract.Type}}TransactorSession) Async{{.Normalized.Name}}(handler func(*types.Receipt, error), {{range $i, $_ := .Normalized.Inputs}}{{if ne $i 0}},{{end}} {{.Name}} {{bindtype .Type $structs}} {{end}}) (*types.Transaction, error) {
		  return _{{$contract.Type}}.Contract.Async{{.Normalized.Name}}(handler, &_{{$contract.Type}}.TransactOpts {{range $i, $_ := .Normalized.Inputs}}, {{.Name}}{{end}})
		}
	{{end}}

	{{range .Events}}
		// {{$contract.Type}}{{.Normalized.Name}} represents a {{.Normalized.Name}} event raised by the {{$contract.Type}} contract.
		type {{$contract.Type}}{{.Normalized.Name}} struct { {{range .Normalized.Inputs}}
			{{capitalise .Name}} {{if .Indexed}}{{bindtopictype .Type $structs}}{{else}}{{bindtype .Type $structs}}{{end}}; {{end}}
			Raw types.Log // Blockchain specific contextual infos
		}

		// Watch{{.Normalized.Name}} is a free log subscription operation binding the contract event 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatevent .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}Filterer) Watch{{.Normalized.Name}}(fromBlock *uint64, handler func(int, []types.Log){{range .Normalized.Inputs}}{{if .Indexed}}, {{.Name}} {{bindtype .Type $structs}}{{end}}{{end}}) (string, error) {
			return _{{$contract.Type}}.contract.WatchLogs(fromBlock, handler, "{{.Original.Name}}"{{range .Normalized.Inputs}}{{if .Indexed}}, {{.Name}}{{end}}{{end}})
		}

		func (_{{$contract.Type}} *{{$contract.Type}}Filterer) WatchAll{{.Normalized.Name}}(fromBlock *uint64, handler func(int, []types.Log)) (string, error) {
			return _{{$contract.Type}}.contract.WatchLogs(fromBlock, handler, "{{.Original.Name}}")
		}

		// Parse{{.Normalized.Name}} is a log parse operation binding the contract event 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatevent .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}Filterer) Parse{{.Normalized.Name}}(log types.Log) (*{{$contract.Type}}{{.Normalized.Name}}, error) {
			event := new({{$contract.Type}}{{.Normalized.Name}})
			if err := _{{$contract.Type}}.contract.UnpackLog(event, "{{.Original.Name}}", log); err != nil {
				return nil, err
			}
			return event, nil
		}

		// Watch{{.Normalized.Name}} is a free log subscription operation binding the contract event 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatevent .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}Session) Watch{{.Normalized.Name}}(fromBlock *uint64, handler func(int, []types.Log){{range .Normalized.Inputs}}{{if .Indexed}}, {{.Name}} {{bindtype .Type $structs}}{{end}}{{end}}) (string, error)  {
			return _{{$contract.Type}}.Contract.Watch{{.Normalized.Name}}(fromBlock, handler {{range .Normalized.Inputs}}{{if .Indexed}}, {{.Name}}{{end}}{{end}})
		}

		func (_{{$contract.Type}} *{{$contract.Type}}Session) WatchAll{{.Normalized.Name}}(fromBlock *uint64, handler func(int, []types.Log)) (string, error) {
			return _{{$contract.Type}}.Contract.WatchAll{{.Normalized.Name}}(fromBlock, handler)
		}

		// Parse{{.Normalized.Name}} is a log parse operation binding the contract event 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{formatevent .Original $structs}}
		func (_{{$contract.Type}} *{{$contract.Type}}Session) Parse{{.Normalized.Name}}(log types.Log) (*{{$contract.Type}}{{.Normalized.Name}}, error) {
			return _{{$contract.Type}}.Contract.Parse{{.Normalized.Name}}(log)
		}

 	{{end}}
{{end}}
`

// tmplSourceJava is the Java source template use to generate the contract binding
// based on.
const tmplSourceJava = `
// This file is an automatically generated Java binding. Do not modify as any
// change will likely be lost upon the next re-generation!

package {{.Package}};

import org.ethereum.geth.*;
import java.util.*;

{{$structs := .Structs}}
{{range $contract := .Contracts}}
{{if not .Library}}public {{end}}class {{.Type}} {
	// ABI is the input ABI used to generate the binding from.
	public final static String ABI = "{{.InputABI}}";
	{{if $contract.FuncSigs}}
		// {{.Type}}FuncSigs maps the 4-byte function signature to its string representation.
		public final static Map<String, String> {{.Type}}FuncSigs;
		static {
			Hashtable<String, String> temp = new Hashtable<String, String>();
			{{range $strsig, $binsig := .FuncSigs}}temp.put("{{$binsig}}", "{{$strsig}}");
			{{end}}
			{{.Type}}FuncSigs = Collections.unmodifiableMap(temp);
		}
	{{end}}
	{{if .InputBin}}
	// BYTECODE is the compiled bytecode used for deploying new contracts.
	public final static String BYTECODE = "0x{{.InputBin}}";

	// deploy deploys a new contract, binding an instance of {{.Type}} to it.
	public static {{.Type}} deploy(TransactOpts auth, EthereumClient client{{range .Constructor.Inputs}}, {{bindtype .Type $structs}} {{.Name}}{{end}}) throws Exception {
		Interfaces args = Geth.newInterfaces({{(len .Constructor.Inputs)}});
		String bytecode = BYTECODE;
		{{if .Libraries}}

		// "link" contract to dependent libraries by deploying them first.
		{{range $pattern, $name := .Libraries}}
		{{capitalise $name}} {{decapitalise $name}}Inst = {{capitalise $name}}.deploy(auth, client);
		bytecode = bytecode.replace("__${{$pattern}}$__", {{decapitalise $name}}Inst.Address.getHex().substring(2));
		{{end}}
		{{end}}
		{{range $index, $element := .Constructor.Inputs}}Interface arg{{$index}} = Geth.newInterface();arg{{$index}}.set{{namedtype (bindtype .Type $structs) .Type}}({{.Name}});args.set({{$index}},arg{{$index}});
		{{end}}
		return new {{.Type}}(Geth.deployContract(auth, ABI, Geth.decodeFromHex(bytecode), client, args));
	}

	// Internal constructor used by contract deployment.
	private {{.Type}}(BoundContract deployment) {
		this.Address  = deployment.getAddress();
		this.Deployer = deployment.getDeployer();
		this.Contract = deployment;
	}
	{{end}}

	// Ethereum address where this contract is located at.
	public final Address Address;

	// Ethereum transaction in which this contract was deployed (if known!).
	public final Transaction Deployer;

	// Contract instance bound to a blockchain address.
	private final BoundContract Contract;

	// Creates a new instance of {{.Type}}, bound to a specific deployed contract.
	public {{.Type}}(Address address, EthereumClient client) throws Exception {
		this(Geth.bindContract(address, ABI, client));
	}

	{{range .Calls}}
	{{if gt (len .Normalized.Outputs) 1}}
	// {{capitalise .Normalized.Name}}Results is the output of a call to {{.Normalized.Name}}.
	public class {{capitalise .Normalized.Name}}Results {
		{{range $index, $item := .Normalized.Outputs}}public {{bindtype .Type $structs}} {{if ne .Name ""}}{{.Name}}{{else}}Return{{$index}}{{end}};
		{{end}}
	}
	{{end}}

	// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.ID}}.
	//
	// Solidity: {{.Original.String}}
	public {{if gt (len .Normalized.Outputs) 1}}{{capitalise .Normalized.Name}}Results{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}}{{end}}{{end}} {{.Normalized.Name}}(CallOpts opts{{range .Normalized.Inputs}}, {{bindtype .Type $structs}} {{.Name}}{{end}}) throws Exception {
		Interfaces args = Geth.newInterfaces({{(len .Normalized.Inputs)}});
		{{range $index, $item := .Normalized.Inputs}}Interface arg{{$index}} = Geth.newInterface();arg{{$index}}.set{{namedtype (bindtype .Type $structs) .Type}}({{.Name}});args.set({{$index}},arg{{$index}});
		{{end}}

		Interfaces results = Geth.newInterfaces({{(len .Normalized.Outputs)}});
		{{range $index, $item := .Normalized.Outputs}}Interface result{{$index}} = Geth.newInterface(); result{{$index}}.setDefault{{namedtype (bindtype .Type $structs) .Type}}(); results.set({{$index}}, result{{$index}});
		{{end}}

		if (opts == null) {
			opts = Geth.newCallOpts();
		}
		this.Contract.call(opts, results, "{{.Original.Name}}", args);
		{{if gt (len .Normalized.Outputs) 1}}
			{{capitalise .Normalized.Name}}Results result = new {{capitalise .Normalized.Name}}Results();
			{{range $index, $item := .Normalized.Outputs}}result.{{if ne .Name ""}}{{.Name}}{{else}}Return{{$index}}{{end}} = results.get({{$index}}).get{{namedtype (bindtype .Type $structs) .Type}}();
			{{end}}
			return result;
		{{else}}{{range .Normalized.Outputs}}return results.get(0).get{{namedtype (bindtype .Type $structs) .Type}}();{{end}}
		{{end}}
	}
	{{end}}

	{{range .Transacts}}
	// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.ID}}.
	//
	// Solidity: {{.Original.String}}
	public Transaction {{.Normalized.Name}}(TransactOpts opts{{range .Normalized.Inputs}}, {{bindtype .Type $structs}} {{.Name}}{{end}}) throws Exception {
		Interfaces args = Geth.newInterfaces({{(len .Normalized.Inputs)}});
		{{range $index, $item := .Normalized.Inputs}}Interface arg{{$index}} = Geth.newInterface();arg{{$index}}.set{{namedtype (bindtype .Type $structs) .Type}}({{.Name}});args.set({{$index}},arg{{$index}});
		{{end}}
		return this.Contract.transact(opts, "{{.Original.Name}}"	, args);
	}
	{{end}}
}
{{end}}
`

const tmplSourceObjc = `// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.
{{$structs := .Structs}}{{range $contract := .Contracts}}
#import "{{.Type}}.h"
#import <FiscoBcosIosSdk/FiscoBcosIosSdk.h>

{{$contractName := .Type}}
{{range $structs}}
// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
@implementation {{$contractName}}_{{.Name}}
@end
{{end}}

@implementation {{.Type}}

/// init
- (instancetype) init:(MobileBcosSDK *)sdk{
    if (self = [super init]){
		self = [self init];
        self.sdk = sdk;
    }
    return self;
}
/// initWithAddress
- (instancetype) initWithAddress:(NSString *)addr
                             sdk:(MobileBcosSDK *)sdk{
	if (self = [super init]){
        self = [self init:sdk];
		self.address = addr;
	}
	return self;
}

{{if .InputBin}}
/// deploy {{range .Constructor.Inputs}}
/// @param {{.Name}} {{.Type}} type argument{{objcPrintArgComment .Type}}{{end}}{{range .Constructor.Outputs}}
/// @return {{.Name}} {{.Type}} type argument{{end}}
- (MobileReceiptResult*) deploy {{range .Constructor.Inputs}}:({{objcFormatStructType .Type $structs $contractName}}) {{.Name}}{{end}}{
	{{if .Constructor.Inputs}}NSArray * __resArr = @[
        {{range $i, $_ :=.Constructor.Inputs}}{{if ne $i 0}},
		{{else}}{{end}}@{
            @"type":@"{{.Type}}",
            @"value":{{objcFormattedValue .Type .Name $structs}}
        }{{end}}
    ];{{else}}{{end}}
	{{if .Constructor.Inputs}}NSString *__params = [self __stringFromArr:__resArr];{{else}}NSString *__params = @"[]";{{end}}
	return [self.sdk deployContract:self.abi contractBin:_bin params:__params];
}
{{end}}

{{range .Calls}}
/// {{.Normalized.Name}}{{range .Normalized.Inputs}}
/// @param {{.Name}} {{.Type}} type argument{{objcPrintArgComment .Type}}{{end}}{{range .Normalized.Outputs}}
/// @return {{.Name}} {{.Type}} type argument{{end}}
- (MobileCallResult *) {{.Normalized.Name}} {{range .Normalized.Inputs}}:({{objcFormatStructType .Type $structs $contractName}}) {{.Name}}{{end}}{
	{{if .Normalized.Inputs}}NSArray * __resArr = @[
        {{range $i, $_ :=.Normalized.Inputs}}{{if ne $i 0}},
		{{else}}{{end}}@{
            @"type":@"{{.Type}}",
            @"value":{{objcFormattedValue .Type .Name $structs}}
        }{{end}}
    ];{{else}}{{end}}
	{{if .Normalized.Inputs}}NSString * __params = [self __stringFromArr:__resArr];{{else}}NSString *__params = @"[]";{{end}}
	return [self.sdk call:self.abi address:self.address method:@"{{.Original.Name}}" params:__params outputNum:{{objcGetOutputNum .Normalized.Outputs}}];
}
{{end}}

{{range .Transacts}}
/// {{.Normalized.Name}}{{range .Normalized.Inputs}}
/// @param {{.Name}} {{.Type}} type argument{{objcPrintArgComment .Type}}{{end}}{{range .Normalized.Outputs}}
/// @return {{.Name}} {{.Type}} type argument{{end}}
- (MobileReceiptResult *) {{.Normalized.Name}} {{range $i, $_ := .Normalized.Inputs}}{{if ne $i 0}}
	{{.Name}}:{{else}} :{{end}}({{objcFormatStructType .Type $structs $contractName}}) {{.Name}}{{end}}{
	{{if .Normalized.Inputs}}NSArray * __resArr = @[
        {{range $i, $_ :=.Normalized.Inputs}}{{if ne $i 0}},
		{{else}}{{end}}@{
            @"type":@"{{.Type}}",
            @"value":{{objcFormattedValue .Type .Name $structs}}
        }{{end}}
    ];{{else}}{{end}}
	{{if .Normalized.Inputs}}NSString *__params = [self __stringFromArr:__resArr];{{else}}NSString *__params = @"[]";{{end}}
	return [self.sdk sendTransaction:self.abi address:self.address method:@"{{.Original.Name}}" params:__params];
}
{{end}}

- (instancetype)init{
	self.abi = @"{{.InputABI}}";
	self.bin = @"0x{{.InputBin}}";
	return self;
}

- (NSString *)__stringFromArr:(NSArray *)arr{
	NSString *jsonString = nil;
    NSError *error;
    NSData *jsonData = [NSJSONSerialization dataWithJSONObject:arr
                                                       options:NSJSONWritingPrettyPrinted
                                                         error:&error];
    if (! jsonData) {
        NSLog(@"Got an error: %@", error);
        return @"[???]";
    } else {
        jsonString = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
        return jsonString;
    }
}


@end
{{end}}
`

const tmplSourceObjcHeader = `#import <Foundation/Foundation.h>
#import <FiscoBcosIosSdk/FiscoBcosIosSdk.h>

NS_ASSUME_NONNULL_BEGIN
{{$structs := .Structs}}{{range $contract := .Contracts}}{{$contractName := .Type}}
{{range $structs}}// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
@interface {{$contractName}}_{{.Name}} : NSObject
{{range $field := .Fields}}@property (nonatomic, assign) {{$field.Type}} {{$field.Name}};
{{end}}
@end
{{end}}
@interface {{.Type}} : NSObject

@property (nonatomic, copy) NSString *address;
@property (nonatomic, copy) NSString *abi;
@property (nonatomic, copy) NSString *bin;
@property (nonatomic, strong) MobileBcosSDK *sdk;

/// init
- (instancetype)init:(MobileBcosSDK *)sdk;
/// initWithAddress
- (instancetype)initWithAddress:(NSString *)addr
	sdk:(MobileBcosSDK *) sdk;
{{if .InputBin}}
/// deploy {{range .Constructor.Inputs}}
/// @param {{.Name}} {{.Type}} type argument{{objcPrintArgComment .Type}}{{end}}{{range .Constructor.Outputs}}
/// @return {{.Name}} {{.Type}} type argument{{end}}
- (MobileReceiptResult*)  deploy {{range $i, $_ := .Constructor.Inputs}}{{if ne $i 0}}
	{{.Name}}:{{else}} :{{end}}({{objcFormatStructType .Type $structs $contractName}}) {{.Name}}{{end}}; {{end}}

{{range .Calls}}
/// {{.Normalized.Name}}{{range .Normalized.Inputs}}
/// @param {{.Name}} {{.Type}} type argument{{objcPrintArgComment .Type}}{{end}} {{range .Normalized.Outputs}}
/// @return {{.Name}} {{.Type}} type argument{{end}}
- (MobileCallResult *) {{.Normalized.Name}} {{range $i, $_ := .Normalized.Inputs}}{{if ne $i 0}}
	{{.Name}}:{{else}} :{{end}}({{objcFormatStructType .Type $structs $contractName}}) {{.Name}}{{end}};
{{end}}

{{range .Transacts}}
/// {{.Normalized.Name}}{{range .Normalized.Inputs}}
/// @param {{.Name}} {{.Type}} type argument{{objcPrintArgComment .Type}}{{end}}{{range .Normalized.Outputs}}
/// @return {{.Name}} {{.Type}} type argument{{end}}
- (MobileReceiptResult *) {{.Normalized.Name}} {{range $i, $_ := .Normalized.Inputs}}{{if ne $i 0}}
	{{.Name}}:{{else}} :{{end}}({{objcFormatStructType .Type $structs $contractName}}) {{.Name}}{{end}};
{{end}}

@end
{{end}}
NS_ASSUME_NONNULL_END
`
