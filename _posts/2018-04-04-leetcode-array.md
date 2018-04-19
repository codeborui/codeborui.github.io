---
layout: post
title:  "leetcode数组相关"
categories: 算法
tags:  leetcode 数组
author: Borui
---

* content
{:toc}

数组最常用的技巧就是首尾双指针，二分和hash表。

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

# [search a 2d matrix](https://leetcode-cn.com/problems/search-a-2d-matrix/description/)
这道题目的核心在于，如何使得我们每次**选择查询的方向是唯一**的。

注意数组的特性，从左到右是有序的，从上到下是有序的，而且每行第一个数大于上一行最后一个数，也就是说每行的数都大于上一行的所有数。当前数小于目标时，可以往左走，也可以往上走。而当前数大于目标时，可以往右走，也可以往下走。

根据我们唯一方向的分析，可以组合出两种策略：
1. 当前数小于目标时，往左走；前数大于目标时，往下走。这时我们的最初位置就是右上角。
2. 当前数小于目标时，往上走；前数大于目标时，往右走。这时我们的最初位置就是左下角。

以上。

# [find minimum in rotated sorted array](https://leetcode-cn.com/problems/find-minimum-in-rotated-sorted-array/description/)
首先，还是观察数组的特征：有序，无重复，旋转后的右侧小于左侧值。

针对这种有序无重复的数组，我们一般会先考虑二分法。由于我们的目标是寻找最小值，因此我们需要不断寻找最小值的范围。这里mid我们向下取整，因此当只有两个数时，中间值将是left的值，我们可以单独考虑两个数的情况，如下只考虑大于等于三个数时的情况：
1. mid大于left值，并且mid大于right值，则情况如下：4,5,6,7,0,1,2,3。目标在右侧。
2. mid大于left值，并且mid小于right值，则情况如下：0,1,2,3,4,5,6,7。目标在左侧。
3. mid小于left值，并且mid小于right值，则情况如下：5,6,7,0,1,2,3,4。目标在左侧。
4. mid小于left值，并且mid大于right值，这种情况不可能，mid小于left说明mid处于翻转后的右半侧，右半侧是有序的，mid只可能小于right。

综上，当mid值小于right值，目标在左半侧。当mid大于left值时，目标在右半侧。

# [find minimum in rotated sorted array ii](https://leetcode-cn.com/problems/find-minimum-in-rotated-sorted-array-ii/description/)
这里存在重复的数据了，不过上述的分析还是有效的，但是关于mid与left，right相等的情况，就需要特别考虑了。

当mid=left=right的时候，这个时候最小值可能在任何一个区间里，如1,1,1,0,1或者1,0,1,1,1，因此只能去顺序遍历了。
1. mid等于left值，并且mid大于right值，则情况如下：3,3,0,1。目标在右侧。
2. mid等于left值，并且mid小于right值，则情况如下：0,0,0,1。目标在左侧。
3. mid等于right值，并且mid大于left值，则情况如下：0,1,1,1。目标在左侧。
4. mid等于right值，并且mid小于left值，则情况如下：2,0,1,1。目标在左侧。

# [spiral-matrix](https://leetcode-cn.com/problems/spiral-matrix/description/)
这道题没有算法难度，但是存在编程难度。

比较精妙的解法是，分别用四个变量来限定边界，然后通过螺旋访问顺序，依次变化边界。

比较通用的方法，就是利用方向数组。

# [spiral matrix ii](https://leetcode-cn.com/problems/spiral-matrix-ii/description/)
和上面的思路一致，通过螺旋访问数组然后赋值。