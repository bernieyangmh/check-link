# check-link

> Golang tool for check website.



## Usage
```
$ checklink -o https://www.google.com https://www.facebook.com -l /home/user/log/todaylog.log -r /home/user/log/todayresult.log 
-o original link
-l log path
-r result path
```
### Meaning of log argument


| log title      |    Meaning  |
| --------: | --------:| :--: |
| aimUrl      |    start to check the url  |
| GET  | GET method  |
| CorrectlyRedict     |   redict is correctly	  |
| ErrorUrl      |    already in map use for avoiding duplicate or is not correct url |
| ErrorPath      |   wrong path  |
| Nil CrawlUrl      |    the contexts from href and src not match correct link |
| put      |   put the url into channel |
| url num is xxx      |    Total number of links  |

### example 
```
2018/03/07 09:14:10 commons.go:83: ErrorUrl		https://www.xxxx.com/aaa/7fabd3a1c59a9e76a8b71df63efdbfbfef1.png
2018/03/07 09:14:10 commons.go:83: ErrorUrl		https://www.xxxx.com/aaa/cc5a42b9338006dd1a26e4730e5e1e.png
2018/03/07 09:14:10 commons.go:83: ErrorUrl		https://www.xxxx.com/aaa/2a6168714522ee19dcb1c1.png
2018/03/07 09:14:10 commons.go:83: ErrorUrl		https://www.xxxx.com/aaa/a700a755dde294347fc0f0.png
2018/03/07 09:14:10 commons.go:83: ErrorUrl		https://www.xxxx.com/aaa/c4322a6168714522ee19dcb1c1.png
2018/03/07 09:14:10 commons.go:83: ErrorUrl		https://www.xxxx.com/aaa/d168ab4f91b5a700a755dde294347fc0f0.png
2018/03/07 09:14:10 commons.go:83: ErrorUrl		https://www.xxxx.com/aaa/c5a42b9338006dd1a26e4730e5e1e.png
2018/03/07 09:14:10 commons.go:83: ErrorUrl		https://xxxx-xxxx.xxxx.com/xxxx.js
155
2018/03/07 09:14:10 nets.go:53: GET		https://www.xxxx.com/aaa/haha
2018/03/07 09:14:10 nets.go:75: GetForBoby		https://www.xxxx.com/aaa/hello
2018/03/07 09:14:11 nets.go:41: aimUrl		https://www.xxxx.com/aaa/hi
2018/03/07 09:14:11 commons.go:229: ErrorUrl		https://www.xxxx.com/
2018/03/07 09:14:11 commons.go:232: ErrorPath			#
2018/03/07 09:14:11 commons.go:229: ErrorUrl		https://www.xxxx.com/aaa/what
2018/03/07 09:14:11 commons.go:229: ErrorUrl		https://www.xxxx.com/aaa/where
2018/03/07 09:14:11 commons.go:229: ErrorUrl		https://www.xxxx.com/aaa/who
2018/03/07 09:14:11 commons.go:229: ErrorUrl		https://www.xxxx.com/aaa/why
2018/03/07 09:14:11 commons.go:229: ErrorUrl		https://www.qiniu.com/aaa/when
```

```
Error Link		https://xxxx.abc.com/frombbbbbb
Ref Link		https://xxxx.abc.com/bbbbb
StatusCode  404
Exception		


Error Link		https://xxxx.abc.com/fromaaaa
Ref Link		https://xxxx.abc.com/aaaa
StatusCode    -2	
Exception		Get https://xxxx.abc.com/fromaaaa x509: certificate has expired or is not yet valid
```
