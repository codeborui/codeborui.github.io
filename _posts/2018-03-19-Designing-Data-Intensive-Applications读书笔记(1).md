---
layout: post
title:  "Designing Data-Intensive Applications读书笔记(1)"
categories: 分布式
tags:  读书笔记 分布式
author: Borui
---

* content
{:toc}

# Reliable, Scalable, and Maintainable Applications
数据密集型应用特征: CPU能力很少是这类应用的限制因素,反而是数据量,数据的复杂度以及数据的变化速度.

数据密集型应用通常是由多种通用功能模块组合而成,这些功能包括:
> 1.数据库<br/>2.缓存<br/>3.索引<br/>4.流处理<br/>5.批处理

![Figure 1-1. One possible architecture for a data system that combines several
components.](https://raw.githubusercontent.com/codeborui/codeborui.github.io/master/img/1.JPG)
