---
layout: post
title:  "Distributed Systems: Concepts and Design, First Edition(1)"
categories: DS-C&D
tags:  分布式 
author: Borui
---

* content
{:toc}

# 分布式系统的特征
1. 分布式系统是组件分布在通过网络相连接的计算机上,组件之间通过消息传递来进行通信和协调动作的系统.
2. 分布式系统重要特征:   
    + 组件的并发性
    + 全局时钟的缺失
    + 组件故障相互独立.
3. 资源共享是构造分布式系统的主要动机.
4. 资源可以被服务器管理,被客户端访问,或者被encapsulate成对象,被其他的客户端对象访问.
5. 构建分布式系统的挑战在于:
    + 组件的异构性(hetero|geneity[异性的|])
    + 开放性(允许增加或替换组件)
    + 安全性
    + 可扩展性(当负载或者用户量增加时,良好工作的能力)
    + 故障处理
    + 组件并发性
    + 透明性(transparency)
    + 提供的服务质量
6. 