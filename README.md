# gengine practice

### Build plugin (for hot update)

```
go build -buildmode=plugin -o plugin/plugin_Dog_d.so plugin/plugin_dog.go
go build -buildmode=plugin -o plugin/plugin_Cat_c.so plugin/plugin_cat.go
```

### Run Test

```
go test -v test/gengine_test.go
```

### Output

```
=== RUN   TestSingle
rule name: rule2 local var: 3 global var: 2
action3
action2
action1
rule name: rule1 local var: 1 global var: 2
action1
action2
action3
    gengine_test.go:104: map[rule1:hello rule2:world]
    gengine_test.go:107: student.Grade=你上升空间巨大，前途无量啊！
--------------------修改参数再次测试
rule name: rule2 local var: 3 global var: 200
action3
action1
action2
rule name: rule1 local var: 1 global var: 200
action1
action2
action3
    gengine_test.go:116: student.Grade=你简直是个天才
--- PASS: TestSingle (0.02s)
=== RUN   TestPool
rule name: rule2 local var: 3 global var: 2
action3
action1
action2
rule name: rule1 local var: 1 global var: 2
action1
action2
action3
    gengine_test.go:147: map[rule1:hello rule2:world]
    gengine_test.go:150: student.Grade=你上升空间巨大，前途无量啊！
--------------------修改参数再次测试
rule name: rule2 local var: 3 global var: 200
action3
action1
action2
rule name: rule1 local var: 1 global var: 200
action1
action2
action3
    gengine_test.go:163: student.Grade=你简直是个天才
--- PASS: TestPool (0.02s)
=== RUN   TestPlugin
汪汪汪
--- PASS: TestPlugin (0.00s)
=== RUN   TestPluginGengine
汪汪汪
--- PASS: TestPluginGengine (0.00s)
PASS
ok  	command-line-arguments	0.057s
```
