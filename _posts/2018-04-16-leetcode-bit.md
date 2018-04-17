---
layout: post
title:  "leetcode位运算相关"
categories: 算法
tags:  leetcode 位运算
author: Borui
---

* content
{:toc}

# [number of 1 bits](https://leetcode-cn.com/problems/number-of-1-bits/description/)
最朴素的想法大家都能想到，就是不停地和0x01做与操作，然后无符号右移一位。

这里有个更好的想法，就是n&(n-1)，会把最右侧的1变成0.因此能执行几次这个操作，就说明有几个1。