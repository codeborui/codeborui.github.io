---
layout: post
title:  "leetcode树相关"
categories: 算法
tags:  leetcode 树
author: Borui
---

* content
{:toc}

树相关的算法一般都需要递归。

# [subtree-of-another-tree](https://leetcode-cn.com/problems/subtree-of-another-tree/description/)
这道题需要特别注意子树的概念，必须是从根节点到叶子结点(包括nil)都要完全匹配，才算是子树。

整体思路就是从被匹配树的根节点开始匹配，匹配不成功，则分别对左子树和右子树进行匹配。这里需要进行递归操作，退出条件就是被匹配的树已经为nil了。

而匹配的过程又是一层递归，先匹配根节点，然后匹配左子树，再匹配右子树。直到匹配到nil节点。

因此这道题最大的难点在于**两层递归**

#[invert binary tree](https://leetcode-cn.com/problems/invert-binary-tree/description/)
通过递归交换左右子树。