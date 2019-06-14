# DRONE-CI SSH 插件

- 使用

使用私钥登录

```yaml
kind: pipeline
name: default
steps:
  - name: publish
    image: droneplugin/ssh
    privileged: true
    settings:
      host: {Your-Remote-Host}
      user: {Your-Remote-User}
      pem: {Your-Private-Key}
      command:
        - ls
        - echo hello
```

使用私钥文件登录

```yaml
kind: pipeline
name: default
steps:
  - name: publish
    image: droneplugin/ssh
    privileged: true
    settings:
      host: {Your-Remote-Host}
      user: {Your-Remote-User}
      pem_file: {Your-Private-Key-Path}
      command:
        - ls
        - echo hello
```

带密码的私钥登录

```yaml
kind: pipeline
name: default
steps:
  - name: publish
    image: droneplugin/ssh
    privileged: true
    settings:
      host: {Your-Remote-Host}
      user: {Your-Remote-User}
      pem: {Your-Private-Key}
      passphrase: {Your-Private-Passphrase}
      command:
        - ls
        - echo hello
```

使用密码登录

```yaml
kind: pipeline
name: default
steps:
  - name: publish
    image: droneplugin/ssh
    privileged: true
    settings:
      host: {Your-Remote-Host}
      user: {Your-Remote-User}
      password: {Your-Remote-Password}
      command:
        - ls
        - echo hello
```