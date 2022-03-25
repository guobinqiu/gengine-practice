package gengine_practice

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"plugin"
	"testing"
)

//定义规则，由运营部门需提供统计指标、决策树
const template = `
rule "rule1" "测试顺序执行模式" salience 1 //格式：rule "规则名" "规则描述" salience 优先级(数字大的先执行)
begin
localvar = 1 //局部变量
println("rule name:", @name, "local var:", localvar, "global var:", globalvar) //输出规则名，局部变量，全局变量
action1() //业务逻辑1，可以是函数，方法或api
action2() //业务逻辑2
action3() //业务逻辑3
return "hello" //规则也可以有返回值
end

rule "rule2" "测试并发执行模式" salience 2
begin
localvar = 3
println("rule name:", @name, "local var:", localvar, "global var:", globalvar)
conc  {//并发执行业务逻辑123，需保证它们线程之间是安全的
action1()
action2()
action3()
}
return "world"
end

//统计学生成绩等级指标
rule "calc-grade-excellent" "" salience 3 //相同优先级之间执行顺序随机，描述可以为空
begin
	if student.Score >= 80 && student.Score <= 100 {
		student.Grade = "你简直是个天才"
	}
end

rule "calc-grade-good" "" salience 3
begin
	if student.Score >=60 && student.Score < 80 {
		student.Grade = "你真是棒棒哒！"
	}
end

rule "calc-grade-potential" "" salience 3
begin
	if student.Score >=0 && student.Score < 60 {
		student.Grade = "你上升空间巨大，前途无量啊！"
	}
end
`

func action1() {
	fmt.Println("action1")
}

func action2() {
	fmt.Println("action2")
}

func action3() {
	fmt.Println("action3")
}

type Student struct {
	Score int
	Grade string
}

func TestSingle(t *testing.T) {
	student := &Student{} //必须是指针类型
	student.Score = 59    //修改分数以得到不同的等级

	dataContext := context.NewDataContext()
	dataContext.Add("action1", action1)     //注入一个自定义函数
	dataContext.Add("action2", action2)     //注入一个自定义函数
	dataContext.Add("action3", action3)     //注入一个自定义函数
	dataContext.Add("println", fmt.Println) //注入一个go的标准函数
	dataContext.Add("globalvar", 2)         //注入一个全局变量，所有规则里都可以拿到它
	dataContext.Add("student", student)     //注入一个结构体变量

	//初始化规则引擎
	eng := engine.NewGengine()

	//构建规则
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	if err := ruleBuilder.BuildRuleFromString(template); err != nil {
		t.Fatalf("build rule err: %v\n", err)
	}

	//执行规则
	if err := eng.Execute(ruleBuilder, true); err != nil {
		t.Fatalf("execute rule error: %v\n", err)
	}

	//输出规则的返回值
	resultMap, _ := eng.GetRulesResultMap()
	t.Log(resultMap)

	//输出学生成绩等级指标
	t.Logf("student.Grade=%s\n", student.Grade)

	//修改参数再次测试
	fmt.Println("--------------------修改参数再次测试")
	student.Score = 100
	dataContext.Add("globalvar", 200) //会覆盖原来的值
	if err := eng.Execute(ruleBuilder, true); err != nil {
		t.Fatalf("execute rule error: %v\n", err)
	}
	t.Logf("student.Grade=%s\n", student.Grade)
}

func TestPool(t *testing.T) {
	student := &Student{} //必须是指针类型
	student.Score = 59    //修改分数以得到不同的等级

	var poolMinLen int64 = 3 //b站的代码这里有问题，你不能把poolMinLen设置的和poolMaxLen一样大（英语也是够chinglish的）
	var poolMaxLen int64 = 4 //合理值：cpu核数 * 2
	var model = 1            //1: sequence model; 2: concurrent model; 3: mix model; 4: inverse model
	apis := map[string]interface{}{
		"action1": action1,     //注入一个自定义函数
		"action2": action2,     //注入一个自定义函数
		"action3": action3,     //注入一个自定义函数
		"println": fmt.Println, //注入一个go的标准函数
	}

	//初始化规则引擎池，提供并发处理能力
	engPool, err := engine.NewGenginePool(poolMinLen, poolMaxLen, model, template, apis)
	if err != nil {
		t.Fatalf("build rule err: %v\n", err)
	}

	data := map[string]interface{}{
		"globalvar": 2,       //注入一个全局变量，所有规则里都可以拿到它
		"student":   student, //注入一个结构体变量
	}
	err, resultMap := engPool.Execute(data, true)
	if err != nil {
		t.Fatalf("execute rule error: %v\n", err)
	}
	t.Log(resultMap)

	//输出学生成绩等级指标
	t.Logf("student.Grade=%s\n", student.Grade)

	//修改参数再次测试
	fmt.Println("--------------------修改参数再次测试")
	student.Score = 100 //修改分数以得到不同的等级
	data = map[string]interface{}{
		"globalvar": 200,     //注入一个全局变量，所有规则里都可以拿到它
		"student":   student, //注入一个结构体变量
	}
	engPool.Execute(data, true)
	if err != nil {
		t.Fatalf("execute rule error: %v\n", err)
	}
	t.Logf("student.Grade=%s\n", student.Grade)
}

//这里一定要定义一个和插件里方法同名的接口
type Animal interface {
	Action()
}

func TestPlugin(t *testing.T) {
	// 1. open the so file to load the symbols
	plug, err := plugin.Open("plugin/plugin_Dog_d.so")
	if err != nil {
		t.Fatal(err)
	}

	// 2. look up a symbol (an exported function or variable)
	sym, err := plug.Lookup("Dog") //大写
	if err != nil {
		t.Fatal(err)
	}

	// 3. Assert that loaded symbol is of a desired type
	dog, ok := sym.(Animal)
	if !ok {
		t.Fatal("unexpected type from module symbol")
	}

	// 4. use the module
	dog.Action()
}

const template2 = `
rule "plugin gengine" "测试gengine基于插件的热加载"
begin
d.Action()
end
`

func TestPluginGengine(t *testing.T) {
	dataContext := context.NewDataContext()

	//PluginLoader这个方法扩展性设计的不太好，插件可能同时导出多个变量
	_, _, e := dataContext.PluginLoader("plugin/plugin_Dog_d.so") //Dog：插件导出的变量名；d：规则内使用的变量名
	if e != nil {
		panic(e)
	}

	//初始化规则引擎
	eng := engine.NewGengine()

	//构建规则
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	if err := ruleBuilder.BuildRuleFromString(template2); err != nil {
		t.Fatalf("build rule err: %v\n", err)
	}

	//执行规则
	if err := eng.Execute(ruleBuilder, true); err != nil {
		t.Fatalf("execute rule error: %v\n", err)
	}
}
