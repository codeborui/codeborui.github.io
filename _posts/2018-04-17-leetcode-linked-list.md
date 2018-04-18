---
layout: post
title:  "leetcode链表相关"
categories: 算法
tags:  leetcode 链表
author: Borui
---

* content
{:toc}

链表结构最常用的技巧就是快慢指针。

# [delete node in a linked list](https://leetcode-cn.com/problems/delete-node-in-a-linked-list/description/)
如果我们从头开始寻找当前节点的前一个节点，那么时间复杂度就是O(n)。
但是题目强调说被删除节点不是尾节点，这意味着当前节点一定有下一个节点。如果我们把下一个节点删除，然后下一个节点的值覆盖当前节点的值，其实效果是一样的。

# [reverse linked list](https://leetcode-cn.com/problems/reverse-linked-list/description/)
链表结构的基本功。通过三个指针，第一个指向新链表头，第二个指向旧链表头，第三个指针用来移动旧链表头。

# [reverse linked list ii](https://leetcode-cn.com/problems/reverse-linked-list-ii/description/)
基本思路和上一题类似，只不过我们需要先遍历到m的位置，然后翻转m到n之间的节点，然后再把n之后的节点连接上。