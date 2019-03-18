# logx

## Example
You can use default Logger to write log
```
func main() {
    ch := logx.NewConsoleHandler()
    logx.AddHandler(ch)
    
    fh, err := logx.NewFileHandler("./main.log")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    logx.AddHandler(fh)
    
    logx.SetLevel(logx.DEBUG)

    logx.Debug("debug")
    logx.Debugf("%s", "debugf")

    logx.Info("info")
    logx.Infof("%s", "infof")
    
    logx.Warn("warn")
    logx.Warnf("%s", "warnf")
    
    logx.Error("error")
    logx.Errorf("%s", "errorf")
}
```
Or you can create new Logger object to write log
```
func main() {
    logger := logx.NewLogger()

    ch := logx.NewConsoleHandler()
    logger.AddHandler(ch)

    fh, err := logx.NewFileHandler("./main.log")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    logger.AddHandler(fh)

    logger.SetLevel(logx.DEBUG)

    logger.Debug("debug")
    logger.Debugf("%s", "debugf")

    logger.Info("info")
    logger.Infof("%s", "infof")
    
    logx.Warn("warn")
    logx.Warnf("%s", "warnf")

    logx.Error("error")
    logx.Errorf("%s", "errorf")
}
```

## Output
```
[DEBUG]:[2019-03-18 15:44:07]:[main.go:22] debug
[DEBUG]:[2019-03-18 15:44:07]:[main.go:23] debugf
[INFO ]:[2019-03-18 15:44:07]:[main.go:25] info
[INFO ]:[2019-03-18 15:44:07]:[main.go:26] infof
[WARN ]:[2019-03-18 15:44:07]:[main.go:28] warn
[WARN ]:[2019-03-18 15:44:07]:[main.go:29] warnf
[ERROR]:[2019-03-18 15:44:07]:[main.go:31] error
[ERROR]:[2019-03-18 15:44:07]:[main.go:32] errorf
```