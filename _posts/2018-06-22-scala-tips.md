---
layout: post
title:  "scala tips"
categories: scala
tags:  tips 
author: Borui
---

* content
{:toc}

## scala: classes and objects
1. 可以将singleton object视为依附于某个object的命名tag.
2. 当singleton object继承*父类*或者*特征*时,singleton object就是该父类或者特征的*实例*,因此可以调用这些父类或特征的方法.
3. class可以传递参数,但singleton object不可以.singleton object被实现为synthetic class的实例,类似java静态类的初始化过程.

## 校验泛型参数的类型
```scala
abstract class RDD[T: ClassTag] {
    if (classOf(RDD[_]).isAssignableFrom(classTag[T].runtimeClass)) {
        ......
    }
}
```
## 简化函数
我们可以使用**下划线**作为一个或多个参数的占位符,只要在函数常量的内部,每个参数只被使用一次.
这里下划线可以代替一个参数,也可以代替整个参数列表.

## 闭包
scala的闭包能够追踪自由变量的变化.
```scala
var num = 1
val a = num + (_:Int)
a(1) // result is 2
num = 33
a(1) // result is 34
```