# gengine practice

### Build plugin (for hot update)

```
go build -buildmode=plugin -o plugin/plugin_Dog_d.so plugin/plugin_dog.go
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
action1
action2
rule name: rule1 local var: 1 global var: 2
action1
action2
action3
    gengine_test.go:103: map[rule1:hello rule2:world]
    gengine_test.go:106: student.Grade=你上升空间巨大，前途无量啊！
--------------------修改参数再次测试
rule name: rule2 local var: 3 global var: 200
action3
action2
action1
rule name: rule1 local var: 1 global var: 200
action1
action2
action3
    gengine_test.go:115: student.Grade=你简直是个天才
--- PASS: TestSingle (0.04s)
=== RUN   TestPool
rule name: rule2 local var: 3 global var: 2
action3
action1
action2
rule name: rule1 local var: 1 global var: 2
action1
action2
action3
    gengine_test.go:146: map[rule1:hello rule2:world]
    gengine_test.go:149: student.Grade=你上升空间巨大，前途无量啊！
--------------------修改参数再次测试
rule name: rule2 local var: 3 global var: 200
action3
action2
action1
rule name: rule1 local var: 1 global var: 200
action1
action2
action3
    gengine_test.go:162: student.Grade=你简直是个天才
--- PASS: TestPool (0.02s)
PASS
```
