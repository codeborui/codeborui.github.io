---
layout: post
title:  "leetcode栈相关"
categories: 算法
tags:  leetcode 数组
author: Borui
---

* content
{:toc}

# [implement queue using stacks](https://leetcode-cn.com/problems/implement-queue-using-stacks/description/)
利用两个栈结构，两次先进后出=>先进先出.

# [implement stack using queues](https://leetcode-cn.com/problems/implement-stack-using-queues/description/)
这里其实就是利用遍历的方式，获取到最后一个元素，然后将之前的遍历放到另一个空队列，之后将交换两个队列的地位。