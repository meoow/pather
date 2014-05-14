# basename, dirname, readlink and more.

```
Usage: {{.Path}} (a|x|b|d|r) [-0] [--] [-|PATH...]

{{.Path}} is a simple file/dir path parsing tool. 

Commands:
  a         Absolute path
  x         File extension
  b         Base name
  d         Parent Directory
  r         Real path of symlink

Arguments:
PATH        path of file or directory

Options:
  -h        Show this help
  -0        Paths are seperated by NUL charater
```