package main

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"os"
	"plugin"
	"testing"
)

type Man interface {
	SaveLive() error
}

func TestPlugin(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// load module 插件您也可以使用go http.Request从远程下载到本地,在加载做到动态的执行不同的功能
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(dir + "/plugin_M_m.so")
	if err != nil {
		panic(err)
	}
	println("plugin opened")

	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Greeter
	m, err := plug.Lookup("M") //大写
	if err != nil {
		panic(err)
	}

	// 3. Assert that loaded symbol is of a desired type
	man, ok := m.(Man)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	// 4. use the module
	if err := man.SaveLive(); err != nil {
		println("use plugin man failed, ", err)
	}

}

func Test_plugin_with_gengine(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dc := context.NewDataContext()
	//3.load plugin into apiName, exportApi
	_, _, e := dc.PluginLoader(dir + "/plugin_M_m.so")
	if e != nil {
		panic(e)
	}

	dc.Add("println", fmt.Println)
	ruleBuilder := builder.NewRuleBuilder(dc)
	err = ruleBuilder.BuildRuleFromString(`
	rule "1"
	begin
	 
	//this method is defined in plugin
	err = m.SaveLive()
	if isNil(err) {
	   println("err is nil")
	}
	end
	`)

	if err != nil {
		panic(err)
	}
	gengine := engine.NewGengine()
	err = gengine.Execute(ruleBuilder, false)

	if err != nil {
		panic(err)
	}
}
