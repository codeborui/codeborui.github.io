---
layout: post
title:  "Java Concurrency in Practice(2)"
categories: Java
tags:  Concurrency 
author: Borui
---

* content
{:toc}

# 线程安全性
线程安全性的核心在于对状态访问操作需要管理,尤其是对**共享**的和**可变**状态的访问.

对象的状态还会包括依赖对象的域,例如容器的状态还包括容器装载的对象状态.

Java中的同步机制包括:
* synchronized
* volatile
* 显示锁
* 原子变量

修复多线程问题主要有三个方向:
* 不在线程之间共享状态变量: 不存在多线程,那么也就不存在多线程问题了.
* 将状态修改为不可变的变量: 不可变变量无论怎么访问,都是正确的.
* 在访问状态变量时使用同步: 当状态是共享的,且可变的时候,那么就只有采用同步机制了.

程序的状态封装的越好,越容易实现程序的线程安全性.

正确的编程理念是:
> 首先使得代码正确运行,再考虑性能问题.而且一定要有足够的证据证明必须提高性能.

## 2.1 什么是线程安全性
最核心的概念是**正确性**.对于某个类来说,正确性意味着符合规范规定.规范通常会定义各种**不变性条件**来约束对象状态和**后验条件**来描述结果.当多个线程访问某个类时,这个类始终表现出正确行为,那么这个类就是线程安全的.

无状态对象一定是线程安全的.

## 2.2 原子性
竞态条件(和数据竞争不是一回事):
* 先检查后执行: 常见例如延迟初始化(饱汉模式)
* 读取-修改-写入: 递增计数器.

复合操作: 包含了一组必须以原子方式执行的操作以确保线程的安全性.

线程安全问题主要是竞态条件导致的,原子性是解决多线程安全问题的基本策略.

## 2.3 加锁机制
要保持状态的一致性,就需要在单个原子操作中更新所有相关的状态变量.

内置锁-synchronized同步代码块:
* 锁对象的引用.
* 锁保护的代码块.

可重入-内置锁是可重入的,因此内置锁的粒度是线程而不是调用.
> 可重入锁的一种实现方式是,每个锁关联一个计数器和所有者线程.

## 2.4 用锁来保护状态
如果需要协调对某个变量的访问,那么所有访问的入口都需要使用同步.如果采用了加锁机制,那么要使用**同一个锁**.

## 2.5 活跃性和性能
简单粗粒度的锁会导致性能问题.

简单性和性能往往有冲突,因此需要权衡同步代码块的大小范围.

执行时间较长的操作一定不能持有锁,这会带来活跃性或性能问题.