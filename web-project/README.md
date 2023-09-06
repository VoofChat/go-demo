### 数据表生成
```shell

tableList := [blog, comment, user]

gen --sqltype=mysql \
    --connstr="root:qwerasdf@tcp(127.0.0.1:3306)/demo" \
    --database="demo" \
    --table="user" \
    --gorm \
    --json \
    --guregu \
    --rest \
    --overwrite
```
