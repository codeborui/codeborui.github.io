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

# [invert binary tree](https://leetcode-cn.com/problems/invert-binary-tree/description/)
通过递归交换左右子树。

# [binary tree level order traversal](https://leetcode-cn.com/problems/binary-tree-level-order-traversal/description/)
两个考察点：一个是树的广度遍历（通过队列实现），还有一个就是如何判断哪些节点在一个层次。

树的广度遍历，首先插入头结点，然后退出条件设置为队列为空的时候，然后左子节点和右子节点不断入队列。

判断层次的话有两种做法：第一种：插入一个nil标记，读到nil的时候说明当前层次遍历完毕，然后将nil重新插回队尾。第二种：在遍历开始前，获悉队列中数据个数，即为当前层次个数，然后循环到次数后进入下一个层次。

# [binary tree zigzag level order traversal](https://leetcode-cn.com/problems/binary-tree-zigzag-level-order-traversal/description/)
这道题直接在上一道题的基础上做改动，设置一个方向flag，指明当前是从左到右遍历，还是从右到左遍历。然后在输出时选择尾插还是头插。

# [binary tree level order traversal ii](https://leetcode-cn.com/problems/binary-tree-level-order-traversal-ii/description/)
类似上一个思路，将结果头插