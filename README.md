# msgboat



OpenFilWallet offline transaction submitter. 



Only transactions encoded by `OpenFilWallet/OpenFilWallet/chain` can be recognized. If the submission is successful, the transaction cid is returned, otherwise an error is returned.



**Build**

```
git clone git@github.com:OpenFilWallet/msgboat.git
cd msgboat
go mod tidy && go build -o msgboat main.go
./msgboad run
```

