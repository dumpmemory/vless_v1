# 本仓库文件已经过时，对vless v1的最新定义请直接查看 [verysimple仓库](https://github.com/hahahrfool/v2ray_simple/blob/main/vless_v1.md)

# 详情
提取自我的v2simple

v2fly社区某些谨慎的人非说有协议问题，那我只好暂且把我自己写的代码提取出来

这样对于v2fly来说绝对百分百没有协议问题



我再全面总结一下v2simple的开源协议问题。

首先声明，退一万步讲，我的vless_v1文档以及整个vless文件夹，以及 udp.go 等文件，都是我自己创建的，基于MIT许可证；

所以现在就可以直接使用我vless_v1项目的所有代码。

有任何不懂的可以直接阅读我的 [v2simple](https://github.com/hahahrfool/v2simple) 代码。


# 完全匿名的找不到实控人的没附带任何开源协议的github项目的fork问题

原v2simple项目并没有附带任何开源许可证，根据下面论证，是可以fork的

论证开始

原作者github里只有一个库，早已跑路，他也未实名，因此法律保护的主体找不到，我们每一个fork的人都可以声称自己是原作者，只是换个号；

此时是无法判定侵权的。现在我就声称我就是原作者，换号的理由就是我忘了密码了。然后因为忘了密码，所以也无法登陆该账号去更新。

如果原作者突然又冒出来了，那么声称自己是原作者的路就封死了；
但是我相信他是不会反对我fork他的代码的，会给予我特许。毕竟我把他的代码发扬光大了。

论证结束

## 司法实践

总之大家要注重司法实践。当一个人完全匿名，找不到踪影时，你是无法说你对他侵权的，因为“他”不存在。因为这个匿名账号的实际控制人你是找不到的。根据著作权法，匿名的著作权是实控人的，而现在实控人都找不到，所以这个东西完全就是公共领域的，相当于大自然的馈赠。逻辑还是一样，任何人都可以声称他是原来账户的实控人，甚至该账户的实控人可能是AI，可能是外星人。

使用了开源协议则有不同，因为开源协议提出了具体要求，而本例原作者没有使用任何协议。

再举个例子，石头打死了小猫，你不能认为“石头”违法，因为法律是针对“人”的，所以要找扔石头的人，如果找不到人，那就是大自然的法则，石头自己从楼上被风吹下砸死的小猫，你总不能说“风”违法，说“大自然”违法吧。 因为这个匿名用户已经找不到，也没有开源协议，所以这个代码就进入了公共领域，进入了大自然。

准确地说，原v2simple代码属于 No Lisence and can't find the author at all. 这和 No Lisence 是不同的， No Lisence 默认尊重原作者的著作权， 但是原作者完全是一个匿名账户，无法追溯实控人，所以 No Lisence 对原作者的尊重也就没法实现。尊重原作者 等价于 尊重大自然。因为原作者 在这种情况下等价于 大自然。

我也是尊重原作者才fork v2simple的。这个代码原理很简单，完全可以自己重写一遍。我只是不想居功而已


# 再说一遍

不用担心使用问题，因为你们只需观察v2simple中我自己commit的部分

我的所有commit都是我自己的，不是原作者的；我的commit基于MIT协议的

还有，就算你担心版权那也是不成立的，因为没有人可以把v2simple的代码直接复制粘贴到 v2ray里然后能直接运行。必须学习其原理，然后自行实现

这也是我希望大佬能指导我的原因，我不懂v2ray的内部结构，导致无法直接应用v2simple的原理
