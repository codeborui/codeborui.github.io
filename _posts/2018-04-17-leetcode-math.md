---
layout: post
title:  "leetcode数学相关"
categories: 算法
tags:  leetcode 数学
author: Borui
---

* content
{:toc}

# [pow x n](https://leetcode-cn.com/problems/powx-n/description/)
这里主要注意两点：
1. 由于x^n里当n为负时，结果为1/(x^-n)，那么当n为x时，需要额外注意。由于x是float类型，判断0的时候需要注意到系统的精度。判0的标准是和0相差小于10^-6。
2. 为了加快运算，x^n=(x^(n/2))^2 * (x^(n%2))

此外，针对除以2的操作，我们采用右移1位。判断奇偶的操作，我们和0x1进行与操作。