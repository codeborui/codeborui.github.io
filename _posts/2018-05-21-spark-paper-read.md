---
layout: post
title:  "spark paper read"
categories: spark
tags:  分布式计算 论文
author: Borui
---

* content
{:toc}

# Spark: Cluster Computing with Working Sets
## **Abstract**
分布式计算领域的map reduce模型已经获得巨大的成功，但该模型更适用于非循环数据流模型。有一类应用则是需要重用跨多个并行操作的工作集数据：比如迭代机器学习算法和交互式数据分析工具。Spark在保留Map reduce的可扩展性和容错性的基础上，提出RDD