---
layout: post
title:  "leetcode数组相关"
categories: 算法
tags:  leetcode 数组
author: Borui
---

* content
{:toc}

# [Two sum](https://leetcode-cn.com/problems/two-sum/description/)
一般对于无序的数组问题,都需要考虑空间换取时间.

# [Median of tow sorted arrays](https://leetcode-cn.com/problems/median-of-two-sorted-arrays/description/)
由于是有序的数组,并且题目指定算法复杂度为O(log(m+n)),因此一定是一个类二分法问题.

此外,这题的难点在于将寻找中位数,转变成寻找第k大小的数.又因为保证复杂度为O(log(m+n)),因此就需要保证每次k值减半.

# [Containers with most water](https://leetcode-cn.com/problems/container-with-most-water/description/)
这道题的关键在于,容器面积的底和高都是变化的,我们需要找到最大的乘积值.

采用首尾双指针的策略,好处就在,寻找过程中底一直在减小(因为双指针不断靠近),因此我们如果想要找到最大的面积,就必须要求高一定要变大.在这个过程中,我们一定会寻找到最大的乘积值.

x*y=n,当x不断变小时,y值的变化趋向向将直接影响n的变化趋向.这样就把两个不确定的变化,转变成了一个不确定的

# [3sum](https://leetcode-cn.com/problems/3sum/description/)
这道题可对比两数之和这道题，这里需要陈列出所有答案。如果继续采用空间换时间，那么就需要固定第一个数，然后演变成在剩下的数组中求取特定两数之和。由于数组无序，因此如何排除重复答案就使得编程实现较为复杂。

此时，不如考虑对数组进行排序，排序的复杂度也就是O(nlog n)。然后固定第一个数，通过首尾指针就可以查找剩下两个数。

编程时，考虑重复的数需要有效地跳过。

# [3sum closest](https://leetcode-cn.com/problems/3sum-closest/description/)
这道题和上一道题基本思路一致，但是编写上需要细心